package feeder

import (
	"sync"
)

type Validator func(string) error

type Manager struct {
	store Store
	validator Validator

	countInvalid int
	mutex sync.Mutex
}

func NewManager(store Store, validator Validator) *Manager {
	return &Manager{
		store: store,
		validator: validator,
		countInvalid: 0,
	}
}

func (m *Manager) HandleMessage(sku string) {
	err := m.validator(sku)
	if err != nil {
		m.mutex.Lock()
		defer m.mutex.Unlock()

		m.countInvalid++
		return
	}

	m.store.Insert(sku)
}
