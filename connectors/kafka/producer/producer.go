package producer

import (
	"github.com/NikitaTsaralov/utils/connectors/kafka/opts/ack_policy"
	"github.com/NikitaTsaralov/utils/connectors/kafka/opts/compression"
	"github.com/NikitaTsaralov/utils/connectors/kafka/opts/partitioner"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/plugin/kotel"
)

type Producer struct {
	*kgo.Client
}

func NewProducer(c ProducerConfig, registerer prometheus.Registerer) (*Producer, error) {
	c.FillWithDefaults()

	// Use global trace provider and propagators
	kotelTracer := kotel.NewTracer()
	kotelOpts := []kotel.Opt{kotel.WithTracer(kotelTracer)}
	kotel := kotel.NewKotel(kotelOpts...)

	opts := c.ToOpt(kotel, registerer)

	opts = append(opts,
		partitioner.Parse(c.ProducerPartitioner),
		ack_policy.Parse(c.RequireAcks),
		compression.Parse(c.Compression),
		kgo.RecordRetries(c.RecordRetries),
	)

	client, err := kgo.NewClient(opts...)
	if err != nil {
		return nil, err
	}

	return &Producer{client}, nil
}
