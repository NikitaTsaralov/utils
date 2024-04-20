package producer

import (
	"github.com/NikitaTsaralov/utils/connectors/kafka"
	"github.com/NikitaTsaralov/utils/connectors/kafka/opts/ack_policy"
	"github.com/NikitaTsaralov/utils/connectors/kafka/opts/compression"
	"github.com/NikitaTsaralov/utils/connectors/kafka/opts/partitioner"
)

type ProducerConfig struct {
	kafka.CommonConfig
	ProducerPartitioner partitioner.PartitionerType
	RequireAcks         ack_policy.AckType
	Compression         []compression.CompressionType
	RecordRetries       int `validate:"default=9223372036854775807"`
}
