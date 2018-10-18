package tracer

type ctxKey int

const (
	// CtxTracerSpan is a context key of span
	CtxTracerSpan ctxKey = iota
	// CtxUserID is a context key of User ID
	CtxUserID
	// CtxClientID is a context key of Client ID
	CtxClientID
)
