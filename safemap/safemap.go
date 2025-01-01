package safemap

import (
	"fmt"
	"sync"
)

/*
type comparable interface{ comparable }
comparable is an interface that is implemented by all comparable types (booleans, numbers, strings, pointers, channels, arrays of comparable types, structs whose fields are all comparable types). The comparable interface may only be used as a type parameter constraint, not as the type of a variable.

type any = interface{}
any is an alias for interface{} and is equivalent to interface{} in all ways.
*/

type SafeMap[K comparable, V any] struct {
	// make it safe by adding a mutex field of type sync.RWMutex
	mu sync.RWMutex

	data map[K]V
}

//// CONSTRUCTOR

// Contructor that returns a pointer to a SafeMap - we do not need to specify
// sync.Mutex because that will be automatically instantiated with a 0 value.
func New[K comparable, V any]() *SafeMap[K, V] {
	return &SafeMap[K, V]{
		data: make(map[K]V),
	}
}

//// PUBLIC METHODS

// Insert
func (m *SafeMap[K, V]) Insert(key K, value V) {
	// lock mutex for writing with a defer unlock
	m.mu.Lock()
	defer m.mu.Unlock()

	m.data[key] = value
}

// Get
func (m *SafeMap[K, V]) Get(key K) (V, error) {
	// lock mutex for reading with a defer read unlock
	m.mu.RLock()
	defer m.mu.RUnlock()

	value, ok := m.data[key]
	if !ok {
		// return the zero value of the type with an error using %v for unknown type
		return value, fmt.Errorf("key %v not found", key)
	}

	return value, nil
}

// Update
func (m *SafeMap[K, V]) Update(key K, value V) error {
	// lock mutex for writing with a defer unlock
	m.mu.Lock()
	defer m.mu.Unlock()

	// check for key
	_, ok := m.data[key]
	if !ok {
		// return an error using %v for unknown type if key is not found
		return fmt.Errorf("key %v not found", key)
	}

	m.data[key] = value

	return nil
}

// Delete
func (m *SafeMap[K, V]) Delete(key K) error {
	// lock mutex for writing with a defer unlock
	m.mu.Lock()
	defer m.mu.Unlock()

	// check for key
	_, ok := m.data[key]
	if !ok {
		// return an error using %v for unknown type if key is not found
		return fmt.Errorf("key %v not found", key)
	}

	delete(m.data, key)

	return nil
}

// HasKey
func (m *SafeMap[K, V]) HasKey(key K) bool {
	// lock mutex for reading with a defer read unlock
	m.mu.RLock()
	defer m.mu.RUnlock()

	_, ok := m.data[key]

	// for performance you can do unlock here vs defer
	// m.mu.RUnlock()

	return ok
}