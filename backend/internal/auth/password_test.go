package auth

import "testing"

func TestHashAndVerify(t *testing.T) {
	password := "secure-password-123"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() error: %v", err)
	}

	if hash == "" {
		t.Fatal("hash is empty")
	}
	if hash == password {
		t.Fatal("hash equals plaintext password")
	}

	if err := VerifyPassword(hash, password); err != nil {
		t.Fatalf("VerifyPassword() error: %v", err)
	}
}

func TestVerifyWrongPassword(t *testing.T) {
	hash, _ := HashPassword("correct-password")

	if err := VerifyPassword(hash, "wrong-password"); err == nil {
		t.Error("expected error for wrong password")
	}
}

func TestHashTooShort(t *testing.T) {
	_, err := HashPassword("short")
	if err == nil {
		t.Error("expected error for password < 8 chars")
	}
}

func TestBcryptCostOnlyLoweredUnderTest(t *testing.T) {
	if productionBcryptCost != 12 {
		t.Fatalf("productionBcryptCost = %d, want 12", productionBcryptCost)
	}
	if bcryptCost == productionBcryptCost {
		t.Fatal("bcryptCost should be lowered inside test binaries")
	}
}
