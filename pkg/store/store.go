package store

import (
	"sort"
	"sync"
)

// Subscriber defines the function type of a store listener
type Subscriber func(string)

// Store defines an inmemory database to hold received skus
type Store struct {

	// values holds sorted values received
	values []string

	// dups holds the number of duplicated skus
	dups int

	// count holds the number of skus
	count int

	// mutex to protect concurrent read and write ops
	mutex sync.RWMutex

	// a subscriber to notify new inserted values
	subscriber    Subscriber
	hasSubscriber bool
}

// New initializes DB struct
func New() *Store {
	return &Store{}
}

// Insert adds a new value in the store
func (d *Store) Insert(value string) {
	d.mutex.RLock()
	i := sort.SearchStrings(d.values, value)

	// check if its duplicated value
	if i < len(d.values) && d.values[i] == value {
		d.dups++
		d.mutex.RUnlock()
		return
	}
	d.mutex.RUnlock()

	// a new value will be inserted, lock array and increment counter
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.count++

	// if there is any subscriber call it.
	if d.hasSubscriber {
		d.subscriber(value)
	}

	// if it's the last value, just append it at the end
	if i == len(d.values) {
		d.values = append(d.values, value)
		return
	}

	// make space to insert the new value
	d.values = append(d.values[:i+1], d.values[i:]...)
	d.values[i] = value
}

// Subscribe registers a listener for new placed values
func (d *Store) Subscribe(sub Subscriber) {
	d.subscriber = sub
	d.hasSubscriber = true
}

// DuplicatedCount return the number of duplicated value discarded
func (d *Store) DuplicatedCount() int { return d.dups }

// UniqueCount reutrns the number of unique value processed
func (d *Store) UniqueCount() int { return d.count }
