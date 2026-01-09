package client

import (
	"fmt"
	"sync"
)

// ClientRegistry handles the storage of authenticated clients.
// It is thread-safe to allow concurrent resource operations.
type ClientRegistry struct {
	clients map[string]Client
	lock    sync.RWMutex
}

func NewRegistry() *ClientRegistry {
	return &ClientRegistry{
		clients: make(map[string]Client),
	}
}

// Store saves a configured client under a specific ID.
func (r *ClientRegistry) Store(id string, c Client) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.clients[id] = c
}

// Load retrieves a client by ID.
func (r *ClientRegistry) Load(id string) (Client, bool) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	c, ok := r.clients[id]
	return c, ok
}

// LoadOrError is a helper to reduce boilerplate in resources.
func (r *ClientRegistry) LoadOrError(id string) (Client, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	c, ok := r.clients[id]
	if !ok {
		return nil, fmt.Errorf("client with id '%s' not found. Ensure the rancher_client resource is configured correctly", id)
	}
	return c, nil
}

// Delete removes a client by ID.
func (r *ClientRegistry) Delete(id string) {
	r.lock.Lock()
	defer r.lock.Unlock()
	delete(r.clients, id)
}
