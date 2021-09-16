package feeder

import (
	"strings"
	"sync"
)

// Validator defines the function type for validating skus
type Validator func(string) error

// Store defines the contract of a service to store valid skus
type Store interface {
	Insert(string)
}

// Manager handles incoming messages, validates and stores them
type Manager struct {
	store     Store
	validator Validator

	countInvalid int
	mutex        sync.Mutex
}

// NewManager returns a new instance of Manager
func NewManager(store Store, validator Validator) *Manager {
	return &Manager{
		store:        store,
		validator:    validator,
		countInvalid: 0,
	}
}

// HandleMessage handles a new incoming message and stores it in store if it is valid
func (m *Manager) HandleMessage(sku string) {
	err := m.validator(sku)
	if err != nil {
		m.mutex.Lock()
		defer m.mutex.Unlock()

		m.countInvalid++
		return
	}
	// store skus in uppercase
	sku = strings.ToUpper(sku)
	m.store.Insert(sku)
}
