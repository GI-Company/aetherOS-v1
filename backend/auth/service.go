package auth

import (
	"context"
	"fmt"
	"log"
	"time"

	"aether/backend/bus"
	"firebase.google.com/go/v4/auth"
	"github.com/golang-jwt/jwt/v4"
)

// AuthService handles user authentication and session management using Firebase.
type AuthService struct {
	bus        *bus.Bus
	client     *bus.Client
	fbAuthClient *auth.Client
	jwtSecret  []byte
}

// NewAuthService creates a new AuthService.
func NewAuthService(bus *bus.Bus, fbAuthClient *auth.Client, jwtSecret string) *AuthService {
	return &AuthService{
		bus:        bus,
		client:     &bus.Client{ID: "auth-service", Receive: make(chan bus.Message, 128)},
		fbAuthClient: fbAuthClient,
		jwtSecret:  []byte(jwtSecret),
	}
}

// Start begins the AuthService's message listening loop.
func (s *AuthService) Start() {
	log.Println("Starting Auth Service")
	s.bus.Subscribe("auth:verify:token", s.client)
	s.bus.Subscribe("auth:verify:jwt", s.client)
	go s.listen()
}

func (s *AuthService) listen() {
	for msg := range s.client.Receive {
		switch msg.Topic {
		case "auth:verify:token":
			go s.handleVerifyToken(msg)
		case "auth:verify:jwt":
			go s.handleVerifyJWT(msg)
		}
	}
}

// handleVerifyToken receives a Firebase ID token from the frontend,
// verifies it, and if successful, issues a custom Aether session token.
func (s *AuthService) handleVerifyToken(msg bus.Message) {
	idToken, ok := msg.Payload.(string)
	if !ok {
		log.Println("Auth service received invalid payload for token verification")
		return
	}

	token, err := s.fbAuthClient.VerifyIDToken(context.Background(), idToken)
	if err != nil {
		log.Printf("Error verifying Firebase ID token: %v", err)
		s.bus.Publish(bus.Message{Topic: "auth:verify:failed", Payload: err.Error()})
		return
	}

	log.Printf("Verified user with Firebase: %s (%s)", token.UID, token.Claims["email"])

	// Generate a new Aether JWT
	aetherToken, err := s.generateAetherJWT(token.UID, token.Claims["email"].(string))
	if err != nil {
		log.Printf("Error generating Aether JWT: %v", err)
		s.bus.Publish(bus.Message{Topic: "auth:verify:failed", Payload: "could not generate session token"})
		return
	}

	s.bus.Publish(bus.Message{Topic: "auth:verify:success", Payload: aetherToken})
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
func (s *AuthService) handleVerifyJWT(msg bus.Message) {
	tokenString, ok := msg.Payload.(string)
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
		s.bus.Publish(bus.Message{Topic: "auth:jwt:invalid", OriginalMessage: &msg, Payload: err.Error()})
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		s.bus.Publish(bus.Message{Topic: "auth:jwt:valid", OriginalMessage: &msg, Payload: claims})
	} else {
		s.bus.Publish(bus.Message{Topic: "auth:jwt:invalid", OriginalMessage: &msg, Payload: "invalid token"})
	}
}
