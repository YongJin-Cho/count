package main

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var SecretKey = []byte("secret")

type Claims struct {
	Permissions []string `json:"permissions"`
	jwt.RegisteredClaims
}

func main() {
	claims := Claims{
		Permissions: []string{"collect"},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			Issuer:    "test",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(SecretKey)
	if err != nil {
		fmt.Printf("Error signing token: %v\n", err)
		return
	}
	fmt.Print(ss)
}
