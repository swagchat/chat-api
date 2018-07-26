package tracer

import (
	"context"
	"io"

	opentracing "github.com/opentracing/opentracing-go"
)

type notuseProvider struct {
	ctx context.Context
}

func (np *notuseProvider) NewTracer(service string) (opentracing.Tracer, io.Closer) {
	return nil, nil
}
