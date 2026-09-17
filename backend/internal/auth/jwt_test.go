package auth

import (
	"sync"
	"testing"
	"time"
)

func TestIssueAndValidateAccess(t *testing.T) {
	svc := NewService("test-secret-key-32-bytes-long!!", 15*time.Minute, 7*24*time.Hour)

	pair, err := svc.IssuePair("user-123", true)
	if err != nil {
		t.Fatalf("IssuePair() error: %v", err)
	}

	if pair.AccessToken == "" {
		t.Error("access token is empty")
	}
	if pair.RefreshToken == "" {
		t.Error("refresh token is empty")
	}
	if pair.ExpiresIn != 900 {
		t.Errorf("ExpiresIn = %d, want 900", pair.ExpiresIn)
	}

	claims, err := svc.ValidateAccess(pair.AccessToken)
	if err != nil {
		t.Fatalf("ValidateAccess() error: %v", err)
	}
	if claims.UserID != "user-123" {
		t.Errorf("UserID = %q, want user-123", claims.UserID)
	}
	if !claims.InstanceAdmin {
		t.Errorf("InstanceAdmin = %v, want true", claims.InstanceAdmin)
	}
}

func TestValidateRefresh(t *testing.T) {
	svc := NewService("test-secret", time.Minute, time.Hour)

	pair, err := svc.IssuePair("user-456", false)
	if err != nil {
		t.Fatalf("IssuePair() error: %v", err)
	}

	claims, err := svc.ValidateRefresh(pair.RefreshToken)
	if err != nil {
		t.Fatalf("ValidateRefresh() error: %v", err)
	}
	if claims.UserID != "user-456" {
		t.Errorf("UserID = %q, want user-456", claims.UserID)
	}
}

func TestTokenPurposesCannotBeExchanged(t *testing.T) {
	svc := NewService("test-secret", time.Minute, time.Hour)
	pair, err := svc.IssuePair("user-456", false)
	if err != nil {
		t.Fatalf("IssuePair() error: %v", err)
	}
	if _, err := svc.ValidateAccess(pair.RefreshToken); err == nil {
		t.Error("refresh token validated as access token")
	}
	if _, err := svc.ValidateRefresh(pair.AccessToken); err == nil {
		t.Error("access token validated as refresh token")
	}
}

func TestElevationRequiresOwner(t *testing.T) {
	svc := NewService("test-secret", time.Minute, time.Hour)
	elev, err := svc.IssueElevation("owner", "connector.delete")
	if err != nil {
		t.Fatalf("IssueElevation() error: %v", err)
	}
	if _, err := svc.ValidateElevation(elev.Token, "connector.delete", "other"); err == nil {
		t.Error("elevation token validated for another user")
	}
}

func TestExpiredAccessToken(t *testing.T) {
	svc := NewService("test-secret", -1*time.Second, time.Hour) // access expires immediately

	pair, err := svc.IssuePair("user-789", false)
	if err != nil {
		t.Fatalf("IssuePair() error: %v", err)
	}

	_, err = svc.ValidateAccess(pair.AccessToken)
	if err == nil {
		t.Error("expected error for expired access token, got nil")
	}
}

func TestWrongSecret(t *testing.T) {
	svc1 := NewService("secret-one", time.Minute, time.Hour)
	svc2 := NewService("secret-two", time.Minute, time.Hour)

	pair, _ := svc1.IssuePair("user-1", false)

	_, err := svc2.ValidateAccess(pair.AccessToken)
	if err == nil {
		t.Error("expected error for token signed with different secret")
	}
}

func TestIssueAndValidateElevation(t *testing.T) {
	svc := NewService("test-secret", time.Minute, time.Hour)

	elev, err := svc.IssueElevation("user-123", "connector.delete")
	if err != nil {
		t.Fatalf("IssueElevation() error: %v", err)
	}

	if elev.Token == "" {
		t.Error("elevation token is empty")
	}
	if time.Until(elev.ExpiresAt) <= 0 {
		t.Error("elevation token already expired")
	}

	claims, err := svc.ValidateElevation(elev.Token, "connector.delete", "user-123")
	if err != nil {
		t.Fatalf("ValidateElevation() error: %v", err)
	}
	if claims.UserID != "user-123" {
		t.Errorf("UserID = %q, want user-123", claims.UserID)
	}
	if claims.Action != "connector.delete" {
		t.Errorf("Action = %q, want connector.delete", claims.Action)
	}
}

func TestElevationWrongAction(t *testing.T) {
	svc := NewService("test-secret", time.Minute, time.Hour)

	elev, _ := svc.IssueElevation("user-123", "connector.delete")

	_, err := svc.ValidateElevation(elev.Token, "user.delete", "user-123")
	if err == nil {
		t.Error("expected error for wrong action on elevation token")
	}
}

// TestConcurrentIssuePairUniqueTokenIDs guards against the data race fixed
// in issue #148: newTokenID() used to be an unsynchronized package-level
// counter, so concurrent issuance could mint tokens with duplicate jti.
// Run with -race to catch a regression to that shape.
func TestConcurrentIssuePairUniqueTokenIDs(t *testing.T) {
	svc := NewService("test-secret", time.Minute, time.Hour)

	const n = 50
	var wg sync.WaitGroup
	ids := make(chan string, n*2)
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			pair, err := svc.IssuePair("user", false)
			if err != nil {
				t.Error(err)
				return
			}
			access, err := svc.ValidateAccess(pair.AccessToken)
			if err != nil {
				t.Error(err)
				return
			}
			refresh, err := svc.ValidateRefresh(pair.RefreshToken)
			if err != nil {
				t.Error(err)
				return
			}
			ids <- access.ID
			ids <- refresh.ID
		}()
	}
	wg.Wait()
	close(ids)

	seen := make(map[string]bool)
	for id := range ids {
		if id == "" {
			t.Fatal("jti is empty")
		}
		if seen[id] {
			t.Fatalf("duplicate jti: %s", id)
		}
		seen[id] = true
	}
}

func TestElevationExpired(t *testing.T) {
	svc := NewService("test-secret", time.Minute, time.Hour)
	svc.elevationTTL = -1 * time.Second // immediately expired

	elev, err := svc.IssueElevation("user-123", "connector.delete")
	if err != nil {
		t.Fatalf("IssueElevation() error: %v", err)
	}

	_, err = svc.ValidateElevation(elev.Token, "connector.delete", "user-123")
	if err == nil {
		t.Error("expected error for expired elevation token")
	}
}
