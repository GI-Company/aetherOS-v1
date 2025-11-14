package server

import (
	"log"
)

// VFSService handles all file system operations
type VFSService struct {
	bus *BusServer
}

// NewVFSService creates a new VFSService
func NewVFSService(bus *BusServer) *VFSService {
	v := &VFSService{
		bus: bus,
	}
	bus.SubscribeServer("vfs:write", v.handleWrite)
	bus.SubscribeServer("vfs:read", v.handleRead)
	bus.SubscribeServer("vfs:list", v.handleList)
	log.Println("VFS Service Started")
	return v
}

func (v *VFSService) handleWrite(env *Envelope) {
	log.Printf("VFS Service received write request: %+v", env.Payload)
	// 1. Validate token caps (Phase 1.4)
	// 2. Resolve path, check permissions
	// 3. Forward message to frontend storage proxy
	//    (This will require the WebSocket bridge to be connected to the bus)
}

func (v *VFSService) handleRead(env *Envelope) {
	log.Printf("VFS Service received read request: %+v", env.Payload)
}

func (v *VFSService) handleList(env *Envelope) {
	log.Printf("VFS Service received list request: %+v", env.Payload)
}
