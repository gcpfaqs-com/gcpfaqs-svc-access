package secure

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTGenerator struct {
	secret []byte
	issuer string
	ttl    time.Duration
}

func NewJWTGenerator(secret []byte, issuer string, ttl time.Duration) *JWTGenerator {
	return &JWTGenerator{
		secret: secret,
		issuer: issuer,
		ttl:    ttl,
	}
}

func (g *JWTGenerator) Generate(subject string) (string, error) {
	now := time.Now().UTC()

	claims := jwt.RegisteredClaims{
		Issuer:    g.issuer,
		Subject:   subject,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(g.ttl)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(g.secret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signed, nil
}
