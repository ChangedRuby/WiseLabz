package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// RunbookRecord represents a row in the runbooks table: a binding from a
// change type, alert severity, or finding check type to operator guidance,
// plus optional pointers to a known-good snapshot and a doc, and zero or
// more steps (RunbookStepRecord) referencing connector lifecycle
// operations. Linking a runbook — and its steps — to a target grants no
// mutation permission by itself: executing a step still requires the
// caller to hold an operator grant on the step's connector and to
// step-up/confirm through the same elevation flow as a direct connector
// restart/start/stop. SnapshotID/DocID remain inert references only.
type RunbookRecord struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Body        string  `json:"body"`
	TargetType  string  `json:"targetType"`
	TargetValue string  `json:"targetValue"`
	SnapshotID  *string `json:"snapshotId"`
	DocID       *string `json:"docId"`
	CreatedAt   string  `json:"createdAt"`
	UpdatedAt   string  `json:"updatedAt"`
}

const runbookColumns = `id, title, body, target_type, target_value, snapshot_id, doc_id, created_at, updated_at`

// CreateRunbook inserts a new runbook.
func (s *Store) CreateRunbook(ctx context.Context, r *RunbookRecord) (*RunbookRecord, error) {
	if r.ID == "" {
		r.ID = uuid.New().String()
	}
	now := time.Now().UTC().Format(time.RFC3339)
	if r.CreatedAt == "" {
		r.CreatedAt = now
	}
	if r.UpdatedAt == "" {
		r.UpdatedAt = now
	}

	_, err := s.db.ExecContext(ctx, `
		INSERT INTO runbooks (id, title, body, target_type, target_value, snapshot_id, doc_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, r.ID, r.Title, r.Body, r.TargetType, r.TargetValue, r.SnapshotID, r.DocID, r.CreatedAt, r.UpdatedAt)
	if err != nil {
		if isUniqueViolation(err) {
			return nil, ErrConflict
		}
		return nil, fmt.Errorf("create runbook: %w", err)
	}
	return r, nil
}

// GetRunbook retrieves a runbook by ID.
func (s *Store) GetRunbook(ctx context.Context, id string) (*RunbookRecord, error) {
	r, err := scanRunbook(s.db.QueryRowContext(ctx, `SELECT `+runbookColumns+` FROM runbooks WHERE id = ?`, id))
	if err != nil {
		return nil, fmt.Errorf("get runbook: %w", err)
	}
	return r, nil
}

// GetRunbookByTarget retrieves the runbook bound to a change type or alert
// severity, using the unique index on (target_type, target_value).
func (s *Store) GetRunbookByTarget(ctx context.Context, targetType, targetValue string) (*RunbookRecord, error) {
	r, err := scanRunbook(s.db.QueryRowContext(ctx,
		`SELECT `+runbookColumns+` FROM runbooks WHERE target_type = ? AND target_value = ?`,
		targetType, targetValue))
	if err != nil {
		return nil, fmt.Errorf("get runbook by target: %w", err)
	}
	return r, nil
}

// ListRunbooks returns all runbooks, newest first.
func (s *Store) ListRunbooks(ctx context.Context) ([]*RunbookRecord, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT `+runbookColumns+` FROM runbooks ORDER BY created_at DESC`)
	if err != nil {
		return nil, fmt.Errorf("list runbooks: %w", err)
	}
	defer rows.Close() //nolint:errcheck

	runbooks := []*RunbookRecord{}
	for rows.Next() {
		r, err := scanRunbook(rows)
		if err != nil {
			return nil, fmt.Errorf("scan runbook: %w", err)
		}
		runbooks = append(runbooks, r)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate runbooks: %w", err)
	}
	return runbooks, nil
}

// UpdateRunbook updates fields on an existing runbook from a partial map of
// column values, then returns the refreshed record.
func (s *Store) UpdateRunbook(ctx context.Context, id string, updates map[string]any) (*RunbookRecord, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	args := []any{now}
	var parts []string

	for k, v := range updates {
		switch k {
		case "title":
			parts = append(parts, "title = ?")
			args = append(args, v)
		case "body":
			parts = append(parts, "body = ?")
			args = append(args, v)
		case "target_type":
			parts = append(parts, "target_type = ?")
			args = append(args, v)
		case "target_value":
			parts = append(parts, "target_value = ?")
			args = append(args, v)
		case "snapshot_id":
			parts = append(parts, "snapshot_id = ?")
			args = append(args, v)
		case "doc_id":
			parts = append(parts, "doc_id = ?")
			args = append(args, v)
		}
	}

	query := "UPDATE runbooks SET updated_at = ?"
	if len(parts) > 0 {
		query += ", " + strings.Join(parts, ", ")
	}
	query += " WHERE id = ?"
	args = append(args, id)

	result, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		if isUniqueViolation(err) {
			return nil, ErrConflict
		}
		return nil, fmt.Errorf("update runbook: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("rows affected: %w", err)
	}
	if rows == 0 {
		return nil, ErrNotFound
	}
	return s.GetRunbook(ctx, id)
}

