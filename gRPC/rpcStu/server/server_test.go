package server

import (
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"testing"
)

func TestServer_Arith(t *testing.T) {
	arith := new(Arith)
	//rpc.Register(arith)
	rpc.RegisterName("Arith", arith)
	rpc.HandleHTTP()
	l, err := net.Listen("tcp", ":1234")
	if err != nil {
		t.Fatal("listen error:", err)
	}
	go http.Serve(l, nil)

	client, err := rpc.DialHTTP("tcp", "localhost:1234")
	if err != nil {
		t.Fatal("dialing:", err)
	}

	// Synchronous call
	args := &Args{7, 8}
	var reply int
	err = client.Call("Arith.Multiply", args, &reply)
	if err != nil {
		t.Fatal("arith error:", err)
	}
	fmt.Printf("Arith: %d*%d=%d\n", args.A, args.B, reply)

	// Asynchronous call
	quotient := new(Quotient)
	divCall := client.Go("Arith.Divide", args, quotient, nil)
	replyCall := <-divCall.Done // will be equal to divCall
	// check errors, print, etc.
	if e := replyCall.Error; e != nil {
		t.Fatal("Asynchronous call error:", e)
	}
	r := replyCall.Reply
	fmt.Printf("Arith: %d*%d=%d\n", args.A, args.B, r)
}
