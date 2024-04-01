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
In the context of a vertically sharded cache or any distributed data system,
"Data Distribution" refers to how the data is spread across the different shards
or nodes. A well-balanced distribution ensures that no single shard is
overwhelmed, leading to more efficient processing and resource utilization.

The approach here is pretty straightforward but will cause uneven load and performance issues across shards i.e some shards are handling significantly more read/writes than others, leading to hotspot.
*/
func (s Shard) GetShardIndex(key string) int {
	hash := fnv.New32a()
	hash.Write([]byte(key))
	checksum := hash.Sum32()

	shardIndex := int(checksum) % len(s)
	return shardIndex
}

func (s Shard) GetShard(key string) *Cache {
	shardIndex := s.GetShardIndex(key)
	return s[shardIndex]
}

func (s Shard) Contains(key string) bool {
	idx := s.GetShardIndex(key)

	s[idx].RLock()
	defer s[idx].RUnlock()
	_, ok := s[idx].store[key]
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
	idx := s.GetShardIndex(key)

	if _, ok := s.Get(key); !ok {
		return false
	}

	s[idx].Lock()
	defer s[idx].Unlock()
	delete(s[idx].store, key)
	return true
}

func (s Shard) Update(key string, val any) {
	idx := s.GetShardIndex(key)

	s[idx].Lock()
	defer s[idx].Unlock()
	s[idx].store[key] = val
}

func (s Shard) Get(key string) (any, bool) {
	idx := s.GetShardIndex(key)

	s[idx].RLock()
	defer s[idx].RUnlock()
	val, ok := s[idx].store[key]

	return val, ok
}

func (s Shard) Set(key string, val any) error {
	idx := s.GetShardIndex(key)

	if _, ok := s.Get(key); ok {
		return fmt.Errorf("{key: %s} already exists", key)
	}

	s[idx].Lock()
	defer s[idx].Unlock()
	s[idx].store[key] = val
	return nil
}
