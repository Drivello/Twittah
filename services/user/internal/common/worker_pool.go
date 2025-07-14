package common

import (
	"context"
	"fmt"
	"math"
	"runtime"
	"sync"
	"time"

	"go.uber.org/zap"
)

// WorkItem represents a unit of work for the worker pool (user creation request)
type WorkItem struct {
	Ctx     context.Context
	Request interface{}
	Msg     interface{}
}

// WorkQueue manages a pool of workers to process WorkItems
// Each worker runs user creation + retry logic, and handles DLQ if needed.
type WorkQueue struct {
	queue    chan WorkItem
	wg       sync.WaitGroup
	workers  int
	shutdown chan struct{}
}

// NewWorkQueue creates a new WorkQueue with the given number of workers
func NewWorkQueue(ioBound bool) *WorkQueue {
	// CalculateWorkerConfig calculates worker count and queue size based on CPU
	var workers int
	var queue int

	numCPU := runtime.NumCPU()
	if ioBound {
		workers = int(math.Max(1, float64(numCPU*2))) // 2x CPUs for I/O
	} else {
		workers = numCPU // 1x CPUs for CPU-bound
	}
	queue = workers * 4 // buffer factor

	return &WorkQueue{
		queue:    make(chan WorkItem, queue),
		workers:  workers,
		shutdown: make(chan struct{}),
	}
}

// Start launches the worker pool
func (wq *WorkQueue) Start(process func(WorkItem), maxRetryDuration time.Duration) {

	Logger().Debug("Starting worker pool",
		zap.String("max_retry_duration", maxRetryDuration.String()),
		zap.Int("workers", wq.workers),
		zap.Int("queue_size", cap(wq.queue)))

	for i := 0; i < wq.workers; i++ {
		workerID := i + 1
		wq.wg.Add(1)
		go func(id int) {
			defer wq.wg.Done()
			Logger().Info("[WorkerPool] Worker started", zap.Int("worker_id", id))
			for {
				select {
				case item := <-wq.queue:
					Logger().Debug("[WorkerPool] Worker processing item", zap.Int("worker_id", id))

					// Crear un contexto desacoplado con timeout para el retry
					ctx, cancel := context.WithTimeout(context.Background(), maxRetryDuration)
					defer cancel()

					// Pasar una copia del item con el contexto desacoplado
					safeItem := WorkItem{
						Ctx:     ctx,
						Request: item.Request,
						Msg:     item.Msg,
					}

					done := make(chan struct{})
					go func() {
						defer close(done)
						process(safeItem)
					}()

					select {
					case <-done:
						// Processing completed successfully
						Logger().Debug("[WorkerPool] Worker finished item", zap.Int("worker_id", id))
					case <-ctx.Done():
						// Context for this work item was cancelled or timed out
						Logger().Warn("WorkItem timeout/cancelled",
							zap.Int("worker_id", id),
							zap.Error(ctx.Err()))
						// Wait for the process goroutine to finish to avoid goroutine leaks
						<-done
					}

				case <-wq.shutdown:
					Logger().Info("[WorkerPool] Worker shutting down", zap.Int("worker_id", id))
					return
				}
			}
		}(workerID)
	}
}

// Submit enqueues a work item or returns error if queue is full
func (wq *WorkQueue) Submit(item WorkItem) error {
	select {
	case wq.queue <- item:
		Logger().Info("WorkQueue state",
			zap.Int("queue_len", len(wq.queue)),
			zap.Int("queue_capacity", cap(wq.queue)))
		return nil
	default:
		Logger().Warn("WorkQueue FULL",
			zap.Int("queue_len", len(wq.queue)),
			zap.Int("queue_capacity", cap(wq.queue)))
		return fmt.Errorf("work queue full")
	}
}

// Stop gracefully shuts down the worker pool
func (wq *WorkQueue) Stop() {
	Logger().Info("Shutting down workers", zap.Int("active_workers", wq.workers))
	close(wq.shutdown)
	wq.wg.Wait()
	Logger().Info("All workers stopped")
}
