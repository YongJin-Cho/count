package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func generateToken(permissions []string, expired bool) string {
	expirationTime := time.Now().Add(5 * time.Minute)
	if expired {
		expirationTime = time.Now().Add(-5 * time.Minute)
	}

	claims := &Claims{
		Permissions: permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(SecretKey)
	return tokenString
}

func TestAuthProvider_ValidateToken(t *testing.T) {
	provider := NewAuthProvider()

	t.Run("Valid token", func(t *testing.T) {
		token := generateToken([]string{"collect"}, false)
		valid, err := provider.ValidateToken(token)
		if err != nil || !valid {
			t.Errorf("Expected valid token, got err=%v, valid=%v", err, valid)
		}
	})

	t.Run("Expired token", func(t *testing.T) {
		token := generateToken([]string{"collect"}, true)
		valid, err := provider.ValidateToken(token)
		if err == nil || valid {
			t.Errorf("Expected invalid token (expired), got err=%v, valid=%v", err, valid)
		}
	})

	t.Run("Malformed token", func(t *testing.T) {
		valid, err := provider.ValidateToken("not.a.token")
		if err == nil || valid {
			t.Errorf("Expected invalid token (malformed), got err=%v, valid=%v", err, valid)
		}
	})
}

func TestAuthProvider_IsAuthorized(t *testing.T) {
	provider := NewAuthProvider()

	t.Run("Authorized", func(t *testing.T) {
		token := generateToken([]string{"collect"}, false)
		authorized, err := provider.IsAuthorized(token, "collect")
		if err != nil || !authorized {
			t.Errorf("Expected authorized, got err=%v, authorized=%v", err, authorized)
		}
	})

	t.Run("Forbidden", func(t *testing.T) {
		token := generateToken([]string{"other"}, false)
		authorized, err := provider.IsAuthorized(token, "collect")
		if err != nil || authorized {
			t.Errorf("Expected forbidden, got err=%v, authorized=%v", err, authorized)
		}
	})
}
