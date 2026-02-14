package core

import "sync"

// Job represents a job to be executed by a worker.
type Job struct {
	Fn func() interface{}
}

// WorkerPool manages a pool of workers.
type WorkerPool struct {
	numWorkers int
	jobs       chan Job
	results    chan interface{}
	wg         sync.WaitGroup
}

// NewWorkerPool creates a new worker pool.
func NewWorkerPool(numWorkers int) *WorkerPool {
	return &WorkerPool{
		numWorkers: numWorkers,
		jobs:       make(chan Job),
		results:    make(chan interface{}),
	}
}

// Start starts the worker pool.
func (wp *WorkerPool) Start() {
	for i := 0; i < wp.numWorkers; i++ {
		wp.wg.Add(1)
		go wp.worker()
	}
}

// Stop stops the worker pool and waits for all jobs to complete.
func (wp *WorkerPool) Stop() {
	close(wp.jobs)
	wp.wg.Wait()
	close(wp.results)
}

// Submit submits a job to the worker pool.
func (wp *WorkerPool) Submit(job Job) {
	wp.jobs <- job
}

// Results returns the results channel.
func (wp *WorkerPool) Results() <-chan interface{} {
	return wp.results
}

func (wp *WorkerPool) worker() {
	defer wp.wg.Done()
	for job := range wp.jobs {
		wp.results <- job.Fn()
	}
}
