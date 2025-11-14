
// =================================
// backend/server/vfs_service.go
// =================================
package server

import (
	"encoding/json"
)

type VFSService struct {
	bus *BusServer
}

func (v *VFSService) Write(path string, content []byte) error {
	msg := map[string]interface{}{ "path": path, "content": content }
	return v.saveToIndexedDBProxy(msg)
}

func (v *VFSService) saveToIndexedDBProxy(msg map[string]interface{}) error {
	b, _ := json.Marshal(msg)
	v.bus.Publish(&Message{Topic: "vfs:write", Payload: string(b)})
	return nil
}

func NewVFSService(bus *BusServer) *VFSService {
    return &VFSService{bus: bus}
}
