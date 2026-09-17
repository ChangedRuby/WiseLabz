package store

import (
	"context"
	"errors"
	"testing"
)

func TestAPIKeyLifecycle(t *testing.T) {
	s := newDocTestStore(t)
	ctx := context.Background()
	user := &User{Username: "api-key-user"}
	if err := s.CreateUser(ctx, user); err != nil {
		t.Fatalf("CreateUser() error: %v", err)
	}

	key := &APIKey{UserID: user.ID, Name: "CI", TokenHash: HashToken("wlz_secret"), Role: "operator"}
	if err := s.CreateAPIKey(ctx, key); err != nil {
		t.Fatalf("CreateAPIKey() error: %v", err)
	}
	if key.ID == "" || key.CreatedAt == "" {
		t.Fatalf("CreateAPIKey() defaults = %#v", key)
	}

	found, err := s.GetAPIKeyByHash(ctx, key.TokenHash)
	if err != nil || found.ID != key.ID || found.TokenHash != key.TokenHash {
		t.Fatalf("GetAPIKeyByHash() = %#v, %v", found, err)
	}
	byID, err := s.GetAPIKeyByID(ctx, key.ID)
	if err != nil || byID.ID != key.ID {
		t.Fatalf("GetAPIKeyByID() = %#v, %v", byID, err)
	}

	keys, err := s.ListAPIKeysForUser(ctx, user.ID)
	if err != nil || len(keys) != 1 {
		t.Fatalf("ListAPIKeysForUser() = %#v, %v; want one", keys, err)
	}
	if err := s.TouchAPIKeyLastUsed(ctx, key.ID); err != nil {
		t.Fatalf("TouchAPIKeyLastUsed() error: %v", err)
	}
	found, _ = s.GetAPIKeyByID(ctx, key.ID)
	if found.LastUsedAt == "" {
		t.Error("TouchAPIKeyLastUsed() did not set lastUsedAt")
	}
	if err := s.RevokeAPIKey(ctx, key.ID); err != nil {
		t.Fatalf("RevokeAPIKey() error: %v", err)
	}
	found, _ = s.GetAPIKeyByID(ctx, key.ID)
	if found.RevokedAt == "" {
		t.Error("RevokeAPIKey() did not set revokedAt")
	}
}

// TestLookupAPIKeyRejectsDisabledUser is a regression test for GHSA-39m2:
// LookupAPIKey used to return the key's stored role in isolation from the
// users table, so a still-valid API key kept working after its owner was
// disabled or demoted. It must now reflect live user state.
func TestLookupAPIKeyRejectsDisabledUser(t *testing.T) {
	s := newDocTestStore(t)
	ctx := context.Background()

	user := &User{Username: "disabled-key-owner", InstanceAdminRole: "admin"}
	if err := s.CreateUser(ctx, user); err != nil {
		t.Fatalf("CreateUser() error: %v", err)
	}
	key := &APIKey{UserID: user.ID, Name: "CI", TokenHash: HashToken("wlz_secret"), Role: "operator"}
	if err := s.CreateAPIKey(ctx, key); err != nil {
		t.Fatalf("CreateAPIKey() error: %v", err)
	}

	// Key is valid and unrevoked before the user is disabled.
	claims, err := s.LookupAPIKey(ctx, key.TokenHash)
	if err != nil || !claims.InstanceAdmin {
		t.Fatalf("LookupAPIKey() before disable = %#v, %v; want InstanceAdmin true, no error", claims, err)
	}

	if err := s.UpdateUser(ctx, user.ID, map[string]any{"disabled": true}); err != nil {
		t.Fatalf("UpdateUser(disabled) error: %v", err)
	}

	// The API key row itself is untouched (still unrevoked), but the owning
	// user is now disabled, so the lookup must reject it.
	if _, err := s.LookupAPIKey(ctx, key.TokenHash); !errors.Is(err, ErrNotFound) {
		t.Fatalf("LookupAPIKey() after disable error = %v, want ErrNotFound", err)
	}
}

// TestLookupAPIKeyReflectsLiveRole is a regression test for GHSA-39m2: a
// role demotion after the key was issued must take effect immediately,
// instead of the key continuing to carry its original stored role.
func TestLookupAPIKeyReflectsLiveRole(t *testing.T) {
	s := newDocTestStore(t)
	ctx := context.Background()

	user := &User{Username: "demoted-key-owner", InstanceAdminRole: "admin"}
	if err := s.CreateUser(ctx, user); err != nil {
		t.Fatalf("CreateUser() error: %v", err)
	}
	key := &APIKey{UserID: user.ID, Name: "CI", TokenHash: HashToken("wlz_secret2"), Role: "operator"}
	if err := s.CreateAPIKey(ctx, key); err != nil {
		t.Fatalf("CreateAPIKey() error: %v", err)
	}

	if err := s.UpdateUser(ctx, user.ID, map[string]any{"instance_admin_role": "user"}); err != nil {
		t.Fatalf("UpdateUser(role) error: %v", err)
	}

	claims, err := s.LookupAPIKey(ctx, key.TokenHash)
	if err != nil || claims.InstanceAdmin {
		t.Fatalf("LookupAPIKey() after demotion = %#v, %v; want InstanceAdmin false, no error", claims, err)
	}
}

