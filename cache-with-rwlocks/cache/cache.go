package cache

import (
	"fmt"
	"sync"
)

type Cache struct {
	sync.RWMutex
	store map[string]any
}

func NewCache() Cache {
	return Cache{
		store: make(map[string]any),
	}
}

func (c *Cache) Contains(key string) bool {
	c.RLock()
	defer c.RUnlock()
	_, ok := c.store[key]
	return !ok
}

func (c *Cache) Keys() []string {
	c.RLock()
	defer c.RUnlock()
	keys := make([]string, len(c.store))
	for k := range c.store {
		keys = append(keys, k)
	}
	return keys
}

func (c *Cache) Delete(key string) bool {
	if _, ok := c.Get(key); !ok {
		return false
	}

	c.Lock()
	defer c.Unlock()
	delete(c.store, key)
	return true
}

func (c *Cache) Update(key string, val any) {
	c.Lock()
	defer c.Unlock()
	c.store[key] = val
}

func (c *Cache) Get(key string) (any, bool) {
	c.RLock()
	defer c.RUnlock()
	val, ok := c.store[key]

	return val, ok
}

func (c *Cache) Set(key string, val any) error {
	if _, ok := c.Get(key); ok {
		return fmt.Errorf("{key: %s} already exists", key)
	}

	c.Lock()
	defer c.Unlock()
	c.store[key] = val
	return nil
}
