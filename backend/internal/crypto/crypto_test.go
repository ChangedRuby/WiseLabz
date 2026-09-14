package crypto

import (
	"encoding/base64"
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	secret := "test-secret-key-for-encryption"
	key := DeriveKey(secret)

	if len(key) != 32 {
		t.Fatalf("DeriveKey returned %d bytes, want 32", len(key))
	}

	plaintext := "sk-ant-api03-abcdefghijklmnopqrstuvwxyz"
	encrypted, err := Encrypt(plaintext, key)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}
	if encrypted == "" {
		t.Fatal("Encrypt returned empty string")
	}
	if encrypted == plaintext {
		t.Fatal("Encrypted text equals plaintext — not encrypted")
	}

	decrypted, err := Decrypt(encrypted, key)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}
	if decrypted != plaintext {
		t.Fatalf("Decrypt returned %q, want %q", decrypted, plaintext)
	}
}

func TestDecryptTampered(t *testing.T) {
	key := DeriveKey("test-secret")
	encrypted, _ := Encrypt("my-api-key", key)

	// Corrupt the ciphertext
	tampered := encrypted[:len(encrypted)-4] + "AAAA"
	_, err := Decrypt(tampered, key)
	if err == nil {
		t.Fatal("Decrypt should fail on tampered ciphertext")
	}
}

func TestDecryptWrongKey(t *testing.T) {
	key1 := DeriveKey("secret-one")
	key2 := DeriveKey("secret-two")

	encrypted, _ := Encrypt("my-api-key", key1)
	_, err := Decrypt(encrypted, key2)
	if err == nil {
		t.Fatal("Decrypt should fail with wrong key")
	}
}

func TestEncryptEmptyPlaintext(t *testing.T) {
	key := DeriveKey("test-secret")
	encrypted, err := Encrypt("", key)
	if err != nil {
		t.Fatalf("Encrypt empty string failed: %v", err)
	}
	decrypted, err := Decrypt(encrypted, key)
	if err != nil {
		t.Fatalf("Decrypt empty string failed: %v", err)
	}
	if decrypted != "" {
		t.Fatalf("Decrypt empty returned %q", decrypted)
	}
}

func TestEncryptBadKeySize(t *testing.T) {
	_, err := Encrypt("test", make([]byte, 16))
	if err == nil {
		t.Fatal("Encrypt should reject 16-byte key")
	}

	_, err = Decrypt("test", make([]byte, 16))
	if err == nil {
		t.Fatal("Decrypt should reject 16-byte key")
	}
}

func TestDeriveKeyDeterministic(t *testing.T) {
	k1 := DeriveKey("same-secret")
	k2 := DeriveKey("same-secret")
	if string(k1) != string(k2) {
		t.Fatal("DeriveKey not deterministic for same input")
	}
}

func TestDeriveKeyDifferent(t *testing.T) {
	k1 := DeriveKey("secret-a")
	k2 := DeriveKey("secret-b")
	if string(k1) == string(k2) {
		t.Fatal("DeriveKey produced same key for different inputs")
	}
}

func TestDecodeKey(t *testing.T) {
	valid32 := base64.StdEncoding.EncodeToString(make([]byte, 32))
	short16 := base64.StdEncoding.EncodeToString(make([]byte, 16))
	long64 := base64.StdEncoding.EncodeToString(make([]byte, 64))

	tests := []struct {
		name    string
		in      string
		wantErr bool
	}{
		{name: "valid 32-byte key", in: valid32, wantErr: false},
		{name: "empty string", in: "", wantErr: true},
		{name: "malformed base64", in: "not-valid-base64!!!", wantErr: true},
		{name: "too short (16 bytes)", in: short16, wantErr: true},
		{name: "too long (64 bytes)", in: long64, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, err := DecodeKey(tt.in)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("DecodeKey(%q) = nil error, want error", tt.in)
				}
				return
			}
			if err != nil {
				t.Fatalf("DecodeKey(%q) unexpected error: %v", tt.in, err)
			}
			if len(key) != 32 {
				t.Fatalf("DecodeKey(%q) returned %d bytes, want 32", tt.in, len(key))
			}
		})
	}
}

func TestDecodeKeyUsableForEncryptDecrypt(t *testing.T) {
	b64 := base64.StdEncoding.EncodeToString(make([]byte, 32))
	key, err := DecodeKey(b64)
	if err != nil {
		t.Fatalf("DecodeKey failed: %v", err)
	}
	encrypted, err := Encrypt("payload", key)
	if err != nil {
		t.Fatalf("Encrypt with decoded key failed: %v", err)
	}
	decrypted, err := Decrypt(encrypted, key)
	if err != nil {
		t.Fatalf("Decrypt with decoded key failed: %v", err)
	}
	if decrypted != "payload" {
		t.Fatalf("Decrypt returned %q, want %q", decrypted, "payload")
	}
}