func TestAPIKeyNotFound(t *testing.T) {
	s := newDocTestStore(t)
	ctx := context.Background()
	if _, err := s.GetAPIKeyByHash(ctx, "missing"); !errors.Is(err, ErrNotFound) {
		t.Errorf("GetAPIKeyByHash() error = %v, want ErrNotFound", err)
	}
	if _, err := s.GetAPIKeyByID(ctx, "missing"); !errors.Is(err, ErrNotFound) {
		t.Errorf("GetAPIKeyByID() error = %v, want ErrNotFound", err)
	}
	if err := s.RevokeAPIKey(ctx, "missing"); !errors.Is(err, ErrNotFound) {
		t.Errorf("RevokeAPIKey() error = %v, want ErrNotFound", err)
	}
}

func TestRevokeAllAPIKeysForUser(t *testing.T) {
	s := newDocTestStore(t)
	ctx := context.Background()
	u := &User{Username: "revoke-all-user"}
	if err := s.CreateUser(ctx, u); err != nil {
		t.Fatalf("CreateUser() error: %v", err)
	}

	t.Run("revoke multiple keys", func(t *testing.T) {
		key1 := &APIKey{UserID: u.ID, Name: "key1", TokenHash: HashToken("secret1"), Role: "viewer"}
		key2 := &APIKey{UserID: u.ID, Name: "key2", TokenHash: HashToken("secret2"), Role: "operator"}
		if err := s.CreateAPIKey(ctx, key1); err != nil {
			t.Fatalf("CreateAPIKey(1) error: %v", err)
		}
		if err := s.CreateAPIKey(ctx, key2); err != nil {
			t.Fatalf("CreateAPIKey(2) error: %v", err)
		}

		// Revoke all
		if err := s.RevokeAllAPIKeysForUser(ctx, u.ID); err != nil {
			t.Fatalf("RevokeAllAPIKeysForUser() error: %v", err)
		}

		// Verify both are revoked
		k1, err := s.GetAPIKeyByID(ctx, key1.ID)
		if err != nil {
			t.Fatalf("GetAPIKeyByID(1) error: %v", err)
		}
		if k1.RevokedAt == "" {
			t.Error("key1 should be revoked")
		}
		k2, err := s.GetAPIKeyByID(ctx, key2.ID)
		if err != nil {
			t.Fatalf("GetAPIKeyByID(2) error: %v", err)
		}
		if k2.RevokedAt == "" {
			t.Error("key2 should be revoked")
		}
	})

	t.Run("no keys for user", func(t *testing.T) {
		u2 := &User{Username: "no-keys-user"}
		if err := s.CreateUser(ctx, u2); err != nil {
			t.Fatalf("CreateUser(u2) error: %v", err)
		}
		// Should not error even if no keys exist
		if err := s.RevokeAllAPIKeysForUser(ctx, u2.ID); err != nil {
			t.Fatalf("RevokeAllAPIKeysForUser(no keys) error: %v", err)
		}
	})

	t.Run("only revokes unrevoked keys", func(t *testing.T) {
		u3 := &User{Username: "mixed-keys-user"}
		if err := s.CreateUser(ctx, u3); err != nil {
			t.Fatalf("CreateUser(u3) error: %v", err)
		}

		key3 := &APIKey{UserID: u3.ID, Name: "active", TokenHash: HashToken("secret3"), Role: "viewer"}
		key4 := &APIKey{UserID: u3.ID, Name: "revoked", TokenHash: HashToken("secret4"), Role: "viewer"}
		if err := s.CreateAPIKey(ctx, key3); err != nil {
			t.Fatalf("CreateAPIKey(3) error: %v", err)
		}
		if err := s.CreateAPIKey(ctx, key4); err != nil {
			t.Fatalf("CreateAPIKey(4) error: %v", err)
		}

		// Manually revoke one
		if err := s.RevokeAPIKey(ctx, key4.ID); err != nil {
			t.Fatalf("RevokeAPIKey(key4) error: %v", err)
		}

		// RevokeAllAPIKeysForUser should only affect the active one
		if err := s.RevokeAllAPIKeysForUser(ctx, u3.ID); err != nil {
			t.Fatalf("RevokeAllAPIKeysForUser() error: %v", err)
		}

		// Both should now be revoked
		k3, _ := s.GetAPIKeyByID(ctx, key3.ID)
		if k3.RevokedAt == "" {
			t.Error("key3 should now be revoked")
		}
		k4, _ := s.GetAPIKeyByID(ctx, key4.ID)
		if k4.RevokedAt == "" {
			t.Error("key4 should still be revoked")
		}
	})
}

func TestTouchAPIKeyLastUsed(t *testing.T) {
	s := newDocTestStore(t)
	ctx := context.Background()
	u := &User{Username: "touch-user"}
	if err := s.CreateUser(ctx, u); err != nil {
		t.Fatalf("CreateUser() error: %v", err)
	}

	key := &APIKey{UserID: u.ID, Name: "test", TokenHash: HashToken("secret"), Role: "viewer"}
	if err := s.CreateAPIKey(ctx, key); err != nil {
		t.Fatalf("CreateAPIKey() error: %v", err)
	}

	t.Run("touch existing key", func(t *testing.T) {
		if err := s.TouchAPIKeyLastUsed(ctx, key.ID); err != nil {
			t.Fatalf("TouchAPIKeyLastUsed() error: %v", err)
		}
		k, _ := s.GetAPIKeyByID(ctx, key.ID)
		if k.LastUsedAt == "" {
			t.Error("LastUsedAt should be set")
		}
	})

	t.Run("nonexistent key", func(t *testing.T) {
		// TouchAPIKeyLastUsed doesn't return an error for nonexistent keys
		// it just silently does nothing (0 rows affected)
		if err := s.TouchAPIKeyLastUsed(ctx, "nonexistent"); err != nil {
			t.Fatalf("TouchAPIKeyLastUsed(nonexistent) error: %v", err)
		}
	})
}
