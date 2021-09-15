package feeder

import (
	"strings"
	"sync"
)

type Validator func(string) error

type Store interface {
	Insert(string)
}

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
	sku = strings.ToUpper(sku)
	m.store.Insert(sku)
}
