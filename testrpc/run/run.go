package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"tour/testrpc/server"
)

func main() {
	arith := new(server.Arith)
	rpc.Register(arith)
	rpc.Register(new(server.Node))

	rpc.HandleHTTP()
	l, e := net.Listen("tcp", "localhost:8000")
	if e != nil {
		log.Fatal("listen error:", e)
	}

	fmt.Println("Starting the server...", l.Addr())

	fmt.Println("Starting the server...")
	log.Fatal(http.Serve(l, nil))
}
