package store

import (
	"context"
	"errors"
	"testing"
)

func mfaTestUser(t *testing.T, s *Store) *User {
	t.Helper()
	u := &User{Username: "mfa-user-" + t.Name()}
	if err := s.CreateUser(context.Background(), u); err != nil {
		t.Fatalf("CreateUser() error: %v", err)
	}
	return u
}

func TestConsumeTOTPStepReplayGuard(t *testing.T) {
	s := newDocTestStore(t)
	ctx := context.Background()
	u := mfaTestUser(t, s)

	factor, err := s.CreatePendingTOTP(ctx, u.ID, "Authenticator", "encrypted-secret")
	if err != nil {
		t.Fatalf("CreatePendingTOTP() error: %v", err)
	}

	ok, err := s.ConsumeTOTPStep(ctx, factor.ID, 100)
	if err != nil || !ok {
		t.Fatalf("ConsumeTOTPStep(100) = %v, %v; want true, nil", ok, err)
	}

	// The same step again must be rejected (replay).
	ok, err = s.ConsumeTOTPStep(ctx, factor.ID, 100)
	if err != nil || ok {
		t.Fatalf("ConsumeTOTPStep(100) replay = %v, %v; want false, nil", ok, err)
	}

	// An older step must also be rejected.
	ok, err = s.ConsumeTOTPStep(ctx, factor.ID, 99)
	if err != nil || ok {
		t.Fatalf("ConsumeTOTPStep(99) = %v, %v; want false, nil", ok, err)
	}

	// A newer step is accepted.
	ok, err = s.ConsumeTOTPStep(ctx, factor.ID, 101)
	if err != nil || !ok {
		t.Fatalf("ConsumeTOTPStep(101) = %v, %v; want true, nil", ok, err)
	}
}

func TestConfirmFactorRejectsSecondConfirmedTOTP(t *testing.T) {
	s := newDocTestStore(t)
	ctx := context.Background()
	u := mfaTestUser(t, s)

	first, err := s.CreatePendingTOTP(ctx, u.ID, "First", "secret-1")
	if err != nil {
		t.Fatalf("CreatePendingTOTP() error: %v", err)
	}
	if _, err := s.ConfirmFactor(ctx, first.ID); err != nil {
		t.Fatalf("ConfirmFactor() first error: %v", err)
	}

	second, err := s.CreatePendingTOTP(ctx, u.ID, "Second", "secret-2")
	if err != nil {
		t.Fatalf("CreatePendingTOTP() error: %v", err)
	}
	if _, err := s.ConfirmFactor(ctx, second.ID); !errors.Is(err, ErrConflict) {
		t.Fatalf("ConfirmFactor() second = %v; want ErrConflict", err)
	}

	// Confirming the same factor twice is also rejected.
	if _, err := s.ConfirmFactor(ctx, first.ID); !errors.Is(err, ErrConflict) {
		t.Fatalf("ConfirmFactor() re-confirm = %v; want ErrConflict", err)
	}
}

func TestUserHasMFAOnlyCountsConfirmedFactors(t *testing.T) {
	s := newDocTestStore(t)
	ctx := context.Background()
	u := mfaTestUser(t, s)

	hasMFA, err := s.UserHasMFA(ctx, u.ID)
	if err != nil || hasMFA {
		t.Fatalf("UserHasMFA() before enrollment = %v, %v; want false, nil", hasMFA, err)
	}

	factor, err := s.CreatePendingTOTP(ctx, u.ID, "Pending", "secret")
	if err != nil {
		t.Fatalf("CreatePendingTOTP() error: %v", err)
	}
	hasMFA, err = s.UserHasMFA(ctx, u.ID)
	if err != nil || hasMFA {
		t.Fatalf("UserHasMFA() with pending factor = %v, %v; want false, nil", hasMFA, err)
	}

	if _, err := s.ConfirmFactor(ctx, factor.ID); err != nil {
		t.Fatalf("ConfirmFactor() error: %v", err)
	}
	hasMFA, err = s.UserHasMFA(ctx, u.ID)
	if err != nil || !hasMFA {
		t.Fatalf("UserHasMFA() after confirm = %v, %v; want true, nil", hasMFA, err)
	}

	factors, err := s.ListUserFactors(ctx, u.ID)
	if err != nil || len(factors) != 1 {
		t.Fatalf("ListUserFactors() = %v, %v; want 1 factor", factors, err)
	}
}

func TestRecoveryCodesSingleUse(t *testing.T) {
	s := newDocTestStore(t)
	ctx := context.Background()
	u := mfaTestUser(t, s)

	if err := s.ReplaceRecoveryCodes(ctx, u.ID, []string{"hash-a", "hash-b", "hash-c"}); err != nil {
		t.Fatalf("ReplaceRecoveryCodes() error: %v", err)
	}

	remaining, err := s.CountRecoveryCodesRemaining(ctx, u.ID)
	if err != nil || remaining != 3 {
		t.Fatalf("CountRecoveryCodesRemaining() = %d, %v; want 3, nil", remaining, err)
	}

	ok, err := s.ConsumeRecoveryCode(ctx, u.ID, "hash-a")
	if err != nil || !ok {
		t.Fatalf("ConsumeRecoveryCode(hash-a) = %v, %v; want true, nil", ok, err)
	}

	// Reusing the same code must fail.
	ok, err = s.ConsumeRecoveryCode(ctx, u.ID, "hash-a")
	if err != nil || ok {
		t.Fatalf("ConsumeRecoveryCode(hash-a) reuse = %v, %v; want false, nil", ok, err)
	}

	// A code that never existed must also fail.
	ok, err = s.ConsumeRecoveryCode(ctx, u.ID, "hash-z")
	if err != nil || ok {
		t.Fatalf("ConsumeRecoveryCode(hash-z) = %v, %v; want false, nil", ok, err)
	}

	remaining, err = s.CountRecoveryCodesRemaining(ctx, u.ID)
	if err != nil || remaining != 2 {
		t.Fatalf("CountRecoveryCodesRemaining() after use = %d, %v; want 2, nil", remaining, err)
	}

	// Replacing wipes the old codes, including unused ones.
	if err := s.ReplaceRecoveryCodes(ctx, u.ID, []string{"hash-d"}); err != nil {
		t.Fatalf("ReplaceRecoveryCodes() second error: %v", err)
	}
	ok, err = s.ConsumeRecoveryCode(ctx, u.ID, "hash-b")
	if err != nil || ok {
		t.Fatalf("ConsumeRecoveryCode(hash-b) after replace = %v, %v; want false, nil", ok, err)
	}
	remaining, err = s.CountRecoveryCodesRemaining(ctx, u.ID)
	if err != nil || remaining != 1 {
		t.Fatalf("CountRecoveryCodesRemaining() after replace = %d, %v; want 1, nil", remaining, err)
	}
}

