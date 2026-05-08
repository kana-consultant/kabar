package user

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type PasswordService struct{}

func NewPasswordService() *PasswordService {
	return &PasswordService{}
}

// Hash hashes a password using bcrypt
func (p *PasswordService) Hash(password string) ([]byte, error) {
	if password == "" {
		return nil, fmt.Errorf("password cannot be empty")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	return hash, nil
}

// Verify checks if a password matches a hash
func (p *PasswordService) Verify(password string, hash []byte) error {
	if password == "" {
		return fmt.Errorf("password cannot be empty")
	}

	if len(hash) == 0 {
		return fmt.Errorf("hash cannot be empty")
	}

	return bcrypt.CompareHashAndPassword(hash, []byte(password))
}

// GenerateRandomPassword generates a secure random password
func (p *PasswordService) GenerateRandomPassword(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random password: %w", err)
	}

	return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}
