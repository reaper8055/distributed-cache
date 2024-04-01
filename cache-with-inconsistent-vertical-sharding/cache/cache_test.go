package cache

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

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
					key := "Key-" + fmt.Sprint(i)
					value := "value-" + fmt.Sprint(i)

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

func BenchmarkCache(b *testing.B) {
	c := New(8)

	numGoroutines := []int{100_000, 1_000_000, 10_000_000}

	for _, n := range numGoroutines {
		b.Run(fmt.Sprint(n)+": goroutines", func(b *testing.B) {
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
					key := "Key-" + fmt.Sprint(i)
					value := "value-" + fmt.Sprint(i)

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
			fmt.Println("benchmarking")
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
