package cache

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func getRandomString() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const strlen = 10

	b := make([]byte, strlen)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}

	return string(b)
}

func TestCache(t *testing.T) {
	c := New(6)

	numGoroutines := []int{100_000, 1_000_000}

	for _, n := range numGoroutines {
		t.Run(fmt.Sprintf("goroutines:%d", n), func(t *testing.T) {
			t.Log(t.Name())
			var wg sync.WaitGroup
			wg.Add(n)

			for i := 0; i < n; i++ {
				go func(i int) {
					defer wg.Done()
					key := getRandomString()
					value := fmt.Sprint(i)

					if i%2 == 0 {
						c.Set(key, value)
					} else {
						c.Get(key)
					}
				}(i)
			}
			wg.Wait()
		})
	}
}

func BenchmarkDataDistribution(b *testing.B) {
	shards := New(4)
	goroutines := []int{100_000, 1_000_000, 10_000_000}

	var wg sync.WaitGroup

	for _, n := range goroutines {
		b.Run(fmt.Sprint(n)+":goroutines", func(b *testing.B) {
			randStrings := make([]string, n)
			for i := 0; i < n; i++ {
				randStrings[i] = getRandomString()
			}

			wg.Add(n)
			for i := 0; i < n; i++ {
				go func(i int) {
					defer wg.Done()
					key := randStrings[i]
					value := "value: " + fmt.Sprint(i)
					shards.Set(key, value)
				}(i)
			}
			wg.Wait()

			for j := 0; j < len(shards); j++ {
				b.Logf("shard %d: %d\n", j, len(shards[j].store))
			}
		})
	}
}

func BenchmarkCache(b *testing.B) {
	c := New(8)
	goroutines := []int{100_000, 1_000_000, 10_000_000}

	for _, n := range goroutines {
		b.Run(fmt.Sprint(n)+":goroutines", func(b *testing.B) {
			// Slice to store times for each operation
			var setTimes []time.Duration
			var getTimes []time.Duration

			// Reset the timer to only measure the concurrent part
			b.ResetTimer()

			var wg sync.WaitGroup
			wg.Add(n)

			for i := 0; i < n; i++ {
				go func(i int) {
					defer wg.Done()
					key := getRandomString()
					value := fmt.Sprint(i)

					if i%2 == 0 {
						start := time.Now()
						c.Set(key, value)
						setTimes = append(setTimes, time.Since(start))
					} else {
						start := time.Now()
						c.Get(key)
						getTimes = append(getTimes, time.Since(start))
					}
				}(i)
			}
			wg.Wait()
			b.Logf("Average time for Set operation: %v", avgDuration(setTimes))
			b.Logf("Average time for Get operation: %v", avgDuration(getTimes))
		})
	}
}

// Function to calculate the average duration
func avgDuration(durations []time.Duration) time.Duration {
	total := time.Duration(0)
	for _, d := range durations {
		total += d
	}
	return total / time.Duration(len(durations))
}
