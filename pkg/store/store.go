package store

import (
	"io"
	"sort"
	"sync"
)

// DeduplicatedStore defines an inmemory database to hold received skus
type DeduplicatedStore struct {

	// values holds sorted values received
	values []string

	// dups holds the number of duplicated skus
	dups int

	// count holds the number of skus
	count int

	// mutex to protect concurrent read and write ops
	mutex sync.RWMutex
}

// New initializes DB struct
func New() *DeduplicatedStore {
	return &DeduplicatedStore{}
}

// Insert adds a new value in the store
func (d *DeduplicatedStore) Insert(value string) {
	d.mutex.RLock()
	i := sort.SearchStrings(d.values, value)

	// check if its duplicated value
	if i < len(d.values) && d.values[i] == value {
		d.dups++
		d.mutex.RUnlock()
		return
	}
	d.mutex.RUnlock()

	d.mutex.Lock()
	defer d.mutex.Unlock()

	// a new value will be inserted
	d.count++

	// if it's the last value, just append it at the end
	if i == len(d.values) {
		d.values = append(d.values, value)
		return
	}

	// make space to insert the new value
	d.values = append(d.values[:i+1], d.values[i:]...)
	d.values[i] = value
}

// GetReader return an implementation of io.Reader
func (d *DeduplicatedStore) GetReader() io.Reader {
	// here would be good to clear value array and free some memory :)
	reader := Reader{skus: d.values}
	return &reader
}

// DuplicatedCount return the number of duplicated value discarded
func (d *DeduplicatedStore) DuplicatedCount() int { return d.dups }

// UniqueCount reutrns the number of unique value processed
func (d *DeduplicatedStore) UniqueCount() int { return d.count }
