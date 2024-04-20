package consumer_group_balancer

import (
	"github.com/twmb/franz-go/pkg/kgo"
)

type (
	BalancerType  int64
	BalancerTypes []BalancerType
)

const (
	CooperativeStickyBalancer BalancerType = iota
	RoundRobinBalancer
	RangeBalancer
	StickyBalancer
)

func Parse(balancers BalancerTypes) kgo.Opt {
	var groupBalancers []kgo.GroupBalancer

	for _, balancer := range balancers {
		switch balancer {
		case CooperativeStickyBalancer:
			groupBalancers = append(groupBalancers, kgo.CooperativeStickyBalancer())
		case StickyBalancer:
			groupBalancers = append(groupBalancers, kgo.StickyBalancer())
		case RangeBalancer:
			groupBalancers = append(groupBalancers, kgo.RangeBalancer())
		case RoundRobinBalancer:
			groupBalancers = append(groupBalancers, kgo.RoundRobinBalancer())
		default:
			continue
		}
	}

	if len(groupBalancers) == 0 {
		groupBalancers = []kgo.GroupBalancer{kgo.CooperativeStickyBalancer()}
	}

	return kgo.Balancers(groupBalancers...)
}
