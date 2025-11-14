package permissions

import (
	"log"

	"aether-broker/backend/bus"
)

// PermissionsService handles all permission and capability checks.
// It now requires a valid Aether JWT for all checks.
type PermissionsService struct {
	bus    *bus.Bus
	client *bus.Client
}

// NewPermissionsService creates a new PermissionsService.
func NewPermissionsService(bus *bus.Bus) *PermissionsService {
	return &PermissionsService{
		bus:    bus,
		client: &bus.Client{ID: "permissions-service", Receive: make(chan bus.Message, 128)},
	}
}

// Start begins the PermissionsService's message listening loop.
func (s *PermissionsService) Start() {
	log.Println("Starting Permissions Service")
	s.bus.Subscribe("permissions:check", s.client)
	go s.listen()
}

func (s *PermissionsService) listen() {
	for msg := range s.client.Receive {
		if msg.Topic == "permissions:check" {
			go s.handlePermissionCheck(msg)
		}
	}
}

// handlePermissionCheck now requires a valid JWT.
// It will ask the auth service to verify the JWT and then, if valid,
// will (for now) approve the request.
func (s *PermissionsService) handlePermissionCheck(msg bus.Message) {
	log.Printf("PermissionsService received check request: %v", msg.Payload)

	// For now, we'll assume the JWT is passed in the payload.
	// In a real system, this would be more structured.
	jwt, ok := msg.Payload.(string)
	if !ok {
		log.Println("Permissions check requires a JWT.")
		s.bus.Publish(bus.Message{Topic: "permissions:check:result", Payload: "deny"})
		return
	}

	// We need to wait for the result of the JWT verification.
	// This requires a synchronous call or a callback mechanism.
	// For simplicity, we'll use a temporary channel for the response.
	resultChan := make(chan bus.Message)
	resultClient := &bus.Client{ID: "permissions-temp-client", Receive: resultChan}
	s.bus.Subscribe("auth:jwt:valid", resultClient)
	s.bus.Subscribe("auth:jwt:invalid", resultClient)

	s.bus.Publish(bus.Message{Topic: "auth:verify:jwt", Payload: jwt})

	// Wait for the auth service to respond.
	result := <-resultChan

	s.bus.Unsubscribe("auth:jwt:valid", resultClient)
	s.bus.Unsubscribe("auth:jwt:invalid", resultClient)

	if result.Topic == "auth:jwt:valid" {
		log.Println("Permissions check successful, JWT is valid.")
		// In the future, we would check the claims in the JWT against
		// the requested action.
		s.bus.Publish(bus.Message{Topic: "permissions:check:result", Payload: "allow"})
	} else {
		log.Println("Permissions check failed, JWT is invalid.")
		s.bus.Publish(bus.Message{Topic: "permissions:check:result", Payload: "deny"})
	}
}
