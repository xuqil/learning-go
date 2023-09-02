package micro

import (
	"context"
	"google.golang.org/grpc"
	"leanring-go/micro/registry"
	"net"
	"time"
)

type ServerOption func(server *Server)

type Server struct {
	name            string
	registry        registry.Registry
	registerTimeout time.Duration
	*grpc.Server
	listener net.Listener
	weight   uint32
	group    string
}

func NewServer(name string, opts ...ServerOption) (*Server, error) {
	res := &Server{
		name:            name,
		Server:          grpc.NewServer(),
		registerTimeout: time.Second * 10,
	}
	for _, opt := range opts {
		opt(res)
	}
	return res, nil
}

func ServerWithWeight(weight uint32) ServerOption {
	return func(server *Server) {
		server.weight = weight
	}
}

func ServerWithGroup(group string) ServerOption {
	return func(server *Server) {
		server.group = group
	}
}

// Start 当用户调用这个方法的时候，说明服务已经准备好了
func (s *Server) Start(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	s.listener = listener
	// 如果有注册中心，开始注册
	if s.registry != nil {
		ctx, cancel := context.WithTimeout(context.Background(), s.registerTimeout)
		defer cancel()
		err = s.registry.Register(ctx, registry.ServiceInstance{
			Name:    s.name,
			Address: listener.Addr().String(),
			Group:   s.group,
		})
		if err != nil {
			return err
		}
		//defer func() {
		//	// 忽略或者 log 一下
		//	_ = s.registry.Close()
		//}()
	}
	err = s.Serve(listener)
	return nil
}

func (s *Server) Close() error {
	if s.registry != nil {
		err := s.registry.Close()
		if err != nil {
			return err
		}
	}
	s.GracefulStop()
	return nil
}

func ServerWithRegistry(r registry.Registry) ServerOption {
	return func(server *Server) {
		server.registry = r
	}
}
