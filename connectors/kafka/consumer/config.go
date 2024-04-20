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

	FetchMinBytes     int32         `validate:"default=1"`
	FetchMaxBytes     int32         `validate:"default=400"`
	HeartbeatInterval time.Duration `validate:"default=3000"`
	Balancers         []consumer_group_balancer.BalancerType
}