// DeleteRunbook deletes a runbook by ID.
func (s *Store) DeleteRunbook(ctx context.Context, id string) error {
	result, err := s.db.ExecContext(ctx, `DELETE FROM runbooks WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete runbook: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

// CreateRunbookWithSteps inserts r and replaces its steps in one
// transaction (all-or-nothing), returning the created runbook and its
// persisted steps in position order.
func (s *Store) CreateRunbookWithSteps(ctx context.Context, r *RunbookRecord, steps []*RunbookStepRecord) (*RunbookRecord, []*RunbookStepRecord, error) {
	var created *RunbookRecord
	var saved []*RunbookStepRecord
	err := s.WithinTransaction(ctx, func(tx *Store) error {
		var err error
		created, err = tx.CreateRunbook(ctx, r)
		if err != nil {
			return err
		}
		saved, err = tx.ReplaceRunbookSteps(ctx, created.ID, steps)
		return err
	})
	if err != nil {
		return nil, nil, err
	}
	return created, saved, nil
}

// UpdateRunbookWithSteps updates fields on an existing runbook and, when
// replaceSteps is true, replaces its steps — both in one transaction. When
// replaceSteps is false, steps is ignored and the runbook's current steps
// are returned unchanged (the "steps" key was absent from the request).
func (s *Store) UpdateRunbookWithSteps(ctx context.Context, id string, updates map[string]any, steps []*RunbookStepRecord, replaceSteps bool) (*RunbookRecord, []*RunbookStepRecord, error) {
	var updated *RunbookRecord
	var saved []*RunbookStepRecord
	err := s.WithinTransaction(ctx, func(tx *Store) error {
		var err error
		updated, err = tx.UpdateRunbook(ctx, id, updates)
		if err != nil {
			return err
		}
		if replaceSteps {
			saved, err = tx.ReplaceRunbookSteps(ctx, id, steps)
			return err
		}
		saved, err = tx.ListRunbookStepsFor(ctx, id)
		return err
	})
	if err != nil {
		return nil, nil, err
	}
	return updated, saved, nil
}

func scanRunbook(row rowScanner) (*RunbookRecord, error) {
	var r RunbookRecord
	var snapshotID, docID sql.NullString
	err := row.Scan(&r.ID, &r.Title, &r.Body, &r.TargetType, &r.TargetValue,
		&snapshotID, &docID, &r.CreatedAt, &r.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	if snapshotID.Valid {
		r.SnapshotID = &snapshotID.String
	}
	if docID.Valid {
		r.DocID = &docID.String
	}
	return &r, nil
}

// RunbookStepRecord represents a row in the runbook_steps table: one
// lifecycle operation (restart/start/stop) a runbook points at, bound to a
// specific connector and, optionally, a specific entity on that connector.
// Authoring a step grants no mutation permission by itself — see
// RunbookRecord's doc comment.
type RunbookStepRecord struct {
	ID          string `json:"id"`
	RunbookID   string `json:"runbookId"`
	Position    int    `json:"position"`
	Title       string `json:"title"`
	ConnectorID string `json:"connectorId"`
	Verb        string `json:"verb"`
	EntityRef   string `json:"entityRef"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

const runbookStepColumns = `id, runbook_id, position, title, connector_id, verb, entity_ref, created_at, updated_at`

// ListRunbookSteps returns the steps belonging to any of runbookIDs,
// grouped by runbook ID and ordered by position within each group. Missing
// or empty runbookIDs yield an empty map.
func (s *Store) ListRunbookSteps(ctx context.Context, runbookIDs []string) (map[string][]*RunbookStepRecord, error) {
	out := map[string][]*RunbookStepRecord{}
	if len(runbookIDs) == 0 {
		return out, nil
	}

	placeholders := make([]string, len(runbookIDs))
	args := make([]any, len(runbookIDs))
	for i, id := range runbookIDs {
		placeholders[i] = "?"
		args[i] = id
	}

	rows, err := s.db.QueryContext(ctx,
		`SELECT `+runbookStepColumns+` FROM runbook_steps WHERE runbook_id IN (`+strings.Join(placeholders, ",")+`) ORDER BY runbook_id, position`,
		args...)
	if err != nil {
		return nil, fmt.Errorf("list runbook steps: %w", err)
	}
	defer rows.Close() //nolint:errcheck

	for rows.Next() {
		st, err := scanRunbookStep(rows)
		if err != nil {
			return nil, fmt.Errorf("scan runbook step: %w", err)
		}
		out[st.RunbookID] = append(out[st.RunbookID], st)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate runbook steps: %w", err)
	}
	return out, nil
}

// ListRunbookStepsFor returns runbookID's steps, ordered by position.
func (s *Store) ListRunbookStepsFor(ctx context.Context, runbookID string) ([]*RunbookStepRecord, error) {
	m, err := s.ListRunbookSteps(ctx, []string{runbookID})
	if err != nil {
		return nil, err
	}
	return m[runbookID], nil
}

// GetRunbookStep retrieves one step scoped to runbookID — a step ID that
// exists but belongs to a different runbook is treated as not found.
func (s *Store) GetRunbookStep(ctx context.Context, runbookID, stepID string) (*RunbookStepRecord, error) {
	st, err := scanRunbookStep(s.db.QueryRowContext(ctx,
		`SELECT `+runbookStepColumns+` FROM runbook_steps WHERE runbook_id = ? AND id = ?`,
		runbookID, stepID))
	if err != nil {
		return nil, fmt.Errorf("get runbook step: %w", err)
	}
	return st, nil
}

// ReplaceRunbookSteps deletes all of runbookID's existing steps and inserts
// steps in their given order (position = slice index). A step whose ID
// already belonged to this runbook is kept; any other ID (empty, or
// belonging to a different runbook or a step no longer present) is
// replaced with a freshly generated one, so callers can't smuggle in
// another runbook's step ID. Callers wanting this atomic with the parent
// runbook write should run it inside s.WithinTransaction (see
// CreateRunbookWithSteps/UpdateRunbookWithSteps).
func (s *Store) ReplaceRunbookSteps(ctx context.Context, runbookID string, steps []*RunbookStepRecord) ([]*RunbookStepRecord, error) {
	existing, err := s.ListRunbookStepsFor(ctx, runbookID)
	if err != nil {
		return nil, err
	}
	existingIDs := make(map[string]bool, len(existing))
	for _, st := range existing {
		existingIDs[st.ID] = true
	}

	if _, err := s.db.ExecContext(ctx, `DELETE FROM runbook_steps WHERE runbook_id = ?`, runbookID); err != nil {
		return nil, fmt.Errorf("delete runbook steps: %w", err)
	}

	now := time.Now().UTC().Format(time.RFC3339)
	saved := make([]*RunbookStepRecord, 0, len(steps))
	for i, st := range steps {
		id := st.ID
		if id == "" || !existingIDs[id] {
			id = uuid.New().String()
		}
		rec := &RunbookStepRecord{
			ID:          id,
			RunbookID:   runbookID,
			Position:    i,
			Title:       st.Title,
			ConnectorID: st.ConnectorID,
			Verb:        st.Verb,
			EntityRef:   st.EntityRef,
			CreatedAt:   now,
			UpdatedAt:   now,
		}
		if _, err := s.db.ExecContext(ctx, `
			INSERT INTO runbook_steps (id, runbook_id, position, title, connector_id, verb, entity_ref, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, rec.ID, rec.RunbookID, rec.Position, rec.Title, rec.ConnectorID, rec.Verb, rec.EntityRef, rec.CreatedAt, rec.UpdatedAt); err != nil {
			return nil, fmt.Errorf("insert runbook step: %w", err)
		}
		saved = append(saved, rec)
	}
	return saved, nil
}

func scanRunbookStep(row rowScanner) (*RunbookStepRecord, error) {
	var st RunbookStepRecord
	err := row.Scan(&st.ID, &st.RunbookID, &st.Position, &st.Title, &st.ConnectorID, &st.Verb, &st.EntityRef, &st.CreatedAt, &st.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &st, nil
}
