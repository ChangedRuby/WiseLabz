package store

import (
	"context"
	"testing"
)

func newTestConnector(t *testing.T, s *Store, name string) *ConnectorRecord {
	t.Helper()
	c := &ConnectorRecord{Name: name, Category: "virtualization", Type: "proxmox", URL: "https://example.com"}
	if err := s.CreateConnector(context.Background(), c); err != nil {
		t.Fatalf("CreateConnector() error: %v", err)
	}
	return c
}

func newTestUser(t *testing.T, s *Store, username string) *User {
	t.Helper()
	u := &User{Username: username}
	if err := s.CreateUser(context.Background(), u); err != nil {
		t.Fatalf("CreateUser() error: %v", err)
	}
	return u
}

func TestUserHasConnectorRoleDefaultDeny(t *testing.T) {
	s := newDocTestStore(t)
	ctx := context.Background()
	u := newTestUser(t, s, "no-grant-user")
	c := newTestConnector(t, s, "conn-a")

	role, err := s.GetUserConnectorRole(ctx, u.ID, c.ID)
	if err != nil {
		t.Fatalf("GetUserConnectorRole() error: %v", err)
	}
	if role != "" {
		t.Errorf("GetUserConnectorRole() = %q, want empty (no grant)", role)
	}

	ok, err := s.UserHasConnectorRole(ctx, u.ID, c.ID, "viewer")
	if err != nil {
		t.Fatalf("UserHasConnectorRole() error: %v", err)
	}
	if ok {
		t.Error("UserHasConnectorRole() = true, want false with no grant")
	}
}

func TestUserHasConnectorRoleHierarchy(t *testing.T) {
	s := newDocTestStore(t)
	ctx := context.Background()
	u := newTestUser(t, s, "hierarchy-user")
	c := newTestConnector(t, s, "conn-b")

	if _, err := s.UpsertConnectorGrant(ctx, u.ID, c.ID, "operator"); err != nil {
		t.Fatalf("UpsertConnectorGrant() error: %v", err)
	}

	for _, minRole := range []string{"viewer", "operator"} {
		ok, err := s.UserHasConnectorRole(ctx, u.ID, c.ID, minRole)
		if err != nil {
			t.Fatalf("UserHasConnectorRole(%s) error: %v", minRole, err)
		}
		if !ok {
			t.Errorf("UserHasConnectorRole(%s) = false, want true for an operator grant", minRole)
		}
	}
}

func TestUserHasConnectorRoleViewerInsufficientForOperator(t *testing.T) {
	s := newDocTestStore(t)
	ctx := context.Background()
	u := newTestUser(t, s, "viewer-only-user")
	c := newTestConnector(t, s, "conn-c")

	if _, err := s.UpsertConnectorGrant(ctx, u.ID, c.ID, "viewer"); err != nil {
		t.Fatalf("UpsertConnectorGrant() error: %v", err)
	}

	ok, err := s.UserHasConnectorRole(ctx, u.ID, c.ID, "operator")
	if err != nil {
		t.Fatalf("UserHasConnectorRole() error: %v", err)
	}
	if ok {
		t.Error("UserHasConnectorRole(operator) = true, want false for a viewer-only grant")
	}
}

func TestUpsertConnectorGrantUpdatesExistingRole(t *testing.T) {
	s := newDocTestStore(t)
	ctx := context.Background()
	u := newTestUser(t, s, "upsert-user")
	c := newTestConnector(t, s, "conn-d")

	first, err := s.UpsertConnectorGrant(ctx, u.ID, c.ID, "viewer")
	if err != nil {
		t.Fatalf("UpsertConnectorGrant() error: %v", err)
	}
	second, err := s.UpsertConnectorGrant(ctx, u.ID, c.ID, "operator")
	if err != nil {
		t.Fatalf("UpsertConnectorGrant() update error: %v", err)
	}
	if second.ID != first.ID {
		t.Errorf("UpsertConnectorGrant() update created a new row: first ID %q, second ID %q", first.ID, second.ID)
	}
	if second.Role != "operator" {
		t.Errorf("UpsertConnectorGrant() update role = %q, want operator", second.Role)
	}

	grants, err := s.ListUserConnectorGrants(ctx, u.ID)
	if err != nil {
		t.Fatalf("ListUserConnectorGrants() error: %v", err)
	}
	if len(grants) != 1 {
		t.Fatalf("ListUserConnectorGrants() = %d grants, want 1 (update, not insert)", len(grants))
	}
}

