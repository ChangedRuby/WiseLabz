package store

import (
	"context"
	"slices"
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

// TestGetUserConnectorRoleHighestAcrossSources covers #279 part 3: a manual
// viewer grant plus an oidc operator grant on the same connector must give
// operator, and deleting the manual row (the only thing DeleteConnectorGrant
// ever touches) falls back to the surviving oidc grant, never to no access.
func TestGetUserConnectorRoleHighestAcrossSources(t *testing.T) {
	s := newDocTestStore(t)
	ctx := context.Background()
	u := newTestUser(t, s, "multi-source-user")
	c := newTestConnector(t, s, "conn-multi-source")

	if _, err := s.UpsertConnectorGrant(ctx, u.ID, c.ID, "viewer"); err != nil {
		t.Fatalf("UpsertConnectorGrant() error: %v", err)
	}
	if _, err := upsertConnectorGrant(ctx, s.db, u.ID, c.ID, "operator", connectorGrantSourceOIDC); err != nil {
		t.Fatalf("upsertConnectorGrant(oidc) error: %v", err)
	}

	role, err := s.GetUserConnectorRole(ctx, u.ID, c.ID)
	if err != nil {
		t.Fatalf("GetUserConnectorRole() error: %v", err)
	}
	if role != "operator" {
		t.Fatalf("GetUserConnectorRole() = %q, want operator (highest of manual viewer + oidc operator)", role)
	}

	grants, err := s.ListUserConnectorGrants(ctx, u.ID)
	if err != nil {
		t.Fatalf("ListUserConnectorGrants() error: %v", err)
	}
	if len(grants) != 2 {
		t.Fatalf("ListUserConnectorGrants() = %d grants, want 2 (one per source)", len(grants))
	}

	// DeleteConnectorGrant (the admin API) only ever touches the manual row.
	if err := s.DeleteConnectorGrant(ctx, u.ID, c.ID); err != nil {
		t.Fatalf("DeleteConnectorGrant() error: %v", err)
	}
	role, err = s.GetUserConnectorRole(ctx, u.ID, c.ID)
	if err != nil {
		t.Fatalf("GetUserConnectorRole() after delete error: %v", err)
	}
	if role != "operator" {
		t.Fatalf("GetUserConnectorRole() after deleting the manual grant = %q, want operator (oidc grant survives)", role)
	}
}

// TestUpsertAndDeleteConnectorGrantNeverTouchOIDCRow covers the brief's
// "keep the manual grant API's behaviour the same" requirement: the admin
// PUT/DELETE endpoints must never create, modify or remove an oidc-sourced row.
func TestUpsertAndDeleteConnectorGrantNeverTouchOIDCRow(t *testing.T) {
	s := newDocTestStore(t)
	ctx := context.Background()
	u := newTestUser(t, s, "manual-api-user")
	c := newTestConnector(t, s, "conn-manual-api")

	oidcGrant, err := upsertConnectorGrant(ctx, s.db, u.ID, c.ID, "viewer", connectorGrantSourceOIDC)
	if err != nil {
		t.Fatalf("upsertConnectorGrant(oidc) error: %v", err)
	}

	if _, err := s.UpsertConnectorGrant(ctx, u.ID, c.ID, "operator"); err != nil {
		t.Fatalf("UpsertConnectorGrant() error: %v", err)
	}
	grants, err := s.ListUserConnectorGrants(ctx, u.ID)
	if err != nil {
		t.Fatalf("ListUserConnectorGrants() error: %v", err)
	}
	if len(grants) != 2 {
		t.Fatalf("ListUserConnectorGrants() = %d grants, want 2 (manual + untouched oidc)", len(grants))
	}
	for _, g := range grants {
		if g.Source == connectorGrantSourceOIDC && g.Role != oidcGrant.Role {
			t.Errorf("manual UpsertConnectorGrant() modified the oidc row: role = %q, want %q", g.Role, oidcGrant.Role)
		}
	}

	if err := s.DeleteConnectorGrant(ctx, u.ID, c.ID); err != nil {
		t.Fatalf("DeleteConnectorGrant() error: %v", err)
	}
	grants, err = s.ListUserConnectorGrants(ctx, u.ID)
	if err != nil {
		t.Fatalf("ListUserConnectorGrants() after delete error: %v", err)
	}
	if len(grants) != 1 || grants[0].Source != connectorGrantSourceOIDC {
		t.Fatalf("ListUserConnectorGrants() after manual delete = %v, want only the oidc row", grants)
	}
}

func TestSyncOIDCConnectorGrants(t *testing.T) {
	s := newDocTestStore(t)
	ctx := context.Background()
	u := newTestUser(t, s, "sync-user")
	c1 := newTestConnector(t, s, "conn-sync-1")
	c2 := newTestConnector(t, s, "conn-sync-2")

	t.Run("add", func(t *testing.T) {
		diff, err := s.SyncOIDCConnectorGrants(ctx, u.ID, map[string]string{c1.ID: "viewer", c2.ID: "operator"})
		if err != nil {
			t.Fatalf("SyncOIDCConnectorGrants() error: %v", err)
		}
		if len(diff.Added) != 2 || len(diff.Removed) != 0 {
			t.Fatalf("SyncOIDCConnectorGrants() diff = %+v, want 2 added, 0 removed", diff)
		}
		role, err := s.GetUserConnectorRole(ctx, u.ID, c2.ID)
		if err != nil {
			t.Fatalf("GetUserConnectorRole() error: %v", err)
		}
		if role != "operator" {
			t.Fatalf("GetUserConnectorRole() = %q, want operator", role)
		}
	})

	t.Run("noop when nothing changed", func(t *testing.T) {
		diff, err := s.SyncOIDCConnectorGrants(ctx, u.ID, map[string]string{c1.ID: "viewer", c2.ID: "operator"})
		if err != nil {
			t.Fatalf("SyncOIDCConnectorGrants() error: %v", err)
		}
		if !diff.Empty() {
			t.Fatalf("SyncOIDCConnectorGrants() diff = %+v, want empty (no change)", diff)
		}
	})

	t.Run("remove", func(t *testing.T) {
		diff, err := s.SyncOIDCConnectorGrants(ctx, u.ID, map[string]string{c1.ID: "viewer"})
		if err != nil {
			t.Fatalf("SyncOIDCConnectorGrants() error: %v", err)
		}
		if len(diff.Added) != 0 || len(diff.Removed) != 1 || diff.Removed[0].ConnectorID != c2.ID {
			t.Fatalf("SyncOIDCConnectorGrants() diff = %+v, want c2 removed only", diff)
		}
		role, err := s.GetUserConnectorRole(ctx, u.ID, c2.ID)
		if err != nil {
			t.Fatalf("GetUserConnectorRole() error: %v", err)
		}
		if role != "" {
			t.Fatalf("GetUserConnectorRole(c2) after removal = %q, want empty", role)
		}
	})

	t.Run("empty desired revokes every oidc grant", func(t *testing.T) {
		diff, err := s.SyncOIDCConnectorGrants(ctx, u.ID, map[string]string{})
		if err != nil {
			t.Fatalf("SyncOIDCConnectorGrants() error: %v", err)
		}
		if len(diff.Removed) != 1 {
			t.Fatalf("SyncOIDCConnectorGrants() diff = %+v, want the last oidc grant removed", diff)
		}
		grants, err := s.ListUserConnectorGrants(ctx, u.ID)
		if err != nil {
			t.Fatalf("ListUserConnectorGrants() error: %v", err)
		}
		if len(grants) != 0 {
			t.Fatalf("ListUserConnectorGrants() = %v, want none left", grants)
		}
	})

	t.Run("never touches a manual grant on the same connector", func(t *testing.T) {
		if _, err := s.UpsertConnectorGrant(ctx, u.ID, c1.ID, "viewer"); err != nil {
			t.Fatalf("UpsertConnectorGrant() error: %v", err)
		}
		if _, err := s.SyncOIDCConnectorGrants(ctx, u.ID, map[string]string{}); err != nil {
			t.Fatalf("SyncOIDCConnectorGrants() error: %v", err)
		}
		grants, err := s.ListUserConnectorGrants(ctx, u.ID)
		if err != nil {
			t.Fatalf("ListUserConnectorGrants() error: %v", err)
		}
		if len(grants) != 1 || grants[0].Source != connectorGrantSourceManual {
			t.Fatalf("SyncOIDCConnectorGrants() with empty desired = %v, want the manual grant untouched", grants)
		}
	})
}

func TestListConnectorIDs(t *testing.T) {
	s := newDocTestStore(t)
	ctx := context.Background()
	c1 := newTestConnector(t, s, "conn-list-ids-1")
	c2 := newTestConnector(t, s, "conn-list-ids-2")

	ids, err := s.ListConnectorIDs(ctx)
	if err != nil {
		t.Fatalf("ListConnectorIDs() error: %v", err)
	}
	found := map[string]bool{}
	for _, id := range ids {
		found[id] = true
	}
	if !found[c1.ID] || !found[c2.ID] {
		t.Fatalf("ListConnectorIDs() = %v, want it to include %q and %q", ids, c1.ID, c2.ID)
	}
}

func TestConnectorReaderIDs(t *testing.T) {
	s := newDocTestStore(t)
	ctx := context.Background()
	c := newTestConnector(t, s, "readers-a")
	other := newTestConnector(t, s, "readers-b")

	both := newTestUser(t, s, "both-sources")
	viewer := newTestUser(t, s, "viewer-only")
	disabled := newTestUser(t, s, "disabled")
	otherOnly := newTestUser(t, s, "other-only")
	admin := newTestUser(t, s, "grantless-admin")
	newTestUser(t, s, "no-grants")

	if err := s.UpdateUser(ctx, admin.ID, map[string]any{"instance_admin_role": "admin"}); err != nil {
		t.Fatalf("promote admin: %v", err)
	}
	if _, err := s.UpsertConnectorGrant(ctx, both.ID, c.ID, "operator"); err != nil {
		t.Fatal(err)
	}
	if _, err := s.SyncOIDCConnectorGrants(ctx, both.ID, map[string]string{c.ID: "viewer"}); err != nil {
		t.Fatalf("oidc grant: %v", err)
	}
	if _, err := s.UpsertConnectorGrant(ctx, viewer.ID, c.ID, "viewer"); err != nil {
		t.Fatal(err)
	}
	if _, err := s.UpsertConnectorGrant(ctx, disabled.ID, c.ID, "viewer"); err != nil {
		t.Fatal(err)
	}
	if err := s.UpdateUser(ctx, disabled.ID, map[string]any{"disabled": true}); err != nil {
		t.Fatalf("disable user: %v", err)
	}
	if _, err := s.UpsertConnectorGrant(ctx, otherOnly.ID, other.ID, "operator"); err != nil {
		t.Fatal(err)
	}

	got, err := s.ConnectorReaderIDs(ctx, c.ID)
	if err != nil {
		t.Fatalf("ConnectorReaderIDs() error: %v", err)
	}
	slices.Sort(got)
	want := []string{both.ID, viewer.ID}
	slices.Sort(want)
	if !slices.Equal(got, want) {
		t.Errorf("ConnectorReaderIDs() = %v, want %v", got, want)
	}
}
