package opentelemetry

type Config struct {
	URL         string `validate:"required"`
	ServiceName string `validate:"required"`
}