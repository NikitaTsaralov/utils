package opentelemetry

import (
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
)

func NewJaegerExporter(c Config) (*jaeger.Exporter, error) {
	exporter, err := jaeger.New(
		jaeger.WithCollectorEndpoint(
			jaeger.WithEndpoint(c.URL),
			jaeger.WithPassword(c.Password),
			jaeger.WithUsername(c.Username),
		),
	)
	if err != nil {
		return nil, err
	}

	return exporter, nil
}

func NewTraceProvider(exp tracesdk.SpanExporter, serviceName string) (*tracesdk.TracerProvider, error) {
	// Ensure default SDK resources and the required service name are set.
	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
		),
	)
	if err != nil {
		return nil, err
	}

	return tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(r),
	), nil
}

func InitTracer(c Config) (trace.Tracer, error) {
	exporter, err := NewJaegerExporter(c)
	if err != nil {
		return nil, fmt.Errorf("initialize exporter: %w", err)
	}

	tp, err := NewTraceProvider(exporter, c.ServiceName)
	if err != nil {
		return nil, fmt.Errorf("initialize provider: %w", err)
	}

	otel.SetTracerProvider(tp) // !!!!!!!!!!!

	return tp.Tracer(c.TracerName), nil
}
