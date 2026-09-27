package store

import (
	"context"
	"errors"
	"testing"
)

func TestCreateAndListSavedViews(t *testing.T) {
	ctx := context.Background()
	s := newDocTestStore(t)

	if err := s.CreateSavedView(ctx, &SavedView{
		UserID: "user-1", Surface: "services", Name: "Down services", Filters: `{"status":"down"}`,
	}); err != nil {
		t.Fatalf("CreateSavedView() error: %v", err)
	}
	if err := s.CreateSavedView(ctx, &SavedView{
		UserID: "user-1", Surface: "services", Name: "All virtualization",
	}); err != nil {
		t.Fatalf("CreateSavedView() error: %v", err)
	}
	// Different user, same surface — must not leak into user-1's list.
	if err := s.CreateSavedView(ctx, &SavedView{
		UserID: "user-2", Surface: "services", Name: "Other user's view",
	}); err != nil {
		t.Fatalf("CreateSavedView() error: %v", err)
	}
	// Same user, different surface — must not leak into the services list.
	if err := s.CreateSavedView(ctx, &SavedView{
		UserID: "user-1", Surface: "changes", Name: "Critical changes",
	}); err != nil {
		t.Fatalf("CreateSavedView() error: %v", err)
	}

	views, err := s.ListSavedViews(ctx, "user-1", "services")
	if err != nil {
		t.Fatalf("ListSavedViews() error: %v", err)
	}
	if len(views) != 2 {
		t.Fatalf("ListSavedViews() = %d views, want 2; got %+v", len(views), views)
	}
	if views[0].ID == "" || views[0].CreatedAt == "" {
		t.Errorf("ListSavedViews()[0] missing generated ID/CreatedAt: %+v", views[0])
	}
	// Both views can share a created_at second, so order is not guaranteed;
	// look the unset-filters view up by name.
	for _, v := range views {
		if v.Name == "All virtualization" && v.Filters != "{}" {
			t.Errorf("Filters default = %q, want {} when unset", v.Filters)
		}
	}

	empty, err := s.ListSavedViews(ctx, "user-1", "alerts")
	if err != nil {
		t.Fatalf("ListSavedViews() error: %v", err)
	}
	if len(empty) != 0 {
		t.Fatalf("ListSavedViews(alerts) = %+v, want empty", empty)
	}
}

func TestDeleteSavedView(t *testing.T) {
	ctx := context.Background()
	s := newDocTestStore(t)

	v := &SavedView{UserID: "user-1", Surface: "changes", Name: "My view"}
	if err := s.CreateSavedView(ctx, v); err != nil {
		t.Fatalf("CreateSavedView() error: %v", err)
	}

	fetched, err := s.GetSavedView(ctx, v.ID)
	if err != nil {
		t.Fatalf("GetSavedView() error: %v", err)
	}
	if fetched.Name != "My view" {
		t.Errorf("GetSavedView().Name = %q, want %q", fetched.Name, "My view")
	}

	if err := s.DeleteSavedView(ctx, v.ID); err != nil {
		t.Fatalf("DeleteSavedView() error: %v", err)
	}

	if _, err := s.GetSavedView(ctx, v.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("GetSavedView() after delete = %v, want ErrNotFound", err)
	}
	if err := s.DeleteSavedView(ctx, v.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("DeleteSavedView() again = %v, want ErrNotFound", err)
	}
}
