package rpc

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestInitClientProxy(t *testing.T) {
	server := NewServer()
	server.RegisterService(&UserServiceServer{})
	go func() {
		err := server.Start("tcp", ":8081")
		t.Log(err)
	}()
	time.Sleep(time.Second * 3)
	usClient := &UserService{}

	err := InitClientProxy(":8081", usClient)
	require.NoError(t, err)
	resp, err := usClient.GetById(context.Background(), &GetByIdReq{
		Id: 123,
	})
	assert.Equal(t, &GetByIdResp{
		Msg: "Hello world",
	}, resp)
}
