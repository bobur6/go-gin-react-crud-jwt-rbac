package auth_test

import (
	"testing"
	"time"

	"assignment3/backend/internal/auth"
	"assignment3/backend/internal/models"
)

func TestJWTServiceRoundTrip(t *testing.T) {
	service := auth.NewJWTService("secret", "issuer", time.Minute)

	user := models.User{
		ID:       "user-1",
		Username: "alice",
		Role:     "user",
	}

	token, err := service.GenerateToken(user)
	if err != nil {
		t.Fatalf("GenerateToken returned error: %v", err)
	}

	claims, err := service.ParseToken(token)
	if err != nil {
		t.Fatalf("ParseToken returned error: %v", err)
	}

	if claims.Username != user.Username || claims.UserID != user.ID || claims.Role != user.Role {
		t.Fatalf("unexpected claims: %+v", claims)
	}
}
