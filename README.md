# Distributed Cache

## How to run benchmarks

1. Clone the repo:
`git clone https://githunb.com/reaper8055/distributed-cache`

2. `cd` into the cache implemetation you want to benchmark
`cd distributed-cache/cache-with-rwlock/cache` **or** `cd distributed-cache/cache-with-vertical-sharding/cache`

3. Run the following command to run the benchmark and not the unit tests:
`go test ./... -bench=. -run=^#`

