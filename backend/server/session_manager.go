
// =================================
// backend/server/session_manager.go
// =================================
package server

import (
	"encoding/json"
)

type SessionManager struct {
	bus *BusServer
	vfs *VFSService
    session *Session
}

type Session struct {
    Username string `json:"username"`
    // Other session data would go here
}

func NewSessionManager(bus *BusServer, vfs *VFSService) *SessionManager {
    sm := &SessionManager{bus: bus, vfs: vfs}
    bus.Subscribe("session:save", sm.handleSaveSession)
    sm.LoadSession() // Try to load the session on startup
    return sm
}

func (s *SessionManager) LoadSession() error {
    data, err := s.vfs.Read("/sessions/latest/session.json")
    if err != nil {
        // If the file doesn't exist, create a default session
        s.session = &Session{Username: "guest"}
        return s.SaveSession()
    }
    var session Session
    if err := json.Unmarshal(data, &session); err != nil {
        return err
    }
    s.session = &session
    return nil
}

func (s *SessionManager) SaveSession() error {
	data, err := json.Marshal(s.session)
    if err != nil {
        return err
    }
	return s.vfs.Write("/sessions/latest/session.json", data)
}

func (s *SessionManager) handleSaveSession(msg *Message) {
    // In a real app, we would update the session from the message payload
    err := s.SaveSession()
    if err != nil {
        s.bus.Reply(msg, map[string]interface{}{"error": err.Error()})
        return
    }
    s.bus.Reply(msg, map[string]interface{}{"success": true})
}
