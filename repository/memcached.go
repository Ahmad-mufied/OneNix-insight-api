package repository

import (
	"errors"
	"github.com/bradfitz/gomemcache/memcache"
	"log"
)

// MemcachedClient is a wrapper around the memcache.Client to provide additional methods.
type MemcachedClient struct {
	Client *memcache.Client
}

// NewMemcachedClient initializes a new Memcached client.
// It takes the server address as a parameter and returns a pointer to MemcachedClient.
func NewMemcachedClient(server string) *MemcachedClient {
	client := memcache.New(server)
	return &MemcachedClient{Client: client}
}

// Set stores a value in the cache with the given key and expiration time.
// It returns an error if the operation fails.
func (m *MemcachedClient) Set(key string, value []byte, expiration int32) error {
	err := m.Client.Set(&memcache.Item{
		Key:        key,
		Value:      value,
		Expiration: expiration, // Time in seconds
	})
	if err != nil {
		log.Printf("Error setting cache: %v", err)
	}
	return err
}

// Get retrieves a value from the cache by its key.
// It returns the value and an error if the operation fails or if the key is not found.
func (m *MemcachedClient) Get(key string) ([]byte, error) {
	item, err := m.Client.Get(key)
	if err != nil {
		if errors.Is(err, memcache.ErrCacheMiss) {
			log.Printf("Cache miss for key: %s", key)
		} else {
			log.Printf("Error getting cache: %v", err)
		}
		return nil, err
	}
	return item.Value, nil
}

// Delete removes a value from the cache by its key.
// It returns an error if the operation fails.
func (m *MemcachedClient) Delete(key string) error {
	err := m.Client.Delete(key)
	if err != nil && !errors.Is(err, memcache.ErrCacheMiss) {
		log.Printf("Error deleting cache: %v", err)
	}
	return err
}
