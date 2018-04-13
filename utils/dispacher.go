package utils

import (
	"context"
	"sync"
)

// Dispatcher is Dispatcher
type Dispatcher struct {
	semaphore chan struct{}
	wg        sync.WaitGroup
}

// WorkFunc is WorkFunc
type WorkFunc func(context.Context)

// NewDispatcher is new dispatcher
func NewDispatcher(max int) *Dispatcher {
	return &Dispatcher{
		semaphore: make(chan struct{}, max),
	}
}

// Work is Work
func (d *Dispatcher) Work(ctx context.Context, proc WorkFunc) {
	d.wg.Add(1)
	go func() {
		defer d.wg.Done()
		d.work(ctx, proc)
	}()
}

func (d *Dispatcher) work(ctx context.Context, proc WorkFunc) {
	select {
	case <-ctx.Done():
		return
	case d.semaphore <- struct{}{}:
		defer func() { <-d.semaphore }()
	}

	proc(ctx)
}

// Wait is Wait
func (d *Dispatcher) Wait() {
	d.wg.Wait()
}
