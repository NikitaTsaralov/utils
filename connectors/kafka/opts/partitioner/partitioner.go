package partitioner

import (
	"github.com/twmb/franz-go/pkg/kgo"
)

type PartitionerType string

const (
	UniformBytes PartitionerType = "uniform_bytes"
	LeastBackup                  = "least_backup"
	Manual                       = "manual"
	RoundRobin                   = "round_robin"
	StickyKey                    = "sticky_key"
	Sticky                       = "sticky"
)

func Parse(partitioner PartitionerType) kgo.Opt {
	selectedPartitioner := kgo.UniformBytesPartitioner(32<<10, true, true, nil)

	switch partitioner {
	case UniformBytes:
		// Default hasher is murmur2.
		selectedPartitioner = kgo.UniformBytesPartitioner(32<<10, true, true, nil)
	case LeastBackup:
		selectedPartitioner = kgo.LeastBackupPartitioner()
	case Manual:
		selectedPartitioner = kgo.ManualPartitioner()
	case RoundRobin:
		selectedPartitioner = kgo.RoundRobinPartitioner()
	case StickyKey:
		// Default hasher murmur2.
		selectedPartitioner = kgo.StickyKeyPartitioner(nil)
	case Sticky:
		selectedPartitioner = kgo.StickyPartitioner()
	}

	return kgo.RecordPartitioner(selectedPartitioner)
}
