package micro

import (
	"context"
	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/resolver"
	"leanring-go/micro/registry"
	"time"
)

type grpcBuilder struct {
	r       registry.Registry
	timeout time.Duration
}

func NewRegistryBuilder(r registry.Registry, timeout time.Duration) (resolver.Builder, error) {
	return &grpcBuilder{
		r:       r,
		timeout: timeout,
	}, nil
}

func (b *grpcBuilder) Build(target resolver.Target, cc resolver.ClientConn,
	opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := &grpcResolver{
		cc:      cc,
		r:       b.r,
		target:  target,
		timeout: b.timeout,
	}
	r.resolve()
	go r.watch()
	return r, nil
}

func (b *grpcBuilder) Scheme() string {
	return "registry"
}

type grpcResolver struct {
	//   - "dns://some_authority/foo.bar"
	//     Target{Scheme: "dns", Authority: "some_authority", Endpoint: "foo.bar"}
	//  registry:///localhost:8081
	target  resolver.Target
	r       registry.Registry
	cc      resolver.ClientConn
	timeout time.Duration
	close   chan struct{}
}

func (g *grpcResolver) ResolveNow(options resolver.ResolveNowOptions) {
	g.resolve()
}

func (g *grpcResolver) watch() {
	events, err := g.r.Subscribe(g.target.Endpoint())
	if err != nil {
		g.cc.ReportError(err)
		return
	}
	for {
		select {
		case <-events:
			g.resolve()
		//case event := <- events:
		//	switch event.Type {
		//	case "DELETE":
		//	// 删除
		//	}
		case <-g.close:
			return
		}
	}
}

func (g *grpcResolver) resolve() {
	ctx, cancel := context.WithTimeout(context.Background(), g.timeout)
	defer cancel()
	instance, err := g.r.ListServices(ctx, g.target.Endpoint())
	if err != nil {
		g.cc.ReportError(err)
		return
	}
	address := make([]resolver.Address, 0, len(instance))
	for _, si := range instance {
		address = append(address, resolver.Address{
			Addr:       si.Address,
			Attributes: attributes.New("weight", si.Weight),
		})
	}
	err = g.cc.UpdateState(resolver.State{
		Addresses: address,
	})
	if err != nil {
		g.cc.ReportError(err)
		return
	}
}

func (g *grpcResolver) Close() {
	close(g.close)
}
