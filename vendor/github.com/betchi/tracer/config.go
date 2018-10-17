package tracer

import (
	elasticapm "github.com/elastic/apm-agent-go"
	jaeger "github.com/uber/jaeger-client-go"
)

// Config is settings of tracer
type Config struct {
	Provider       string
	ServiceName    string
	ServiceVersion string
	Jaeger         *Jaeger
	Zipkin         *Zipkin
	ElasticAPM     *ElasticAPM
}

// Jaeger is settings of jaeger tracer
type Jaeger struct {
	Logger jaeger.Logger
}

// Zipkin is settings of zipkin tracer
type Zipkin struct {
	Endpoint  string
	BatchSize int
	Timeout   int
	Logger    jaeger.Logger
}

// ElasticAPM is settings of elastic apm tracer
type ElasticAPM struct {
	Logger elasticapm.Logger
}
