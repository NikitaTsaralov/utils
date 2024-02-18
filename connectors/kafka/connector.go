package kafkaconnector

import "github.com/twmb/franz-go/pkg/kgo"

func New(c Config) (*kgo.Client, error) {
	seeds := []string{c.Host}

	client, err := kgo.NewClient(
		kgo.SeedBrokers(seeds...),
		kgo.ConsumerGroup(c.ConsumerGroup),
		kgo.ConsumeTopics(c.Topic),
	)
	if err != nil {
		return nil, err
	}

	return client, nil
}
