
// ===========================
// backend/server/vfs.go
// ===========================
package server

import (
	"fmt"
	"strings"
	"sync"
)

type VFSService struct {
	bus      *BusServer
	mu       sync.RWMutex
	root     *VFSNode
}

type VFSNode struct {
	name     string
	isDir    bool
	content  []byte
	children map[string]*VFSNode
}

func NewVFSService(bus *BusServer) *VFSService {
	vfs := &VFSService{
		bus: bus,
		root: &VFSNode{name: "/", isDir: true, children: make(map[string]*VFSNode)},
	}
	bus.Subscribe("vfs:read", vfs.handleRead)
	bus.Subscribe("vfs:write", vfs.handleWrite)
	return vfs
}

func (vfs *VFSService) Read(path string) ([]byte, error) {
	vfs.mu.RLock()
	defer vfs.mu.RUnlock()

	parts := strings.Split(path, "/")
	curr := vfs.root

	for i, part := range parts {
		if part == "" {
			continue
		}
		node, ok := curr.children[part]
		if !ok {
			return nil, fmt.Errorf("path not found: %s", path)
		}
		if i == len(parts)-1 {
			if node.isDir {
				return nil, fmt.Errorf("path is a directory: %s", path)
			}
			return node.content, nil
		}
		curr = node
	}
	return nil, fmt.Errorf("path not found: %s", path)
}

func (vfs *VFSService) Write(path string, data []byte) error {
	vfs.mu.Lock()
	defer vfs.mu.Unlock()

	parts := strings.Split(path, "/")
	curr := vfs.root

	for i, part := range parts {
		if part == "" {
			continue
		}
		if i == len(parts)-1 {
			curr.children[part] = &VFSNode{name: part, content: data}
		} else {
			node, ok := curr.children[part]
			if !ok || !node.isDir {
                // For simplicity, we automatically create directories that don't exist
                newNode := &VFSNode{name: part, isDir: true, children: make(map[string]*VFSNode)}
                curr.children[part] = newNode
                curr = newNode
			} else {
                curr = node
            }
		}
	}
	return nil
}


func (vfs *VFSService) handleRead(msg *Message) {
    path := msg.Payload["path"].(string)
    // TODO: Auth checks
    content, err := vfs.Read(path)
    if err != nil {
        vfs.bus.Reply(msg, map[string]interface{}{"error": err.Error()})
        return
    }
    vfs.bus.Reply(msg, map[string]interface{}{"content": content})
}

func (vfs *VFSService) handleWrite(msg *Message) {
    path := msg.Payload["path"].(string)
    content := msg.Payload["content"].([]byte)
    // TODO: Auth checks
    err := vfs.Write(path, content)
    if err != nil {
        vfs.bus.Reply(msg, map[string]interface{}{"error": err.Error()})
        return
    }
    vfs.bus.Reply(msg, map[string]interface{}{"success": true})
}
