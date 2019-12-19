package nodes

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"tour/ajax"
	"tour/distro"
)

type Worker struct {
	id            string // ID returned by Master
	l             net.Listener
	serverAddress string
}

type ClientArgs struct {
	ClientAddress string
}

func StartWorkers(workers int, serverAddress string) {

	for i := 0; i < workers; i++ {
		go startRPCServer(serverAddress)
	}

}

func (w *Worker) getIDFromMaster(serverAddress string) {
	var s StartResult

	Call(serverAddress, "Master.Start", ClientArgs{w.l.Addr().String()}, &s)
	w.id = s.ID
	fmt.Println("ID returned by Master for worker ", w.l.Addr(), " : "+w.id)
}

func (w *Worker) Shutdown(_, _ *struct{}) error {
	fmt.Println("In Worker Shutdown method...", w.l.Addr())

	w.l.Close()
	return nil
}

// Relay method checks if the image block sent from Master was meant for itself
// then it sends the image block to ajax.Chan (it goes to browser)
// otherwise it asks for a random worker address from Master and forwards the image block
// to that worker.
func (w *Worker) Relay(args RelayArgs, _ *struct{}) error {

	if args.ID != w.id {

		var a Address
		Call(w.serverAddress, "Master.Random", struct{}{}, &a)

		hops := append(args.Hop, a.Addr)
		args.Hop = hops
		Call(a.Addr, "Worker.Relay", args, new(struct{}))

	} else {

		ajax.Chan <- distro.Fragment{
			X: args.X, Y: args.Y,
			URL: distro.MyPost(args.Pic)}
	}

	return nil
}

func startRPCServer(serverAddress string) {
	w := new(Worker)
	w.serverAddress = serverAddress
	rpcs := rpc.NewServer()
	rpcs.Register(w)

	l, e := net.Listen("tcp", "127.0.0.1:")
	if e != nil {
		log.Fatal("listen error in Worker node :", e)
	}
	w.l = l

	fmt.Println("Worker: chosen port ", w.l.Addr())
	w.getIDFromMaster(serverAddress)

	for {

		conn, err := l.Accept()
		if err != nil {
			log.Printf("Worker: Error in Accept: %v\n", err)
			break
		}

		go rpcs.ServeConn(conn)
	}

	fmt.Printf("Worker with address %s Stopped.\n", w.l.Addr())

}
