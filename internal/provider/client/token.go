package client

import (
	"sync"
)

// TokenStore safely manages concurrent access to the auth token.
type TokenStore struct {
	mu    sync.RWMutex
	token string
}

// SetToken safely writes the token (used by rancher_login).
func (ts *TokenStore) SetToken(token string) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	ts.token = token
}

// GetToken safely reads the token (used by the clients).
func (ts *TokenStore) GetToken() string {
	ts.mu.RLock()
	defer ts.mu.RUnlock()
	return ts.token
}

func (ts *TokenStore) ClearToken() {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	ts.token = ""
}
