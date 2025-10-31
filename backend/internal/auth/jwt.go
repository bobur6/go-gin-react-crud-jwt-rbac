package auth

import (
	"fmt"
	"time"

	"assignment3/backend/internal/models"

	"github.com/golang-jwt/jwt/v5"
)

// Claims wraps the user identity information included in JWT tokens.
type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// JWTService manages token generation and verification.
type JWTService struct {
	secret []byte
	issuer string
	expiry time.Duration
}

// NewJWTService constructs a JWT service instance.
func NewJWTService(secret, issuer string, expiry time.Duration) *JWTService {
	return &JWTService{
		secret: []byte(secret),
		issuer: issuer,
		expiry: expiry,
	}
}

// GenerateToken creates a signed JWT for the provided user.
func (j *JWTService) GenerateToken(user models.User) (string, error) {
	now := time.Now().UTC()
	claims := Claims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(j.expiry)),
			Subject:   user.ID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(j.secret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}
	return signed, nil
}

// ParseToken validates and parses a JWT string.
func (j *JWTService) ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.secret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
