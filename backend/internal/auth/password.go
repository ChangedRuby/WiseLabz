package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const bcryptCost = 12

// maxPasswordBytes is bcrypt's input limit.
const maxPasswordBytes = 72

// HashPassword returns the bcrypt hash of a password.
func HashPassword(password string) (string, error) {
	if len(password) < 8 {
		return "", fmt.Errorf("password must be at least 8 characters")
	}
	// bcrypt silently truncates input beyond 72 bytes; reject rather than
	// let distinct long passwords collide.
	if len(password) > maxPasswordBytes {
		return "", fmt.Errorf("password must be at most %d bytes", maxPasswordBytes)
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return "", fmt.Errorf("hash password: %w", err)
	}
	return string(hash), nil
}

// VerifyPassword compares a password against its bcrypt hash.
// Returns nil if they match.
func VerifyPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
