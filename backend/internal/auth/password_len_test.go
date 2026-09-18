package auth

import (
	"strings"
	"testing"
)

func TestHashPasswordRejectsOver72Bytes(t *testing.T) {
	if _, err := HashPassword(strings.Repeat("a", 73)); err == nil {
		t.Fatal("expected error for 73-byte password")
	}
	if _, err := HashPassword(strings.Repeat("a", 72)); err != nil {
		t.Fatalf("72-byte password should be accepted: %v", err)
	}
}
