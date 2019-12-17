package nodes

import (
	"log"
	"net/rpc"
)

type StartArgs struct {
	Distributed bool
}

type StartResult struct {
	ID string
}

type Address struct {
	Addr string
}

var Ch chan struct{}

func Call(serverPort, method string, args, reply interface{}) {
	//fmt.Println("In call...", serverPort, method, args)

	client, err := rpc.Dial("tcp", ":"+serverPort)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	defer client.Close()

	// Synchronous call
	err = client.Call(method, args, reply)
	if err != nil {
		log.Printf("Error in RPC call : method %s  %v:", method, err)
	}

}
