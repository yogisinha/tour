package nodes

import (
	"tour/distro"
	"tour/ajax"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"sync"
)

const ServerPort = "8000"

type Worker struct {
	port string
	sync.Mutex
	c  int
	id string // ID returned by Master
	// pending image blocks to process
	l net.Listener
}

type ClientArgs struct {
	ClientPort string
}

type RelayArgs struct {
	ID   string
	X, Y int
	Pic  []byte
	Hop []string
}


func StartWorkers(ports []string) {

	for _, port := range ports {
		//go w.startRPCServer(port)
		go startRPCServer(port)
	}

	// time.Sleep(time.Hour)
	// fmt.Println("startWorkers exiting ...")

}

func (w *Worker) register(serverPort string) {
	var b bool

	Call(serverPort, "Master.Register", &ClientArgs{w.port}, &b)

	if b {
		fmt.Println("Client registered with server ", w.port)
	} else {
		fmt.Println("Client Not registered with server ", w.port)
	}

}

func (w *Worker) getIDFromMaster(serverPort string) {
	w.Lock()
	defer w.Unlock()
	var s StartResult

	Call(serverPort, "Master.Start", ClientArgs{w.port}, &s)
	w.id = s.ID
	fmt.Println("ID returned by Master " + w.id + " for port " + w.port)
}

func (w *Worker) Shutdown(_, _ *struct{}) error {
	w.Lock()
	defer w.Unlock()
	fmt.Println("In Worker Shutdown method..." + w.port)
	w.c = 1

	w.l.Close()
	return nil
}

func (w *Worker) Relay(args RelayArgs, _ *struct{}) error {
	
	// ajax.Chan <- distro.Fragment{
	// 	X: args.X, Y: args.Y,
	// 	URL: distro.Post(args.Pic)}

	// w.Lock()
	// defer w.Unlock()

	if args.ID != w.id {
		fmt.Printf("Client Id: %s and port: %s, Request for other Client: %s\n", w.id, w.port, args.ID)
		
		var a Address
		Call(ServerPort, "Master.Random", struct{}{}, &a)
		fmt.Println("Afte Master.Random call, remote address: ", a.Addr, w.port)

		hops := append(args.Hop, a.Addr)
		args.Hop = hops
		Call(a.Addr, "Worker.Relay", args, new(struct{}))

	} else {

		fmt.Printf("Client Id: %s and port: %s, Request for same Client: %s, total Hops:%v \n", w.id, w.port, args.ID, args.Hop)
		ajax.Chan <- distro.Fragment{
			X: args.X, Y: args.Y,
			URL: distro.MyPost(args.Pic)}
	}

	return nil
}

func startRPCServer(port string) {
	w := new(Worker)
	w.port = port
	rpcs := rpc.NewServer()
	rpcs.Register(w)

	//rpc.HandleHTTP()

	//w.register(ServerPort) // register with server. sending its own port so that server can call back RPCs on workers

	fmt.Println("net Listen on port " + w.port)
	l, e := net.Listen("tcp", ":"+w.port)
	if e != nil {
		log.Fatal("listen error in Worker node :", e)
	}
	w.l = l
	w.getIDFromMaster(ServerPort)

	for {
		// select {
		// case <-Ch:
		// 	break loop
		// default:
		// }

		fmt.Println("Before Accept in Worker..." + w.port)
		conn, err := l.Accept()
		if err != nil {
			log.Printf("Worker: Error in Accept: %v", err)
			break
		} else {
			// w.Lock()
			// fmt.Printf("After Accept %d\n", w.c)
			// w.c = 100
			// fmt.Printf("After Accept assigning %d\n", w.c)

			// w.Unlock()
		}

		//fmt.Println("Before serving the connection...")
		go rpcs.ServeConn(conn)
	}

	fmt.Println("Exiting the serving loop for worker..." + port)

	// fmt.Println("Starting the RPC server on Worker...", l.Addr())
	// log.Fatal(http.Serve(l, nil))

}
