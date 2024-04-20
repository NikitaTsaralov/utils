package ack_policy

import (
	"github.com/twmb/franz-go/pkg/kgo"
)

type AckType int64

const (
	AllAck AckType = iota
	LeaderAck
	NoAck
)

func Parse(ackPolicy AckType) kgo.Opt {
	selectedAckPolicy := kgo.AllISRAcks() //

	switch ackPolicy {
	case NoAck:
		selectedAckPolicy = kgo.NoAck()
	case LeaderAck:
		selectedAckPolicy = kgo.LeaderAck()
	case AllAck:
		selectedAckPolicy = kgo.AllISRAcks()
	}

	return kgo.RequiredAcks(selectedAckPolicy)
}
