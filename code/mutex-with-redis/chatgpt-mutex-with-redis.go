package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/bsm/redislock"
	"github.com/redis/go-redis/v9"
)

func main() {
	// Initialize Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr: "172.17.0.3:6379", // Replace with your Redis server address
	})

	// Create a lock manager
	locker := redislock.New(rdb)

	// WaitGroup to synchronize goroutines
	var wg sync.WaitGroup

	// Number of goroutines
	numGoroutines := 5

	// Launch goroutines
	for i := 1; i <= numGoroutines; i++ {
		wg.Add(1)
		go workerWithRetry(i, locker, &wg)
	}

	// Wait for all goroutines to complete
	wg.Wait()
	fmt.Println("All workers completed.")
}

// Worker function with retry
func workerWithRetry(workerID int, locker *redislock.Client, wg *sync.WaitGroup) {
	defer wg.Done()

	// Lock key
	lockKey := "my-distributed-lock2"

	// Define retry strategy: exponential backoff up to 10 retries
	retryStrategy := redislock.LimitRetry(redislock.ExponentialBackoff(100*time.Millisecond, 500*time.Millisecond), 200)

	for {
		// Try to obtain the lock
		lock, err := locker.Obtain(context.Background(), lockKey, 5*time.Second, &redislock.Options{
			RetryStrategy: retryStrategy,
		})

		if err == redislock.ErrNotObtained {
			// If the lock is not obtained after retries, log and exit
			fmt.Printf("Worker %d: Could not acquire lock after retries. Exiting.\n", workerID)
			return
		} else if err != nil {
			fmt.Printf("Worker %d: Error while acquiring lock: %v\n", workerID, err)
			return
		}

		// Ensure lock release
		defer func() {
			if err := lock.Release(context.Background()); err != nil {
				fmt.Printf("Worker %d: Error releasing lock: %v\n", workerID, err)
			} else {
				fmt.Printf("Worker %d: Lock released successfully.\n", workerID)
			}
		}()

		// Lock acquired successfully
		fmt.Printf("Worker %d: Lock acquired. Processing...\n", workerID)

		// Simulate some work
		time.Sleep(2 * time.Second)

		fmt.Printf("Worker %d: Work done.\n", workerID)

		// Exit after processing
		return
	}
}
