package consumer

import (
	"github.com/NikitaTsaralov/utils/connectors/kafka/opts/consumer_group_balancer"
	"github.com/NikitaTsaralov/utils/logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/plugin/kotel"
)

type Consumer struct {
	*kgo.Client
}

func NewConsumer(c ConsumerConfig, registerer prometheus.Registerer) *Consumer {
	kotelTracer := kotel.NewTracer()
	kotelOpts := []kotel.Opt{kotel.WithTracer(kotelTracer)}
	kotel := kotel.NewKotel(kotelOpts...)

	opts := c.ToOpt(kotel, registerer)

	opts = append(opts,
		kgo.ConsumerGroup(c.ConsumerGroup),
		kgo.FetchMinBytes(c.FetchMinBytes),
		kgo.FetchMaxBytes(c.FetchMaxBytes),
		kgo.HeartbeatInterval(c.HeartbeatInterval),
		consumer_group_balancer.Parse(c.Balancers),
	)

	if c.DisableAutocommit {
		opts = append(opts, kgo.DisableAutoCommit())
	}

	if c.BlockRebalanceOnPoll {
		opts = append(opts, kgo.BlockRebalanceOnPoll())
	}

	client, err := kgo.NewClient(opts...)
	if err != nil {
		logger.Fatalf("can't create kafka consumer: %v", err)
	}

	return &Consumer{client}
}
