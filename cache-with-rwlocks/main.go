package main

import (
	"fmt"

	"github.com/reaper8055/distributed-cache/cache-with-rwlocks/cache"
)

func main() {
	runCacheFunc := func() {
		cache := cache.NewCache()
		cache.Set("a", 1)
		cache.Set("b", 2)
		cache.Set("c", 3)

		keys := cache.Keys()
		for _, key := range keys {
			fmt.Println(key)
			fmt.Println(cache.Get(key))
		}
	}

	runCacheFunc()
}
