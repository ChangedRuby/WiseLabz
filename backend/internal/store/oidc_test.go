package store

import (
	"context"
	"testing"
)

func TestOIDCProviderFlags(t *testing.T) {
	s := newDocTestStore(t)
	ctx := context.Background()
	if err := s.SetOIDCProviderEnabled(ctx, "authentik", false); err != nil {
		t.Fatalf("SetOIDCProviderEnabled(false) error: %v", err)
	}
	flags, err := s.GetOIDCProviderFlags(ctx)
	if err != nil || flags["authentik"] {
		t.Fatalf("GetOIDCProviderFlags() = %v, %v; want authentik=false", flags, err)
	}
	if err := s.SetOIDCProviderEnabled(ctx, "authentik", true); err != nil {
		t.Fatalf("SetOIDCProviderEnabled(true) error: %v", err)
	}
	flags, err = s.GetOIDCProviderFlags(ctx)
	if err != nil || !flags["authentik"] {
		t.Fatalf("GetOIDCProviderFlags() = %v, %v; want authentik=true", flags, err)
	}
}

func TestSetOIDCProviderEnabledErrorHandling(t *testing.T) {
	s := newDocTestStore(t)
	ctx := context.Background()

	t.Run("upsert multiple times", func(t *testing.T) {
		if err := s.SetOIDCProviderEnabled(ctx, "okta", true); err != nil {
			t.Fatalf("SetOIDCProviderEnabled(true) error: %v", err)
		}
		if err := s.SetOIDCProviderEnabled(ctx, "okta", false); err != nil {
			t.Fatalf("SetOIDCProviderEnabled(false) error: %v", err)
		}
		flags, err := s.GetOIDCProviderFlags(ctx)
		if err != nil || flags["okta"] {
			t.Fatalf("GetOIDCProviderFlags() = %v, %v; want okta=false", flags, err)
		}
	})

	t.Run("closed store", func(t *testing.T) {
		closedStore := newDocTestStore(t)
		if err := closedStore.Close(); err != nil {
			t.Fatalf("Close() error: %v", err)
		}
		if err := closedStore.SetOIDCProviderEnabled(ctx, "test", true); err == nil {
			t.Fatal("SetOIDCProviderEnabled() on closed store should error")
		}
	})
}

func TestGetOIDCProviderFlagsErrorHandling(t *testing.T) {
	s := newDocTestStore(t)
	ctx := context.Background()

	t.Run("happy path", func(t *testing.T) {
		flags, err := s.GetOIDCProviderFlags(ctx)
		if err != nil {
			t.Fatalf("GetOIDCProviderFlags() error: %v", err)
		}
		if flags == nil {
			t.Fatal("GetOIDCProviderFlags() returned nil map, want empty map")
		}
	})

	t.Run("closed store", func(t *testing.T) {
		closedStore := newDocTestStore(t)
		if err := closedStore.Close(); err != nil {
			t.Fatalf("Close() error: %v", err)
		}
		_, err := closedStore.GetOIDCProviderFlags(ctx)
		if err == nil {
			t.Fatal("GetOIDCProviderFlags() on closed store should error")
		}
	})
}
