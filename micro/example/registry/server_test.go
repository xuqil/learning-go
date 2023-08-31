package grpc_resolver

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	clientv3 "go.etcd.io/etcd/client/v3"
	"leanring-go/micro"
	"leanring-go/micro/proto/gen"
	"leanring-go/micro/registry/etcd"
	"testing"
)

func TestServer(t *testing.T) {
	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"localhost:2379"},
	})
	require.NoError(t, err)
	r, err := etcd.NewRegistry(etcdClient)
	us := &UserServiceServer{}
	server, err := micro.NewServer("user-service", micro.ServerWithRegistry(r))
	require.NoError(t, err)
	gen.RegisterUserServiceServer(server, us)

	// 意味着 us 完全准备好了
	err = server.Start(":8081")
	t.Log(err)
}

type UserServiceServer struct {
	gen.UnimplementedUserServiceServer
}

func (s UserServiceServer) GetById(ctx context.Context, req *gen.GetByIdReq) (*gen.GetByIdResp, error) {
	fmt.Println("req")
	return &gen.GetByIdResp{
		User: &gen.User{
			Name: "hello, world",
		},
	}, nil
}
