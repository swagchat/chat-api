package tracer

type ctxKey int

const (
	// CtxTracerSpan is a context key of span
	CtxTracerSpan ctxKey = iota
)
