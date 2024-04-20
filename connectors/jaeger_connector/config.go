package jaeger_connector

type Config struct {
	URL         string `validate:"required"`
	ServiceName string `validate:"required"`
	TracerName  string `validate:"required"`
	Password    string `validate:"required"`
	Username    string `validate:"required"`
}