func TestDeleteConnectorGrant(t *testing.T) {
	s := newDocTestStore(t)
	ctx := context.Background()
	u := newTestUser(t, s, "delete-grant-user")
	c := newTestConnector(t, s, "conn-e")

	if _, err := s.UpsertConnectorGrant(ctx, u.ID, c.ID, "viewer"); err != nil {
		t.Fatalf("UpsertConnectorGrant() error: %v", err)
	}
	if err := s.DeleteConnectorGrant(ctx, u.ID, c.ID); err != nil {
		t.Fatalf("DeleteConnectorGrant() error: %v", err)
	}
	role, err := s.GetUserConnectorRole(ctx, u.ID, c.ID)
	if err != nil {
		t.Fatalf("GetUserConnectorRole() error: %v", err)
	}
	if role != "" {
		t.Errorf("GetUserConnectorRole() after delete = %q, want empty", role)
	}
	if err := s.DeleteConnectorGrant(ctx, u.ID, c.ID); err != ErrNotFound {
		t.Errorf("DeleteConnectorGrant() on missing grant = %v, want ErrNotFound", err)
	}
}

func TestListConnectorGrants(t *testing.T) {
	s := newDocTestStore(t)
	ctx := context.Background()
	c := newTestConnector(t, s, "conn-f")
	u1 := newTestUser(t, s, "grant-list-user-1")
	u2 := newTestUser(t, s, "grant-list-user-2")

	if _, err := s.UpsertConnectorGrant(ctx, u1.ID, c.ID, "viewer"); err != nil {
		t.Fatalf("UpsertConnectorGrant() error: %v", err)
	}
	if _, err := s.UpsertConnectorGrant(ctx, u2.ID, c.ID, "operator"); err != nil {
		t.Fatalf("UpsertConnectorGrant() error: %v", err)
	}

	grants, err := s.ListConnectorGrants(ctx, c.ID)
	if err != nil {
		t.Fatalf("ListConnectorGrants() error: %v", err)
	}
	if len(grants) != 2 {
		t.Fatalf("ListConnectorGrants() = %d grants, want 2", len(grants))
	}
}

func TestFilterConnectorIDsByGrant(t *testing.T) {
	s := newDocTestStore(t)
	ctx := context.Background()
	u := newTestUser(t, s, "filter-user")
	viewerConn := newTestConnector(t, s, "conn-viewer")
	operatorConn := newTestConnector(t, s, "conn-operator")
	noGrantConn := newTestConnector(t, s, "conn-none")

	if _, err := s.UpsertConnectorGrant(ctx, u.ID, viewerConn.ID, "viewer"); err != nil {
		t.Fatalf("UpsertConnectorGrant() error: %v", err)
	}
	if _, err := s.UpsertConnectorGrant(ctx, u.ID, operatorConn.ID, "operator"); err != nil {
		t.Fatalf("UpsertConnectorGrant() error: %v", err)
	}

	t.Run("empty ids returns empty", func(t *testing.T) {
		got, err := s.FilterConnectorIDsByGrant(ctx, u.ID, nil, "viewer")
		if err != nil {
			t.Fatalf("FilterConnectorIDsByGrant() error: %v", err)
		}
		if len(got) != 0 {
			t.Errorf("FilterConnectorIDsByGrant(nil) = %v, want empty", got)
		}
	})

	t.Run("viewer minimum keeps both grants, drops ungranted", func(t *testing.T) {
		got, err := s.FilterConnectorIDsByGrant(ctx, u.ID, []string{viewerConn.ID, operatorConn.ID, noGrantConn.ID}, "viewer")
		if err != nil {
			t.Fatalf("FilterConnectorIDsByGrant() error: %v", err)
		}
		if len(got) != 2 {
			t.Fatalf("FilterConnectorIDsByGrant(viewer) = %v, want 2 ids", got)
		}
	})

	t.Run("operator minimum keeps only the operator grant", func(t *testing.T) {
		got, err := s.FilterConnectorIDsByGrant(ctx, u.ID, []string{viewerConn.ID, operatorConn.ID}, "operator")
		if err != nil {
			t.Fatalf("FilterConnectorIDsByGrant() error: %v", err)
		}
		if len(got) != 1 || got[0] != operatorConn.ID {
			t.Errorf("FilterConnectorIDsByGrant(operator) = %v, want only %q", got, operatorConn.ID)
		}
	})

	t.Run("user with no grants at all", func(t *testing.T) {
		other := newTestUser(t, s, "no-grants-at-all")
		got, err := s.FilterConnectorIDsByGrant(ctx, other.ID, []string{viewerConn.ID, operatorConn.ID}, "viewer")
		if err != nil {
			t.Fatalf("FilterConnectorIDsByGrant() error: %v", err)
		}
		if len(got) != 0 {
			t.Errorf("FilterConnectorIDsByGrant() for ungranted user = %v, want empty", got)
		}
	})
}