func TestDeleteUserCascadesMFAFactorsAndRecoveryCodes(t *testing.T) {
	// newDocTestStore opens SQLite without the foreign_keys pragma; use
	// OpenDB (production's own path) so ON DELETE CASCADE actually fires.
	db, err := OpenDB("sqlite", "file:"+migratedSQLite(t))
	if err != nil {
		t.Fatalf("OpenDB() error: %v", err)
	}
	defer db.Close() //nolint:errcheck
	s := New(db, "sqlite")
	ctx := context.Background()
	u := mfaTestUser(t, s)

	factor, err := s.CreatePendingTOTP(ctx, u.ID, "Authenticator", "secret")
	if err != nil {
		t.Fatalf("CreatePendingTOTP() error: %v", err)
	}
	if _, err := s.ConfirmFactor(ctx, factor.ID); err != nil {
		t.Fatalf("ConfirmFactor() error: %v", err)
	}
	if err := s.ReplaceRecoveryCodes(ctx, u.ID, []string{"hash-a"}); err != nil {
		t.Fatalf("ReplaceRecoveryCodes() error: %v", err)
	}

	if err := s.DeleteUser(ctx, u.ID); err != nil {
		t.Fatalf("DeleteUser() error: %v", err)
	}

	if _, err := s.GetFactor(ctx, factor.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("GetFactor() after user delete = %v; want ErrNotFound", err)
	}
	remaining, err := s.CountRecoveryCodesRemaining(ctx, u.ID)
	if err != nil || remaining != 0 {
		t.Fatalf("CountRecoveryCodesRemaining() after user delete = %d, %v; want 0, nil", remaining, err)
	}
}

func TestGetRequire2FADefaultsToNone(t *testing.T) {
	s := newDocTestStore(t)
	ctx := context.Background()

	// Absent the auth_config singleton row entirely (newDocTestStore doesn't
	// seed it), GetRequire2FA still falls back to "none".
	require2FA, err := s.GetRequire2FA(ctx)
	if err != nil || require2FA != "none" {
		t.Fatalf("GetRequire2FA() = %q, %v; want none, nil", require2FA, err)
	}

	if err := s.initSingletons(ctx); err != nil {
		t.Fatalf("initSingletons() error: %v", err)
	}
	if _, err := s.DB().ExecContext(ctx, `UPDATE auth_config SET require_2fa = 'all' WHERE id = 1`); err != nil {
		t.Fatalf("update auth_config: %v", err)
	}
	require2FA, err = s.GetRequire2FA(ctx)
	if err != nil || require2FA != "all" {
		t.Fatalf("GetRequire2FA() after update = %q, %v; want all, nil", require2FA, err)
	}
}

func TestDeleteFactorLastOneAlsoLeavesRecoveryCodesForCallerToWipe(t *testing.T) {
	s := newDocTestStore(t)
	ctx := context.Background()
	u := mfaTestUser(t, s)

	factor, err := s.CreatePendingTOTP(ctx, u.ID, "Authenticator", "secret")
	if err != nil {
		t.Fatalf("CreatePendingTOTP() error: %v", err)
	}
	if _, err := s.ConfirmFactor(ctx, factor.ID); err != nil {
		t.Fatalf("ConfirmFactor() error: %v", err)
	}
	if err := s.DeleteFactor(ctx, factor.ID); err != nil {
		t.Fatalf("DeleteFactor() error: %v", err)
	}
	if _, err := s.GetFactor(ctx, factor.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("GetFactor() after delete = %v; want ErrNotFound", err)
	}
	hasMFA, err := s.UserHasMFA(ctx, u.ID)
	if err != nil || hasMFA {
		t.Fatalf("UserHasMFA() after deleting only factor = %v, %v; want false, nil", hasMFA, err)
	}
}

func TestDeleteUserFactorsForAdminReset(t *testing.T) {
	s := newDocTestStore(t)
	ctx := context.Background()
	u := mfaTestUser(t, s)

	factor, err := s.CreatePendingTOTP(ctx, u.ID, "Authenticator", "secret")
	if err != nil {
		t.Fatalf("CreatePendingTOTP() error: %v", err)
	}
	if _, err := s.ConfirmFactor(ctx, factor.ID); err != nil {
		t.Fatalf("ConfirmFactor() error: %v", err)
	}
	if err := s.DeleteUserFactors(ctx, u.ID); err != nil {
		t.Fatalf("DeleteUserFactors() error: %v", err)
	}
	hasMFA, err := s.UserHasMFA(ctx, u.ID)
	if err != nil || hasMFA {
		t.Fatalf("UserHasMFA() after admin reset = %v, %v; want false, nil", hasMFA, err)
	}
}
