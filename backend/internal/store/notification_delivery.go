package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// DeliveryStatus is the lifecycle state of one channel delivery attempt.
type DeliveryStatus string

// Delivery lifecycle states.
const (
	DeliveryStatusPending DeliveryStatus = "pending"
	DeliveryStatusSent    DeliveryStatus = "sent"
	DeliveryStatusFailed  DeliveryStatus = "failed"
)

// DeliveryRecord represents a row in the notification_deliveries table:
// one delivery attempt of one notification over one channel.
type DeliveryRecord struct {
	ID             string         `json:"id"`
	NotificationID string         `json:"notificationId"`
	Channel        string         `json:"channel"`
	Status         DeliveryStatus `json:"status"`
	Attempts       int            `json:"attempts"`
	LastError      string         `json:"lastError,omitempty"`
	NextAttemptAt  string         `json:"nextAttemptAt,omitempty"`
	CreatedAt      string         `json:"createdAt"`
	UpdatedAt      string         `json:"updatedAt"`
}

// CreateDelivery inserts a new delivery attempt row. A created delivery
// always represents at least one attempt, so Attempts defaults to 1 when unset.
func (s *Store) CreateDelivery(ctx context.Context, d *DeliveryRecord) error {
	if d.ID == "" {
		d.ID = uuid.New().String()
	}
	if d.Attempts == 0 {
		d.Attempts = 1
	}
	now := time.Now().UTC().Format(time.RFC3339)
	if d.CreatedAt == "" {
		d.CreatedAt = now
	}
	if d.UpdatedAt == "" {
		d.UpdatedAt = now
	}

	_, err := s.db.ExecContext(ctx, `
		INSERT INTO notification_deliveries (id, notification_id, channel, status, attempts, last_error, next_attempt_at, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, d.ID, d.NotificationID, d.Channel, d.Status, d.Attempts, d.LastError, nilToStr(d.NextAttemptAt), d.CreatedAt, d.UpdatedAt)
	if err != nil {
		return fmt.Errorf("create delivery: %w", err)
	}
	return nil
}

// UpdateDeliveryResult records the outcome of a delivery attempt.
func (s *Store) UpdateDeliveryResult(ctx context.Context, id string, status DeliveryStatus, attempts int, lastError, nextAttemptAt string) error {
	result, err := s.db.ExecContext(ctx, `
		UPDATE notification_deliveries
		SET status = ?, attempts = ?, last_error = ?, next_attempt_at = ?, updated_at = ?
		WHERE id = ?
	`, status, attempts, lastError, nilToStr(nextAttemptAt), time.Now().UTC().Format(time.RFC3339), id)
	if err != nil {
		return fmt.Errorf("update delivery result: %w", err)
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

// ListDueDeliveries returns failed deliveries whose next retry is due, ordered
// soonest-first. Never returns a nil slice.
func (s *Store) ListDueDeliveries(ctx context.Context, now string, limit int) ([]DeliveryRecord, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, notification_id, channel, status, attempts, last_error, next_attempt_at, created_at, updated_at
		FROM notification_deliveries
		WHERE status = 'failed' AND next_attempt_at IS NOT NULL AND next_attempt_at <= ?
		ORDER BY next_attempt_at ASC LIMIT ?
	`, now, limit)
	if err != nil {
		return nil, fmt.Errorf("list due deliveries: %w", err)
	}
	defer rows.Close() //nolint:errcheck

	deliveries := []DeliveryRecord{}
	for rows.Next() {
		d, err := scanDelivery(rows)
		if err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		deliveries = append(deliveries, d)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate due deliveries: %w", err)
	}
	return deliveries, nil
}

// ListDeliveries returns a paginated list of deliveries, most recent first,
// optionally filtered by status. Never returns a nil slice.
func (s *Store) ListDeliveries(ctx context.Context, status string, offset, limit int) ([]DeliveryRecord, int, error) {
	where := ""
	var args []any
	if status != "" {
		where = "WHERE status = ?"
		args = append(args, status)
	}

	return paginatedQuery(ctx, s.db, "notification_deliveries",
		"id, notification_id, channel, status, attempts, last_error, next_attempt_at, created_at, updated_at",
		where, args, "created_at DESC", limit, offset, scanDelivery)
}

func scanDelivery(row rowScanner) (DeliveryRecord, error) {
	var d DeliveryRecord
	var lastError sql.NullString
	var nextAttemptAt sql.NullString
	err := row.Scan(&d.ID, &d.NotificationID, &d.Channel, &d.Status, &d.Attempts, &lastError, &nextAttemptAt, &d.CreatedAt, &d.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return DeliveryRecord{}, ErrNotFound
	}
	if err != nil {
		return DeliveryRecord{}, err
	}
	d.LastError = nullStrToStr(lastError)
	d.NextAttemptAt = nullStrToStr(nextAttemptAt)
	return d, nil
}
