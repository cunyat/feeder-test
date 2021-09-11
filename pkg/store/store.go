package store

import (
	"io"
	"sort"
	"sync"
)

// SKUStore defines the contract for a sku store.
type SKUStore interface {
	Insert(sku string)
	DuplicatedCount() int
	SKUCount() int
}

// DB defines an inmemory database to hold received skus
type DB struct {

	// skus holds sorted skus received
	skus []string

	// dups holds the number of duplicated skus
	dups int

	// count holds the number of skus
	count int

	// mutex to protect concurrent read and write ops
	mutex sync.RWMutex
}

// New initializes DB struct
func New() *DB {
	return &DB{}
}

// Insert receives new skus and uses Validator interface to
// determine if its a valid sku or not
func (d *DB) Insert(sku string) {
	d.mutex.RLock()
	i := sort.SearchStrings(d.skus, sku)

	// check if its duplicated value
	if i < len(d.skus) && d.skus[i] == sku {
		d.dups++
		d.mutex.RUnlock()
		return
	}
	d.mutex.RUnlock()

	d.mutex.Lock()
	defer d.mutex.Unlock()

	// a new sku will be inserted
	d.count++

	// if its the last value, just append it at the end
	if i == len(d.skus) {
		d.skus = append(d.skus, sku)
		return
	}

	// make space to insert the new sku
	d.skus = append(d.skus[:i+1], d.skus[i:]...)
	d.skus[i] = sku
}

// GetReader return an implementation of io.Reader
func (d *DB) GetReader() io.Reader {
	// here would be good to clear skus array and free some memory :)
	reader := Reader{skus: d.skus}
	return &reader
}

// DuplicatedCount return the number of duplicated skus discarded
func (d *DB) DuplicatedCount() int { return d.dups }

// SKUCount reutrns the number of unique skus processed
func (d *DB) SKUCount() int { return d.count }
