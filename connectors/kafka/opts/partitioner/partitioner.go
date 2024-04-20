package partitioner

import (
	"github.com/twmb/franz-go/pkg/kgo"
)

type PartitionerType int64

const (
	UniformBytes PartitionerType = iota
	LeastBackup
	Manual
	RoundRobin
	StickyKey
	Sticky
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
