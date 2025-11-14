package permissions

import (
	"log"
	"time"

	"aether/backend/server"
)

// PermissionsService handles all permission and capability checks.
type PermissionsService struct {
	bus *server.BusServer
}

// NewPermissionsService creates a new PermissionsService.
func NewPermissionsService(bus *server.BusServer) *PermissionsService {
	ps := &PermissionsService{
		bus: bus,
	}
	bus.SubscribeServer("permissions:check", ps.handlePermissionCheck)
	return ps
}

// handlePermissionCheck requires a valid JWT.
// It will ask the auth service to verify the JWT and then, if valid,
// will (for now) approve the request.
func (s *PermissionsService) handlePermissionCheck(env *server.Envelope) {
	log.Printf("PermissionsService received check request: %v", env.Payload)

	jwt, ok := env.Payload["jwt"].(string)
	if !ok {
		log.Println("Permissions check requires a JWT.")
		s.replyToRequest(env, "deny", "JWT not found in payload")
		return
	}

	// Use the bus to make a request to the auth service to verify the JWT
	authRequest := &server.Envelope{
		Topic: "auth:verify:jwt",
		Payload: map[string]interface{}{
			"token": jwt,
		},
	}

	// Use Request/Reply to wait for the auth service response
	authResponse, ok := s.bus.Request(authRequest, 2*time.Second)
	if !ok {
		log.Println("Permissions check failed: no response from auth service")
		s.replyToRequest(env, "deny", "auth service timeout")
		return
	}

	if authResponse.Topic == "auth:jwt:valid" {
		log.Println("Permissions check successful, JWT is valid.")
		// In the future, we would check the claims in the JWT against
		// the requested action.
		s.replyToRequest(env, "allow", "")
	} else {
		log.Println("Permissions check failed, JWT is invalid.")
		s.replyToRequest(env, "deny", "invalid JWT")
	}
}

func (s *PermissionsService) replyToRequest(originalEnv *server.Envelope, result string, errorMsg string) {
	payload := map[string]interface{}{"result": result}
	if errorMsg != "" {
		payload["error"] = errorMsg
	}

	reply := &server.Envelope{
		Topic:   "permissions:check:result",
		Payload: payload,
	}

	if requestID, ok := originalEnv.Payload["_request_id"].(string); ok {
		reply.Payload["_reply_to"] = requestID
	}

	s.bus.Publish(reply)
}
