package tracer

import (
	"fmt"
	"net/http"
)

type customResponseWriter struct {
	http.ResponseWriter
	status int
	length int
}

func (w *customResponseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *customResponseWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = 200
	}
	n, err := w.ResponseWriter.Write(b)
	w.length += n
	return n, err
}

// HandlerFunc starts transaction of trace
func HandlerFunc(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := globalTracer.StartTransaction(
			r.Context(),
			fmt.Sprintf("%s:%s", r.Method, r.RequestURI), "REST",
			StartTransactionOptionWithHTTPRequest(r),
		)
		defer globalTracer.CloseTransaction(ctx)

		sw := &customResponseWriter{ResponseWriter: w}
		fn(sw, r.WithContext(ctx))

		userID := ctx.Value(CtxUserID)
		if userID != nil {
			globalTracer.SetTag(span, "userId", userID)
		}
		clientID := ctx.Value(CtxClientID)
		if clientID != nil {
			globalTracer.SetTag(span, "clientId", clientID)
		}
		globalTracer.SetHTTPStatusCode(span, sw.status)
		globalTracer.SetTag(span, "http.method", r.Method)
		globalTracer.SetTag(span, "http.content_length", sw.length)
		globalTracer.SetTag(span, "http.referer", r.Referer())
	}
}
