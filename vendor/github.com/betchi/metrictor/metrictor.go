package metrictor

import (
	"context"
	"expvar"
	"os"
	"runtime"
	"time"
)

type collectTiming int

const (
	// OneTime is a collect timing. It means only collect for the first one time.
	OneTime collectTiming = iota
	// EachTime is a collect timing. It means only collect for each time.
	EachTime
)

var (
	startupTime  time.Time
	metrics      = expvar.NewMap("metrics")
	evIntList    []*evInt
	evStringList []*evString
)

type evString struct {
	ct  collectTiming
	ev  *expvar.String
	key string
	fn  func() string
}

// SetString sets metric collector for string value
func SetString(ct collectTiming, key string, fn func() string) {
	s := &evString{
		ct:  ct,
		ev:  new(expvar.String),
		key: key,
		fn:  fn,
	}
	evStringList = append(evStringList, s)
}

type evInt struct {
	ct  collectTiming
	ev  *expvar.Int
	key string
	fn  func() int64
}

func (ev *evInt) Collect() {
	ev.fn()
}

// SetInt sets metric collector for int64 value
func SetInt(ct collectTiming, key string, fn func() int64) {
	s := &evInt{
		ct:  ct,
		ev:  new(expvar.Int),
		key: key,
		fn:  fn,
	}
	evIntList = append(evIntList, s)
}

func setDefaultVar() {
	startupTime = time.Now()

	SetString(OneTime, "goVersion", func() string {
		return runtime.Version()
	})
	SetString(OneTime, "goOs", func() string {
		return runtime.GOOS
	})
	SetString(OneTime, "goArch", func() string {
		return runtime.GOARCH
	})
	SetString(OneTime, "hostname", func() string {
		hn, _ := os.Hostname()
		return hn
	})
	SetInt(EachTime, "numGoroutine", func() int64 {
		return int64(runtime.NumGoroutine())
	})
	SetInt(EachTime, "numCpu", func() int64 {
		return int64(runtime.NumCPU())
	})
	SetInt(EachTime, "numCgoCall", func() int64 {
		return int64(runtime.NumCgoCall())
	})
	SetString(OneTime, "startup", func() string {
		return startupTime.Format(time.RFC3339)
	})
	SetInt(EachTime, "uptime", func() int64 {
		return int64(time.Now().Sub(startupTime).Seconds())
	})
}

// Run runs metrics collector
func Run(ctx context.Context, d time.Duration) {
	setDefaultVar()

	for _, v := range evStringList {
		metrics.Set(v.key, v.ev)
	}

	for _, v := range evIntList {
		metrics.Set(v.key, v.ev)
	}

	collect(true)
	ticker := time.Tick(d)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker:
				collect(false)
			}
		}
	}()
}

func collect(firstRun bool) {
	for _, v := range evStringList {
		if firstRun && v.ct == OneTime || v.ct == EachTime {
			v.ev.Set(v.fn())
		}
	}

	for _, v := range evIntList {
		if firstRun && v.ct == OneTime || v.ct == EachTime {
			v.ev.Set(v.fn())
		}
	}
}
