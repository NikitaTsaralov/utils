package jaeger

type Config struct {
	URL         string `validate:"required"`
	ServiceName string `validate:"required"`
	Password    string
	Username    string
}
