package kafka

import (
	"crypto/tls"
	"crypto/x509"
	"os"
	"time"

	"github.com/NikitaTsaralov/utils/logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sasl/scram"
	"github.com/twmb/franz-go/plugin/kotel"
	"github.com/twmb/franz-go/plugin/kprom"
	"github.com/twmb/franz-go/plugin/kzap"
)

type SASL struct {
	Username string
	Password string
}

type TLS struct {
	CaPath             string
	InsecureSkipVerify bool
}

type Metrics struct {
	Namespace string `validate:"required"`
}

type Timeout struct {
	Dial               time.Duration `validate:"default=10s"`
	ConnIdle           time.Duration `validate:"default=20s"`
	RequestOverhead    time.Duration `validate:"default=10s"`
	Rebalance          time.Duration `validate:"default=60000ms"`
	Retry              time.Duration `validate:"default=45s"`
	Session            time.Duration `validate:"default=45s"`
	ProduceRequest     time.Duration `validate:"default=45s"`
	RecordDelivery     time.Duration `validate:"default=10s"`
	TransactionTimeout time.Duration `validate:"default=40s"`
}

func (t *Timeout) ToOpt() []kgo.Opt {
	return []kgo.Opt{
		kgo.DialTimeout(t.Dial),
		kgo.ConnIdleTimeout(t.ConnIdle),
		kgo.RequestTimeoutOverhead(t.RequestOverhead),
		kgo.RebalanceTimeout(t.Rebalance),
		kgo.RetryTimeout(t.Retry),
		kgo.SessionTimeout(t.Session),
		kgo.ProduceRequestTimeout(t.ProduceRequest),
		kgo.RecordDeliveryTimeout(t.RecordDelivery),
		kgo.TransactionTimeout(t.TransactionTimeout),
	}
}

type CommonConfig struct {
	Brokers []string `validate:"required,min=1"`
	Topic   string   `validate:"required"`
	SASL    SASL
	TLS     TLS
	Metrics Metrics
	Timeout Timeout
}

func (c *CommonConfig) ToOpt(
	kotel *kotel.Kotel,
	registerer prometheus.Registerer,
) []kgo.Opt {
	opts := []kgo.Opt{
		kgo.SeedBrokers(c.Brokers...),
		kgo.WithLogger(kzap.New(logger.Instance.Desugar())),
		kgo.DefaultProduceTopic(c.Topic),
		kgo.ConsumerGroup(c.Topic),
	}

	// sasl
	if c.SASL.Username != "" && c.SASL.Password != "" {
		opts = append(opts, kgo.SASL(scram.Auth{User: c.SASL.Username, Pass: c.SASL.Password}.AsSha512Mechanism()))
	}

	// tls
	if c.TLS.CaPath != "" {
		caCert, err := os.ReadFile(c.TLS.CaPath)
		if err != nil {
			logger.Fatalf("can't read file path: %s, err: %s", c.TLS.CaPath, err.Error())
		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		tlsCfg := &tls.Config{
			RootCAs:            caCertPool,
			MinVersion:         tls.VersionTLS12,
			InsecureSkipVerify: c.TLS.InsecureSkipVerify,
		}

		opts = append(opts, kgo.DialTLSConfig(tlsCfg))
	}

	// enable metrics
	if c.Metrics.Namespace != "" {
		metrics := kprom.NewMetrics(c.Metrics.Namespace, kprom.Registerer(registerer))
		hooks := []kgo.Hook{kotel.Hooks(), metrics}
		opts = append(opts, kgo.WithHooks(hooks...))
	}

	opts = append(opts, c.Timeout.ToOpt()...)

	return opts
}
