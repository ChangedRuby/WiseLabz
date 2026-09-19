package store

import (
	"context"
	"errors"
	"testing"
	"time"
)

func createTestNotification(ctx context.Context, t *testing.T, s *Store) string {
	t.Helper()
	mustCreateUser(t, s, "user-1")
	n := &NotificationRecord{UserID: "user-1", EventType: "alert.created", Title: "Disk full"}
	if err := s.CreateNotification(ctx, n); err != nil {
		t.Fatalf("CreateNotification() error: %v", err)
	}
	list, _, err := s.ListNotifications(ctx, "user-1", false, 0, 20)
	if err != nil {
		t.Fatalf("ListNotifications() error: %v", err)
	}
	return list[0].ID
}

func TestDeliveryCreateAndList(t *testing.T) {
	ctx := context.Background()
	s := newDocTestStore(t)
	notificationID := createTestNotification(ctx, t, s)

	if err := s.CreateDelivery(ctx, &DeliveryRecord{
		NotificationID: notificationID, Channel: "webhook", Status: DeliveryStatusPending,
	}); err != nil {
		t.Fatalf("CreateDelivery() error: %v", err)
	}

	deliveries, total, err := s.ListDeliveries(ctx, "", 0, 20)
	if err != nil {
		t.Fatalf("ListDeliveries() error: %v", err)
	}
	if total != 1 || len(deliveries) != 1 {
		t.Fatalf("ListDeliveries() total=%d len=%d, want 1/1", total, len(deliveries))
	}
	d := deliveries[0]
	if d.NotificationID != notificationID || d.Channel != "webhook" || d.Status != DeliveryStatusPending {
		t.Fatalf("ListDeliveries()[0] = %+v, want notificationID=%s channel=webhook status=pending", d, notificationID)
	}
	if d.Attempts != 1 {
		t.Fatalf("ListDeliveries()[0].Attempts = %d, want 1 (default)", d.Attempts)
	}
}

func TestUpdateDeliveryResultTransitionsAndClearsNextAttempt(t *testing.T) {
	ctx := context.Background()
	s := newDocTestStore(t)
	notificationID := createTestNotification(ctx, t, s)

	d := &DeliveryRecord{NotificationID: notificationID, Channel: "smtp", Status: DeliveryStatusPending}
	if err := s.CreateDelivery(ctx, d); err != nil {
		t.Fatalf("CreateDelivery() error: %v", err)
	}

	future := time.Now().UTC().Add(time.Hour).Format(time.RFC3339)
	if err := s.UpdateDeliveryResult(ctx, d.ID, DeliveryStatusFailed, 2, "smtp timeout", future); err != nil {
		t.Fatalf("UpdateDeliveryResult() error: %v", err)
	}

	deliveries, _, err := s.ListDeliveries(ctx, "", 0, 20)
	if err != nil {
		t.Fatalf("ListDeliveries() error: %v", err)
	}
	got := deliveries[0]
	if got.Status != DeliveryStatusFailed || got.Attempts != 2 || got.LastError != "smtp timeout" || got.NextAttemptAt != future {
		t.Fatalf("after failed update = %+v, want status=failed attempts=2 lastError=smtp timeout nextAttemptAt=%s", got, future)
	}

	if err := s.UpdateDeliveryResult(ctx, d.ID, DeliveryStatusSent, 2, "", ""); err != nil {
		t.Fatalf("UpdateDeliveryResult() error: %v", err)
	}

	deliveries, _, err = s.ListDeliveries(ctx, "", 0, 20)
	if err != nil {
		t.Fatalf("ListDeliveries() error: %v", err)
	}
	got = deliveries[0]
	if got.Status != DeliveryStatusSent || got.NextAttemptAt != "" {
		t.Fatalf("after sent update = %+v, want status=sent nextAttemptAt=\"\" (NULL, not literal empty string)", got)
	}
}

func TestUpdateDeliveryResultNotFound(t *testing.T) {
	ctx := context.Background()
	s := newDocTestStore(t)

	err := s.UpdateDeliveryResult(ctx, "does-not-exist", DeliveryStatusSent, 1, "", "")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("UpdateDeliveryResult() error = %v, want ErrNotFound", err)
	}
}

