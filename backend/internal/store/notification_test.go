package store

import (
	"context"
	"errors"
	"testing"
)

func TestNotificationCreateAndList(t *testing.T) {
	ctx := context.Background()
	s := newDocTestStore(t)
	mustCreateUser(t, s, "user-1")
	mustCreateUser(t, s, "user-2")

	if err := s.CreateNotification(ctx, &NotificationRecord{
		UserID: "user-1", AlertID: "alert-1", EventType: "alert.created", Title: "Disk full",
	}); err != nil {
		t.Fatalf("CreateNotification() error: %v", err)
	}

	notifications, total, err := s.ListNotifications(ctx, "user-1", false, 0, 20)
	if err != nil {
		t.Fatalf("ListNotifications() error: %v", err)
	}
	if total != 1 || len(notifications) != 1 {
		t.Fatalf("ListNotifications() total=%d len=%d, want 1/1", total, len(notifications))
	}
	if notifications[0].AlertID != "alert-1" || notifications[0].Read {
		t.Fatalf("ListNotifications()[0] = %+v, want AlertID=alert-1 Read=false", notifications[0])
	}
}

func TestNotificationUnreadOnlyFilter(t *testing.T) {
	ctx := context.Background()
	s := newDocTestStore(t)
	mustCreateUser(t, s, "user-1")
	mustCreateUser(t, s, "user-2")

	if err := s.CreateNotification(ctx, &NotificationRecord{UserID: "user-1", EventType: "alert.created", Title: "A"}); err != nil {
		t.Fatalf("CreateNotification() error: %v", err)
	}
	if err := s.CreateNotification(ctx, &NotificationRecord{UserID: "user-1", EventType: "alert.created", Title: "B"}); err != nil {
		t.Fatalf("CreateNotification() error: %v", err)
	}

	all, total, err := s.ListNotifications(ctx, "user-1", false, 0, 20)
	if err != nil {
		t.Fatalf("ListNotifications() error: %v", err)
	}
	if total != 2 {
		t.Fatalf("total = %d, want 2", total)
	}

	if _, err := s.MarkNotificationRead(ctx, "user-1", all[0].ID); err != nil {
		t.Fatalf("MarkNotificationRead() error: %v", err)
	}

	unread, unreadTotal, err := s.ListNotifications(ctx, "user-1", true, 0, 20)
	if err != nil {
		t.Fatalf("ListNotifications(unreadOnly) error: %v", err)
	}
	if unreadTotal != 1 || len(unread) != 1 {
		t.Fatalf("unread total=%d len=%d, want 1/1", unreadTotal, len(unread))
	}
	if unread[0].ID != all[1].ID {
		t.Fatalf("unread notification id = %q, want %q", unread[0].ID, all[1].ID)
	}
}

func TestMarkNotificationReadIdempotentAndScoped(t *testing.T) {
	ctx := context.Background()
	s := newDocTestStore(t)
	mustCreateUser(t, s, "user-1")
	mustCreateUser(t, s, "user-2")

	if err := s.CreateNotification(ctx, &NotificationRecord{UserID: "user-1", EventType: "alert.created", Title: "A"}); err != nil {
		t.Fatalf("CreateNotification() error: %v", err)
	}
	list, _, err := s.ListNotifications(ctx, "user-1", false, 0, 20)
	if err != nil {
		t.Fatalf("ListNotifications() error: %v", err)
	}
	id := list[0].ID

	got, err := s.MarkNotificationRead(ctx, "user-1", id)
	if err != nil {
		t.Fatalf("MarkNotificationRead() error: %v", err)
	}
	if !got.Read {
		t.Fatalf("MarkNotificationRead() Read = false, want true")
	}

	// Idempotent: marking an already-read notification again still succeeds.
	if _, err := s.MarkNotificationRead(ctx, "user-1", id); err != nil {
		t.Fatalf("MarkNotificationRead() second call error: %v", err)
	}

	// Wrong owner: indistinguishable from a nonexistent id.
	if _, err := s.MarkNotificationRead(ctx, "user-2", id); !errors.Is(err, ErrNotFound) {
		t.Fatalf("MarkNotificationRead() wrong owner error = %v, want ErrNotFound", err)
	}
	if _, err := s.MarkNotificationRead(ctx, "user-1", "does-not-exist"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("MarkNotificationRead() missing id error = %v, want ErrNotFound", err)
	}
}

func TestMarkAllNotificationsRead(t *testing.T) {
	ctx := context.Background()
	s := newDocTestStore(t)
	mustCreateUser(t, s, "user-1")
	mustCreateUser(t, s, "user-2")

	for range 3 {
		if err := s.CreateNotification(ctx, &NotificationRecord{UserID: "user-1", EventType: "alert.created", Title: "A"}); err != nil {
			t.Fatalf("CreateNotification() error: %v", err)
		}
	}
	if err := s.CreateNotification(ctx, &NotificationRecord{UserID: "user-2", EventType: "alert.created", Title: "B"}); err != nil {
		t.Fatalf("CreateNotification() error: %v", err)
	}

	if err := s.MarkAllNotificationsRead(ctx, "user-1"); err != nil {
		t.Fatalf("MarkAllNotificationsRead() error: %v", err)
	}

	unread1, total1, err := s.ListNotifications(ctx, "user-1", true, 0, 20)
	if err != nil {
		t.Fatalf("ListNotifications() error: %v", err)
	}
	if total1 != 0 || len(unread1) != 0 {
		t.Fatalf("user-1 unread total=%d len=%d, want 0/0", total1, len(unread1))
	}

	unread2, total2, err := s.ListNotifications(ctx, "user-2", true, 0, 20)
	if err != nil {
		t.Fatalf("ListNotifications() error: %v", err)
	}
	if total2 != 1 || len(unread2) != 1 {
		t.Fatalf("user-2 unread total=%d len=%d, want 1/1 (must not be touched by user-1's mark-all)", total2, len(unread2))
	}

	// Nothing to mark is a valid outcome, not an error.
	if err := s.MarkAllNotificationsRead(ctx, "user-1"); err != nil {
		t.Fatalf("MarkAllNotificationsRead() on already-clear user error: %v", err)
	}
}
