package secure

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const defaultBCryptCost = 12

type BCryptHasher struct {
	cost int
}

func NewBCryptHasher() *BCryptHasher {
	return &BCryptHasher{
		cost: defaultBCryptCost,
	}
}

func (h *BCryptHasher) Hash(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return string(hashedBytes), nil
}

func (h *BCryptHasher) Verify(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
