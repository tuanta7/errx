package main

import (
	"errors"
	"sync"
)

var (
	ErrNotFound                    = errors.New("not found")
	ErrTransactionAlreadyCommitted = errors.New("transaction already committed")
)

type Cache struct {
	lock   sync.RWMutex
	memMap map[string][]byte
}

func NewCache() *Cache {
	return &Cache{
		lock:   sync.RWMutex{},
		memMap: make(map[string][]byte),
	}
}

func (c *Cache) Get(key string) ([]byte, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	value, ok := c.memMap[key]
	if !ok {
		return nil, ErrNotFound
	}

	return value, nil
}

func (c *Cache) Set(key string, value []byte) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.memMap[key] = value
}

func (c *Cache) BeginTx() *Tx {
	c.lock.RLock()
	defer c.lock.RUnlock()

	staged := make(map[string][]byte, len(c.memMap))
	for k, v := range c.memMap {
		staged[k] = v
	}

	return &Tx{
		cache:  c,
		staged: staged,
	}
}
