package server

import "sync"

type AuthService struct {
	mu     sync.RWMutex
	tokens map[string][]string // token -> caps
}

func NewAuthService() *AuthService {
	return &AuthService{tokens: make(map[string][]string)}
}

func (a *AuthService) IssueToken(caps []string) string {
	t := genUUID()
	a.mu.Lock()
	a.tokens[t] = caps
	a.mu.Unlock()
	return t
}

func (a *AuthService) ValidateCaps(token string, need string) bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	caps, ok := a.tokens[token]
	if !ok {
		return false
	}
	for _, c := range caps {
		if c == need {
			return true
		}
	}
	return false
}
