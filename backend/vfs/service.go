package vfs

import (
	"log"

	"aether/backend/bus"
)

// VFSService handles all file system operations
type VFSService struct {
	bus    *bus.Bus
	client *bus.Client
}

// NewVFSService creates a new VFSService
func NewVFSService(bus *bus.Bus) *VFSService {
	return &VFSService{
		bus: bus,
		client: &bus.Client{
			ID:      "vfs-service",
			Receive: make(chan bus.Message, 128), // Buffered channel
		},
	}
}

// Start begins the VFS service's message listening loop
func (v *VFSService) Start() {
	log.Println("Starting VFS Service")
	v.bus.Subscribe("vfs:write", v.client)
	v.bus.Subscribe("vfs:read", v.client)
	v.bus.Subscribe("vfs:list", v.client)

	go v.listen()
}

func (v *VFSService) listen() {
	for msg := range v.client.Receive {
		switch msg.Topic {
		case "vfs:write":
			v.handleWrite(msg)
		case "vfs:read":
			v.handleRead(msg)
		case "vfs:list":
			v.handleList(msg)
		}
	}
}

func (v *VFSService) handleWrite(msg bus.Message) {
	log.Printf("VFS Service received write request: %+v", msg.Payload)
	// 1. Validate token caps (Phase 1.4)
	// 2. Resolve path, check permissions
	// 3. Forward message to frontend storage proxy
	//    (This will require the WebSocket bridge to be connected to the bus)
}

func (v *VFSService) handleRead(msg bus.Message) {
	// Implementation for reading files
}

func (v *VFSService) handleList(msg bus.Message) {
	// Implementation for listing files
}
