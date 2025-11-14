package server

import (
	"encoding/json"
)

type SessionManager struct {
	cache *PersistentCache
}

func NewSessionManager(cache *PersistentCache) *SessionManager {
	return &SessionManager{cache: cache}
}

func (s *SessionManager) SaveSession(sessionID string, data interface{}) error {
	s.cache.Set("session:"+sessionID, data)
	s.cache.Snapshot("session:" + sessionID)
	return nil
}

func (s *SessionManager) LoadSession(sessionID string) (interface{}, bool) {
	v, ok := s.cache.Get("session:" + sessionID)
	return v, ok
}

func (s *SessionManager) SaveAppSnapshot(appId string, snap *AppSnapshot) error {
	b, _ := json.Marshal(snap)
	s.cache.Set("app:"+appId, string(b))
	s.cache.Snapshot("app:" + appId)
	return nil
}
