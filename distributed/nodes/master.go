package nodes

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"log"
	"math"
	"math/rand"
	"net"
	"net/rpc"
	"os"
	"sync"
	"time"
)

type Master struct {
	sync.Mutex
	count   int
	l       net.Listener
	workers []string // port nos for workers
	portMap map[int]string
}

// This function will start the RPC server
// It then divides the image blocks in parts and assign each block to some random id (between 1 - n (n is no of workers) )
// It will randomly hands off a block to some worker
//  method Master.Register  rpc: can't find service Master.Register: (Can we change go library code to show some helpful message
// showing the url of the specification or listing out the points for a method to be exported as a RPC end point)
func StartServer() *Master {
	m := new(Master)
	m.portMap = make(map[int]string)
	go m.startRPCServer(ServerPort) // test ur own rpc code with this code..

	filePath := string(`C:\Users\600131444\Pictures\vehicleImage.jpg`)
	go m.sendFileFragments(filePath)

	return m
}

func (m *Master) Start(args ClientArgs, reply *StartResult) error {
	fmt.Println("Master.Start called on server...")
	m.Lock()
	defer m.Unlock()

	m.workers = append(m.workers, args.ClientPort)
	//m.portMap[m.count] = args.ClientPort
	reply.ID = fmt.Sprintf("%d", m.count)
	m.count++

	fmt.Println(m.workers)

	return nil
}

func getBlocks(X, Y int, m, n int) []image.Rectangle {
	var xsteps, ysteps int
	if X%m == 0 {
		xsteps = X / m
	} else {
		xsteps = X/m + 1
	}

	if Y%n == 0 {
		ysteps = Y / n
	} else {
		ysteps = Y/n + 1
	}

	if m >= X {
		xsteps = 1
	}
	if n >= Y {
		ysteps = 1
	}

	var x0, y0 int
	var blocks []image.Rectangle
	for y := 0; y < ysteps; y++ {
		for x := 0; x < xsteps; x++ {
			frstRect := image.Rect(x0, y0, x0+m, y0+n)
			secondRect := image.Rect(x0, y0,
				int(math.Min(float64(x0+m), float64(X))),
				int(math.Min(float64(y0+n), float64(Y))))

			blocks = append(blocks, frstRect.Intersect(secondRect))
			x0 += m
		}

		x0 = 0
		y0 += n
	}

	rand.Shuffle(len(blocks), func(i, j int) {
		blocks[i], blocks[j] = blocks[j], blocks[i]
	})

	return blocks

}

// Sends the image fragments to client
func (mr *Master) sendFileFragments(filePath string) {
	time.Sleep(1 * time.Second)
	// open the file
	// make the RGBA
	// get the pix data
	// make the Fragment and send to ajax.Chan
	f, err := os.Open(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in opening file %s: %v\n", filePath, err)
		return
	}
	defer f.Close()

	img, err := jpeg.Decode(f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in converting to image: %v\n", err)
		return
	}

	b := img.Bounds()
	new_m := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(new_m, new_m.Bounds(), img, b.Min, draw.Src)
	X, Y := 247, 165
	m, n := 10, 10

	for _, block := range getBlocks(X, Y, m, n) {
		m1 := new_m.SubImage(block)
		var buf bytes.Buffer
		jpeg.Encode(&buf, m1, nil)

		mr.Lock()
		numWorkers := len(mr.workers)
		clientID := fmt.Sprintf("%d", rand.Intn(numWorkers))
		randomPort := mr.workers[rand.Intn(numWorkers)]
		mr.Unlock()

		//fmt.Printf("Sending to client Id: %s and random port: %s\n", clientID, randomPort)

		args := RelayArgs{
			ID:  clientID,
			X:   block.Min.X,
			Y:   block.Min.Y,
			Pic: buf.Bytes()}

		hops := append(args.Hop, randomPort)
		args.Hop = hops

		Call(randomPort, "Worker.Relay", args, new(struct{}))

		//ajax.Chan <- distro.Fragment{X: block.Min.X, Y: block.Min.Y, URL: distro.Post(buf.Bytes())}
	}

	fmt.Println("End of sendFilefragments...")
}

func (m *Master) Register(args *ClientArgs, b *bool) error {
	fmt.Println("Master.Register called on server...")
	m.Lock()
	defer m.Unlock()

	m.workers = append(m.workers, args.ClientPort)
	fmt.Println("Master.Register ", m.workers)
	*b = true
	return nil
}

func (mr *Master) Random(_ struct{}, a *Address) error {	
	mr.Lock()
	defer mr.Unlock()

	numWorkers := len(mr.workers)
	randomPort := mr.workers[rand.Intn(numWorkers)]

	a.Addr = randomPort
	return nil
}

func (m *Master) KillWorkers() {
	fmt.Println("In Master KillWorkers...")
	m.Lock()
	defer m.Unlock()

	for _, clientPort := range m.workers {
		Call(clientPort, "Worker.Shutdown", new(struct{}), new(struct{}))
	}
}

func (m *Master) StopMaster() {
	fmt.Println("In Master StopMaster...")
	m.Lock()
	defer m.Unlock()

	Call(ServerPort, "Master.Shutdown", new(struct{}), new(struct{}))
}

func (m *Master) Shutdown(_, _ *struct{}) error {
	fmt.Println("In Master Shutdown method...")
	m.Lock()
	defer m.Unlock()
	m.l.Close()
	return nil
}

func (m *Master) startRPCServer(port string) {
	//master := new(Master)
	//rpcs := rpc.NewServer()
	rpc.Register(m)

	//rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":"+port)
	if e != nil {
		log.Fatal("listen error in Master node :", e)
	}
	m.l = l

	// go func() {
	// 	for {
	// 		conn, err := l.Accept()
	// 		if err != nil {
	// 			log.Printf("Server: Error in Accept: %v", err)
	// 			continue
	// 		}

	// 		go rpc.ServeConn(conn)
	// 	}

	// }()

	//loop:
	for {
		// select {
		// case <-Ch:
		// 	break loop
		// default:
		// }

		fmt.Println("Before Accept in Master...")
		conn, err := l.Accept()
		if err != nil {
			log.Printf("Server: Error in Accept: %v", err)
			break
		}

		go rpc.ServeConn(conn)
	}

	fmt.Println("Exiting the serving loop for master...")

	// fmt.Println("Starting the RPC server on Master...", l.Addr())
	// log.Fatal(http.Serve(l, nil))

}
