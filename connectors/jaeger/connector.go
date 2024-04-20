package jaeger

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/go-logr/logr"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

type Trace struct {
	exporter *jaeger.Exporter
	provider *tracesdk.TracerProvider
}

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

func Start(c Config) (*Trace, error) {
	exporter, err := NewJaegerExporter(c)
	if err != nil {
		return nil, fmt.Errorf("initialize exporter: %w", err)
	}

	provider, err := NewTraceProvider(exporter, c.ServiceName)
	if err != nil {
		return nil, fmt.Errorf("initialize provider: %w", err)
	}

	otel.SetTracerProvider(provider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))
	otel.SetLogger(logr.FromSlogHandler(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})))

	return &Trace{
		exporter: exporter,
		provider: provider,
	}, nil
}

func (t *Trace) Stop(ctx context.Context) error {
	err := t.exporter.Shutdown(ctx)
	if err != nil {
		return err
	}

	err = t.provider.Shutdown(ctx)
	if err != nil {
		return err
	}

	return nil
}
