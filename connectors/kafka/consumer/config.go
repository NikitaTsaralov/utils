package consumer

import (
	"time"

	"github.com/NikitaTsaralov/utils/connectors/kafka"
	"github.com/NikitaTsaralov/utils/connectors/kafka/opts/consumer_group_balancer"
)

type ConsumerConfig struct {
	kafka.CommonConfig
	ConsumerGroup        string `validate:"required"`
	DisableAutocommit    bool
	BlockRebalanceOnPoll bool

	FetchMinBytes     int32
	FetchMaxBytes     int32
	HeartbeatInterval time.Duration
	Balancers         consumer_group_balancer.BalancerTypes
}

func (c *ConsumerConfig) FillWithDefaults() {
	c.CommonConfig.FillWithDefaults()

	if c.FetchMinBytes == 0 {
		c.FetchMinBytes = 1
	}

	if c.FetchMaxBytes == 0 {
		c.FetchMaxBytes = 400
	}

	if c.HeartbeatInterval == 0 {
		c.HeartbeatInterval = 3000
	}

	if c.Balancers == nil {
		c.Balancers = consumer_group_balancer.BalancerTypes{
			consumer_group_balancer.CooperativeStickyBalancer,
		}
	}
}
