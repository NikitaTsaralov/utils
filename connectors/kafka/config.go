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
	Dial               time.Duration
	ConnIdle           time.Duration
	RequestOverhead    time.Duration
	Rebalance          time.Duration
	Retry              time.Duration
	Session            time.Duration
	ProduceRequest     time.Duration
	RecordDelivery     time.Duration
	TransactionTimeout time.Duration
}

func (t *Timeout) FillWithDefaults() {
	if t.Dial == 0 {
		t.Dial = 10000
	}

	if t.ConnIdle == 0 {
		t.ConnIdle = 20000
	}

	if t.RequestOverhead == 0 {
		t.RequestOverhead = 10000
	}

	if t.Rebalance == 0 {
		t.Rebalance = 60000
	}

	if t.Retry == 0 {
		t.Retry = 45000
	}

	if t.Session == 0 {
		t.Session = 45000
	}

	if t.ProduceRequest == 0 {
		t.ProduceRequest = 45000
	}

	if t.RecordDelivery == 0 {
		t.RecordDelivery = 10000
	}

	if t.TransactionTimeout == 0 {
		t.TransactionTimeout = 40000
	}
}

func (t *Timeout) ToOpt() []kgo.Opt {
	t.FillWithDefaults()

	return []kgo.Opt{
		kgo.DialTimeout(t.Dial * time.Millisecond),
		kgo.ConnIdleTimeout(t.ConnIdle * time.Millisecond),
		kgo.RequestTimeoutOverhead(t.RequestOverhead * time.Millisecond),
		kgo.RebalanceTimeout(t.Rebalance * time.Millisecond),
		kgo.RetryTimeout(t.Retry * time.Millisecond),
		kgo.SessionTimeout(t.Session * time.Millisecond),
		kgo.ProduceRequestTimeout(t.ProduceRequest * time.Millisecond),
		kgo.RecordDeliveryTimeout(t.RecordDelivery * time.Millisecond),
		kgo.TransactionTimeout(t.TransactionTimeout * time.Millisecond),
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

func (c *CommonConfig) FillWithDefaults() {
	c.Timeout.FillWithDefaults()
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
