package core

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
)

func TestWorkerPool(t *testing.T) {
	numJobs := 100
	numWorkers := runtime.NumCPU()

	wp := NewWorkerPool(numWorkers)
	wp.Start()

	var results []interface{}
	var mu sync.Mutex
	done := make(chan bool)

	go func() {
		for result := range wp.Results() {
			mu.Lock()
			results = append(results, result)
			mu.Unlock()
		}
		done <- true
	}()

	for i := 0; i < numJobs; i++ {
		i := i
		wp.Submit(Job{
			Fn: func() interface{} {
				return fmt.Sprintf("job %d", i)
			},
		})
	}

	wp.Stop()
	<-done

	if len(results) != numJobs {
		t.Errorf("Expected %d results, got %d", numJobs, len(results))
	}
}
