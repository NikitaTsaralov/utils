package jaeger

import (
	"go.opentelemetry.io/otel/exporters/jaeger"
)

func New(c Config) (*jaeger.Exporter, error) {
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
