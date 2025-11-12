package aethersdk

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

// For testing purposes only.
var jwtSecret = []byte("very-secret-key")

// NewJWT generates a new JWT for a given user ID.
// This is for testing and example purposes only.
func NewJWT(userID string, validity time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(validity).Unix(),
	})

	return token.SignedString(jwtSecret)
}
