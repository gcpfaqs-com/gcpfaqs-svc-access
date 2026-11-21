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

func (g *JWTGenerator) DecodeToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{},
		func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return g.secret, nil
		})

	if err != nil {
		return "", fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		return claims.Subject, nil
	}

	return "", fmt.Errorf("invalid token")
}