func TestListDueDeliveries(t *testing.T) {
	ctx := context.Background()
	s := newDocTestStore(t)
	notificationID := createTestNotification(ctx, t, s)

	now := time.Now().UTC()
	past := now.Add(-time.Hour).Format(time.RFC3339)
	future := now.Add(time.Hour).Format(time.RFC3339)
	nowStr := now.Format(time.RFC3339)

	// Due: failed, next_attempt_at in the past.
	due := &DeliveryRecord{NotificationID: notificationID, Channel: "webhook", Status: DeliveryStatusFailed, NextAttemptAt: past}
	if err := s.CreateDelivery(ctx, due); err != nil {
		t.Fatalf("CreateDelivery(due) error: %v", err)
	}
	// Not due: failed, but next_attempt_at in the future.
	notDueFuture := &DeliveryRecord{NotificationID: notificationID, Channel: "webhook", Status: DeliveryStatusFailed, NextAttemptAt: future}
	if err := s.CreateDelivery(ctx, notDueFuture); err != nil {
		t.Fatalf("CreateDelivery(notDueFuture) error: %v", err)
	}
	// Not due: failed, but next_attempt_at NULL.
	notDueNull := &DeliveryRecord{NotificationID: notificationID, Channel: "webhook", Status: DeliveryStatusFailed}
	if err := s.CreateDelivery(ctx, notDueNull); err != nil {
		t.Fatalf("CreateDelivery(notDueNull) error: %v", err)
	}
	// Not due: pending.
	pending := &DeliveryRecord{NotificationID: notificationID, Channel: "webhook", Status: DeliveryStatusPending, NextAttemptAt: past}
	if err := s.CreateDelivery(ctx, pending); err != nil {
		t.Fatalf("CreateDelivery(pending) error: %v", err)
	}
	// Not due: sent.
	sent := &DeliveryRecord{NotificationID: notificationID, Channel: "webhook", Status: DeliveryStatusSent, NextAttemptAt: past}
	if err := s.CreateDelivery(ctx, sent); err != nil {
		t.Fatalf("CreateDelivery(sent) error: %v", err)
	}

	dueList, err := s.ListDueDeliveries(ctx, nowStr, 20)
	if err != nil {
		t.Fatalf("ListDueDeliveries() error: %v", err)
	}
	if len(dueList) != 1 || dueList[0].ID != due.ID {
		t.Fatalf("ListDueDeliveries() = %+v, want only %s", dueList, due.ID)
	}
}

func TestListDeliveriesStatusFilterAndPagination(t *testing.T) {
	ctx := context.Background()
	s := newDocTestStore(t)
	notificationID := createTestNotification(ctx, t, s)

	for range 2 {
		if err := s.CreateDelivery(ctx, &DeliveryRecord{NotificationID: notificationID, Channel: "webhook", Status: DeliveryStatusFailed}); err != nil {
			t.Fatalf("CreateDelivery(failed) error: %v", err)
		}
	}
	if err := s.CreateDelivery(ctx, &DeliveryRecord{NotificationID: notificationID, Channel: "in_app", Status: DeliveryStatusSent}); err != nil {
		t.Fatalf("CreateDelivery(sent) error: %v", err)
	}

	failed, failedTotal, err := s.ListDeliveries(ctx, "failed", 0, 20)
	if err != nil {
		t.Fatalf("ListDeliveries(failed) error: %v", err)
	}
	if failedTotal != 2 || len(failed) != 2 {
		t.Fatalf("ListDeliveries(failed) total=%d len=%d, want 2/2", failedTotal, len(failed))
	}
	for _, d := range failed {
		if d.Status != DeliveryStatusFailed {
			t.Fatalf("ListDeliveries(failed) contained status=%s", d.Status)
		}
	}

	all, allTotal, err := s.ListDeliveries(ctx, "", 0, 20)
	if err != nil {
		t.Fatalf("ListDeliveries(\"\") error: %v", err)
	}
	if allTotal != 3 || len(all) != 3 {
		t.Fatalf("ListDeliveries(\"\") total=%d len=%d, want 3/3", allTotal, len(all))
	}

	page, pageTotal, err := s.ListDeliveries(ctx, "", 0, 1)
	if err != nil {
		t.Fatalf("ListDeliveries(paginated) error: %v", err)
	}
	if pageTotal != 3 || len(page) != 1 {
		t.Fatalf("ListDeliveries(paginated) total=%d len=%d, want 3/1", pageTotal, len(page))
	}
}
