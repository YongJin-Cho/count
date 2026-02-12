package auth

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
)

var SecretKey = []byte("secret")

type Claims struct {
	Permissions []string `json:"permissions"`
	jwt.RegisteredClaims
}

type AuthProvider struct{}

func NewAuthProvider() *AuthProvider {
	return &AuthProvider{}
}

func (p *AuthProvider) ValidateToken(tokenString string) (bool, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return SecretKey, nil
	})

	if err != nil {
		return false, err
	}

	return token.Valid, nil
}

func (p *AuthProvider) IsAuthorized(tokenString string, permission string) (bool, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})

	if err != nil {
		return false, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		for _, p := range claims.Permissions {
			if p == permission {
				return true, nil
			}
		}
		return false, nil
	}

	return false, errors.New("invalid claims")
}
