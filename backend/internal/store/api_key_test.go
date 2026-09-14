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

	user := &User{Username: "disabled-key-owner", Role: "operator"}
	if err := s.CreateUser(ctx, user); err != nil {
		t.Fatalf("CreateUser() error: %v", err)
	}
	key := &APIKey{UserID: user.ID, Name: "CI", TokenHash: HashToken("wlz_secret"), Role: "operator"}
	if err := s.CreateAPIKey(ctx, key); err != nil {
		t.Fatalf("CreateAPIKey() error: %v", err)
	}

	// Key is valid and unrevoked before the user is disabled.
	claims, err := s.LookupAPIKey(ctx, key.TokenHash)
	if err != nil || claims.Role != "operator" {
		t.Fatalf("LookupAPIKey() before disable = %#v, %v; want role operator, no error", claims, err)
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

	user := &User{Username: "demoted-key-owner", Role: "operator"}
	if err := s.CreateUser(ctx, user); err != nil {
		t.Fatalf("CreateUser() error: %v", err)
	}
	key := &APIKey{UserID: user.ID, Name: "CI", TokenHash: HashToken("wlz_secret2"), Role: "operator"}
	if err := s.CreateAPIKey(ctx, key); err != nil {
		t.Fatalf("CreateAPIKey() error: %v", err)
	}

	if err := s.UpdateUser(ctx, user.ID, map[string]any{"role": "viewer"}); err != nil {
		t.Fatalf("UpdateUser(role) error: %v", err)
	}

	claims, err := s.LookupAPIKey(ctx, key.TokenHash)
	if err != nil || claims.Role != "viewer" {
		t.Fatalf("LookupAPIKey() after demotion = %#v, %v; want role viewer, no error", claims, err)
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
