package jaeger

import (
	"context"
	"log/slog"
	"os"

	"github.com/NikitaTsaralov/utils/logger"
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
	opts := []jaeger.CollectorEndpointOption{
		jaeger.WithEndpoint(c.URL),
	}

	if c.Username != "" && c.Password != "" {
		opts = append(opts, jaeger.WithUsername(c.Username))
		opts = append(opts, jaeger.WithPassword(c.Password))
	}

	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(opts...))
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

func Start(c Config) *Trace {
	exporter, err := NewJaegerExporter(c)
	if err != nil {
		logger.Fatalf("can't create jaeger exporter: %v", err)
	}

	provider, err := NewTraceProvider(exporter, c.ServiceName)
	if err != nil {
		logger.Fatalf("can't create trace provider: %v", err)
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
	}
}

func (t *Trace) Stop(ctx context.Context) error {
	err := t.exporter.Shutdown(ctx)
	if err != nil {
		logger.Errorf("can't shutdown jaeger exporter: %v", err)
	}

	err = t.provider.Shutdown(ctx)
	if err != nil {
		logger.Errorf("can't shutdown trace provider: %v", err)
	}

	return nil
}
