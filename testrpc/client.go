package main

import (
	"fmt"
	"log"
	"net/rpc"
	"sync"
	"tour/testrpc/server"
)

func main() {
	client, err := rpc.DialHTTP("tcp", "localhost:8000")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	// Synchronous call
	args := &server.Args{7, 8}
	var reply int
	err = client.Call("Arith.Multiply", args, &reply)
	if err != nil {
		log.Fatal("arith error:", err)
	}
	fmt.Printf("Arith: %d*%d=%d", args.A, args.B, reply)

	arg := &server.Arg{}
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		go func() {
			wg.Add(1)
			defer wg.Done()
			var res server.Result
			err = client.Call("Node.Start", arg, &res)
			if err != nil {
				log.Fatal("Node.Start error:", err)
			}
			fmt.Printf("Result %d\n", res.ID)

		}()
	}

	wg.Wait()
	var res server.Result
	err = client.Call("Node.GetID", arg, &res)
	if err != nil {
		log.Fatal("Node.Start error:", err)
	}
	fmt.Printf("GetID %d\n", res.ID)

}
