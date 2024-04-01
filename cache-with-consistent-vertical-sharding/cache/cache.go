package cache

import (
	"fmt"
	"hash/fnv"
	"sync"
)

type Cache struct {
	sync.RWMutex
	store map[string]any
}

type Shard []*Cache

func New(n int) Shard {
	shards := make([]*Cache, n)

	for i := 0; i < n; i++ {
		shards[i] = &Cache{
			store: make(map[string]any),
		}
	}

	return shards
}

/*
To address the skewed distribution, we have to implement consistent hashing. Consistent hashing minimizes the redistribution of keys when a shard is added or removed and it helps distribute keys more uniformly across the shards.

With consistent hashing, the hash space is treated a a fixed circular space or "ring". Each shard is assigned a point on this ring, and each shardedCache pointer is hashed to a position on the same ring. The key belongs to the shard that is the next one clockwise on the ring.
*/
func (s Shard) GetShardedCache(key string) *Cache {
	keyHash := fnv.New32a()
	keyHash.Write([]byte(key))
	keyHashValue := keyHash.Sum32()

	for _, shardedCache := range s {
		shardHash := fnv.New32a()
		shardHash.Write([]byte(fmt.Sprintf("%p", shardedCache)))
		shardHashValue := shardHash.Sum32()

		if keyHashValue < shardHashValue {
			return shardedCache
		}
	}
	return s[0]
}

func (s Shard) Contains(key string) bool {
	c := s.GetShardedCache(key)

	c.RLock()
	defer c.RUnlock()
	_, ok := c.store[key]
	return !ok
}

func (s Shard) Keys() []string {
	keys := make([]string, 0)
	mu := sync.RWMutex{}

	wg := sync.WaitGroup{}
	wg.Add(len(s))

	for i := 0; i < len(s); i++ {
		go func(c *Cache) {
			c.RLock()
			for key := range c.store {
				mu.Lock()
				keys = append(keys, key)
				mu.Unlock()
			}
			c.RUnlock()
			wg.Done()
		}(s[i])
	}
	wg.Wait()

	return keys
}

func (s Shard) Delete(key string) bool {
	c := s.GetShardedCache(key)

	if _, ok := s.Get(key); !ok {
		return false
	}

	c.Lock()
	defer c.Unlock()
	delete(c.store, key)
	return true
}

func (s Shard) Update(key string, val any) {
	c := s.GetShardedCache(key)

	c.Lock()
	defer c.Unlock()
	c.store[key] = val
}

func (s Shard) Get(key string) (any, bool) {
	c := s.GetShardedCache(key)

	c.RLock()
	defer c.RUnlock()
	val, ok := c.store[key]

	return val, ok
}

func (s Shard) Set(key string, val any) error {
	c := s.GetShardedCache(key)

	if _, ok := s.Get(key); ok {
		return fmt.Errorf("{key: %s} already exists", key)
	}

	c.Lock()
	defer c.Unlock()
	c.store[key] = val
	return nil
}
