package auth

import (
	"context"
	"fmt"
	"log"
	"time"

	"aether/backend/server"
	"firebase.google.com/go/v4/auth"
	jwt "github.com/golang-jwt/jwt/v4"
)

// AuthService handles user authentication and session management using Firebase.
type AuthService struct {
	bus        *server.BusServer
	fbAuthClient *auth.Client
	jwtSecret  []byte
}

// NewAuthService creates a new AuthService.
func NewAuthService(bus *server.BusServer, fbAuthClient *auth.Client, jwtSecret string) *AuthService {
	a := &AuthService{
		bus:        bus,
		fbAuthClient: fbAuthClient,
		jwtSecret:  []byte(jwtSecret),
	}
	a.bus.SubscribeServer("auth:verify:token", a.handleVerifyToken)
	a.bus.SubscribeServer("auth:verify:jwt", a.handleVerifyJWT)
	return a
}


// handleVerifyToken receives a Firebase ID token from the frontend,
// verifies it, and if successful, issues a custom Aether session token.
func (s *AuthService) handleVerifyToken(env *server.Envelope) {
	idToken, ok := env.Payload["token"].(string)
	if !ok {
		log.Println("Auth service received invalid payload for token verification")
		return
	}

	token, err := s.fbAuthClient.VerifyIDToken(context.Background(), idToken)
	if err != nil {
		log.Printf("Error verifying Firebase ID token: %v", err)
		// Ideally, send a reply with an error
		return
	}

	log.Printf("Verified user with Firebase: %s (%s)", token.UID, token.Claims["email"])

	// Generate a new Aether JWT
	aetherToken, err := s.generateAetherJWT(token.UID, token.Claims["email"].(string))
	if err != nil {
		log.Printf("Error generating Aether JWT: %v", err)
		// Ideally, send a reply with an error
		return
	}

	reply := &server.Envelope{
		Topic: "auth:verify:success",
		Payload: map[string]interface{}{"token": aetherToken},
	}
	if requestID, ok := env.Payload["_request_id"].(string); ok {
		reply.Payload["_reply_to"] = requestID
	}
	s.bus.Publish(reply)
}

func (s *AuthService) generateAetherJWT(uid, email string) (string, error) {
	claims := jwt.MapClaims{
		"uid":   uid,
		"email": email,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
		"iat":   time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

// handleVerifyJWT verifies an Aether session JWT.
func (s *AuthService) handleVerifyJWT(env *server.Envelope) {
	tokenString, ok := env.Payload["token"].(string)
	if !ok {
		log.Println("Auth service received invalid payload for JWT verification")
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		// Ideally, send a reply with an error
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		reply := &server.Envelope{
			Topic: "auth:jwt:valid",
			Payload: map[string]interface{}{"claims": claims},
		}
		if requestID, ok := env.Payload["_request_id"].(string); ok {
			reply.Payload["_reply_to"] = requestID
		}
		s.bus.Publish(reply)
	} else {
		// Ideally, send a reply with an error
	}
}
