package round_robin

import (
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/resolver"
	"leanring-go/micro/route"
	"sync/atomic"
)

type Balancer struct {
	index       int32
	connections []subConn
	length      int32
	filter      route.Filter
}

func (b *Balancer) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	candidates := make([]subConn, 0, len(b.connections))
	for _, c := range b.connections {
		if b.filter != nil && !b.filter(info, c.addr) {
			continue
		}
		candidates = append(candidates, c)
	}
	if len(candidates) == 0 {
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}
	idx := atomic.AddInt32(&b.index, 1)
	c := candidates[int(idx)%len(candidates)]
	return balancer.PickResult{
		SubConn: c.c,
		Done: func(info balancer.DoneInfo) {

		},
	}, nil
}

type Builder struct {
	Filter route.Filter
}

func (b *Builder) Build(info base.PickerBuildInfo) balancer.Picker {
	connection := make([]subConn, 0, len(info.ReadySCs))
	for c, ci := range info.ReadySCs {
		connection = append(connection, subConn{
			c:    c,
			addr: ci.Address,
		})
	}

	return &Balancer{
		connections: connection,
		index:       -1,
		length:      int32(len(connection)),
		filter:      b.Filter,
	}
}

type subConn struct {
	c    balancer.SubConn
	addr resolver.Address
}
