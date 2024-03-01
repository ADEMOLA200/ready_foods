package util

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

const secretKey = "secret"

func GenerateJwt(issuer string) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    issuer,
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})

	return claims.SignedString([]byte(secretKey))
}

func ParseJwt(cookie string) (string, error) {
	// Define the claims type as a pointer to jwt.StandardClaims
	claims := &jwt.StandardClaims{}

	token, err := jwt.ParseWithClaims(cookie, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil || !token.Valid {
		return "", err
	}

	return claims.Issuer, nil
}
