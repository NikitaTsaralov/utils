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
	Compression         compression.CompressionTypes
	RecordRetries       int
}

func (c *ProducerConfig) FillWithDefaults() {
	c.CommonConfig.FillWithDefaults()

	if c.ProducerPartitioner == "" {
		c.ProducerPartitioner = "uniform_bytes"
	}

	if c.RequireAcks == "" {
		c.RequireAcks = "leader"
	}

	if c.Compression == nil {
		c.Compression = compression.CompressionTypes{compression.SnappyCompression, compression.NoneCompression}
	}

	if c.RecordRetries == 0 {
		c.RecordRetries = 9223372036854775807
	}
}
