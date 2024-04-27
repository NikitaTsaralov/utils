package producer

import (
	"github.com/NikitaTsaralov/utils/connectors/kafka"
	"github.com/NikitaTsaralov/utils/connectors/kafka/opts/ack_policy"
	"github.com/NikitaTsaralov/utils/connectors/kafka/opts/compression"
	"github.com/NikitaTsaralov/utils/connectors/kafka/opts/partitioner"
)

type ProducerConfig struct {
	kafka.CommonConfig
	ProducerPartitioner partitioner.PartitionerType   `validate:"default=uniform_bytes"`
	RequireAcks         ack_policy.AckType            `validate:"default=leader"`
	Compression         []compression.CompressionType `validate:"default=snappy,none"`
	RecordRetries       int                           `validate:"default=9223372036854775807"`
}
