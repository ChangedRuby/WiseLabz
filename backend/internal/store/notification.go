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

// NotificationRecord represents a row in the in_app_notifications table.
type NotificationRecord struct {
	ID        string `json:"id"`
	UserID    string `json:"-"`
	AlertID   string `json:"alertId,omitempty"` // empty = NULL, no deep-link target
	EventType string `json:"eventType"`
	Title     string `json:"title"`
	Message   string `json:"message"`
	Read      bool   `json:"read"`
	CreatedAt string `json:"createdAt"`
}

// CreateNotification inserts a new in-app notification.
func (s *Store) CreateNotification(ctx context.Context, n *NotificationRecord) error {
	if n.ID == "" {
		n.ID = uuid.New().String()
	}
	if n.CreatedAt == "" {
		n.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	}

	_, err := s.db.ExecContext(ctx, `
		INSERT INTO in_app_notifications (id, user_id, alert_id, event_type, title, message, read, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, n.ID, n.UserID, nilToStr(n.AlertID), n.EventType, n.Title, n.Message, boolToInt(n.Read), n.CreatedAt)
	if err != nil {
		return fmt.Errorf("create notification: %w", err)
	}
	return nil
}

// ListNotifications returns a paginated list of a user's notifications, most recent first.
func (s *Store) ListNotifications(ctx context.Context, userID string, unreadOnly bool, offset, limit int) ([]NotificationRecord, int, error) {
	where := "WHERE user_id = ?"
	args := []any{userID}
	if unreadOnly {
		where += " AND read = 0"
	}

	return paginatedQuery(ctx, s.db, "in_app_notifications", notificationColumns, where, args, "created_at DESC", limit, offset, scanNotification)
}

// notificationColumns is the shared column list for every
// in_app_notifications SELECT.
const notificationColumns = `id, user_id, alert_id, event_type, title, message, read, created_at`

// ListNotificationsSince returns a user's notifications with the given event
// types created at or after `since` (RFC3339), most recent first. Used by the
// digest sweep to pull everything accumulated since the user's last digest.
// If since is empty, returns all notifications of the given event types.
func (s *Store) ListNotificationsSince(ctx context.Context, userID, since string, eventTypes []string) ([]NotificationRecord, error) {
	if len(eventTypes) == 0 {
		return []NotificationRecord{}, nil
	}
	placeholders := make([]string, len(eventTypes))
	args := []any{userID}
	for i, et := range eventTypes {
		placeholders[i] = "?"
		args = append(args, et)
	}

	query := `SELECT ` + notificationColumns + `
		FROM in_app_notifications
		WHERE user_id = ? AND event_type IN (` + strings.Join(placeholders, ",") + `)`

	if since != "" {
		query += ` AND created_at >= ?`
		args = append(args, since)
	}

	query += ` ORDER BY created_at DESC`

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list notifications since: %w", err)
	}
	defer rows.Close() //nolint:errcheck

	var notifications []NotificationRecord
	for rows.Next() {
		n, err := scanNotification(rows)
		if err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		notifications = append(notifications, n)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate notifications: %w", err)
	}
	if notifications == nil {
		notifications = []NotificationRecord{}
	}
	return notifications, nil
}

// MarkNotificationRead marks one of a user's notifications as read and returns
// the updated record. Scoping the UPDATE by userID means a wrong owner and a
// nonexistent id are indistinguishable (both ErrNotFound) — existence isn't
// leaked across users.
func (s *Store) MarkNotificationRead(ctx context.Context, userID, id string) (*NotificationRecord, error) {
	result, err := s.db.ExecContext(ctx, `
		UPDATE in_app_notifications SET read = 1 WHERE id = ? AND user_id = ?
	`, id, userID)
	if err != nil {
		return nil, fmt.Errorf("mark notification read: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("rows affected: %w", err)
	}
	if rows == 0 {
		return nil, ErrNotFound
	}

	row := s.db.QueryRowContext(ctx, `
		SELECT id, user_id, alert_id, event_type, title, message, read, created_at
		FROM in_app_notifications WHERE id = ?
	`, id)
	n, err := scanNotification(row)
	if err != nil {
		return nil, fmt.Errorf("get notification: %w", err)
	}
	return &n, nil
}

// MarkAllNotificationsRead marks all of a user's unread notifications as read.
// Zero rows affected (nothing unread) is a valid outcome, not an error.
func (s *Store) MarkAllNotificationsRead(ctx context.Context, userID string) error {
	_, err := s.db.ExecContext(ctx, `
		UPDATE in_app_notifications SET read = 1 WHERE user_id = ? AND read = 0
	`, userID)
	if err != nil {
		return fmt.Errorf("mark all notifications read: %w", err)
	}
	return nil
}

// GetNotificationByID returns a notification by id, unscoped by user (for
// internal/background use like delivery retries).
func (s *Store) GetNotificationByID(ctx context.Context, id string) (*NotificationRecord, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT id, user_id, alert_id, event_type, title, message, read, created_at
		FROM in_app_notifications WHERE id = ?
	`, id)
	n, err := scanNotification(row)
	if err != nil {
		return nil, fmt.Errorf("get notification: %w", err)
	}
	return &n, nil
}

// rowScanner abstracts *sql.Row and *sql.Rows so scanNotification serves both
// the single-row (mark-read) and multi-row (list) query paths.
type rowScanner interface {
	Scan(dest ...any) error
}

func scanNotification(row rowScanner) (NotificationRecord, error) {
	var n NotificationRecord
	var alertID sql.NullString
	err := row.Scan(&n.ID, &n.UserID, &alertID, &n.EventType, &n.Title, &n.Message, &n.Read, &n.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return NotificationRecord{}, ErrNotFound
	}
	if err != nil {
		return NotificationRecord{}, err
	}
	n.AlertID = alertID.String
	return n, nil
}
