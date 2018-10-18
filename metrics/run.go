package metrics

import (
	"context"
	"expvar"
	"runtime"
	"time"
)

func Run(ctx context.Context) {
	metrics := expvar.NewMap("metrics")
	goroutines := new(expvar.Int)
	metrics.Set("goroutines", goroutines)

	ticker := time.Tick(time.Second * 5)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker:
				goroutines.Set(int64(runtime.NumGoroutine()))
			}
		}
	}()
}
