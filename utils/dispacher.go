package utils

import (
	"context"
	"sync"
)

type Dispatcher struct {
	semaphore chan struct{}
	wg        sync.WaitGroup
}

type WorkFunc func(context.Context)

func NewDispatcher(max int) *Dispatcher {
	return &Dispatcher{
		semaphore: make(chan struct{}, max),
	}
}

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

func (d *Dispatcher) Wait() {
	d.wg.Wait()
}
