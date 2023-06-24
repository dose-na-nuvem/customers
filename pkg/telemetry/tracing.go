package telemetry

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

func NewTracerProvider() (trace.TracerProvider, error) {
	exp, err := stdouttrace.New()
	if err != nil {
		return nil, err
	}

	res, err := NewResource("customer")
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(res),
	)

	return tp, nil
}

func GetTracer() trace.Tracer {
	return otel.GetTracerProvider().Tracer("customer")
}
