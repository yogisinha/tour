package nodes

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"log"
	"math/rand"
	"net"
	"net/rpc"
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

type Master struct {
	sync.Mutex
	count         int
	l             net.Listener
	workers       []string // addresses (host:port) for workers so that Master can connect to them.
	done          chan bool
	shutdown      chan bool
	ServerAddress string
}

type RelayArgs struct {
	ID   string
	X, Y int
	Pic  []byte
	Hop  []string
}

// StartMaster will start the RPC server for the Master.
// It then divides the image blocks in parts and assign each block to some random id (between 1 - n (n is no of workers) )
// It will randomly hands off a block to some worker
//  method Master.Register  rpc: can't find service Master.Register: (Can we change go library code to show some helpful message
// showing the url of the specification or listing out the points for a method to be exported as a RPC end point)
func StartMaster(imagePath string, m, n int) *Master {
	mr := new(Master)
	mr.done = make(chan bool)
	mr.shutdown = make(chan bool)

	serverAddress := mr.connect() // starts listening on random port
	mr.ServerAddress = serverAddress

	go mr.serveConnections()

	go mr.sendFileFragments(imagePath, m, n)

	return mr
}

func (m *Master) connect() string {

	rpc.Register(m)

	l, e := net.Listen("tcp", "127.0.0.1:")

	if e != nil {
		log.Fatal("listen error in Master node :", e)
	}
	m.l = l
	fmt.Println("Master: chosen port ", m.l.Addr().String())

	return m.l.Addr().String()
}

func (m *Master) serveConnections() {

loop:
	for {
		select {
		case <-m.done:
			break loop
		default:
		}

		conn, err := m.l.Accept()
		if err != nil {
			log.Printf("Master: Error in Accept: %v\n", err)
			break
		}

		go rpc.ServeConn(conn)
	}

	fmt.Println("Master Stopped.")

}

func (m *Master) Done() {
	m.done <- true
}

func (m *Master) Wait() {
	<-m.done
}

func (m *Master) Start(args ClientArgs, reply *StartResult) error {
	fmt.Println("Master.Start called on server...")
	m.Lock()
	defer m.Unlock()

	m.workers = append(m.workers, args.ClientAddress)
	reply.ID = fmt.Sprintf("%d", m.count)
	m.count++

	fmt.Println("Workers registered with Master ", m.workers)

	return nil
}

// Sends the image fragments to client
func (m *Master) sendFileFragments(imagePath string, x, y int) {
	time.Sleep(2 * time.Second) // We will wait to Workers to start so that sends from Master won't get lost

	// open the file
	// make the RGBA
	// get the pix data
	// make the Fragment and send to ajax.Chan

	f, err := os.Open(imagePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in opening file %s: %v\n", imagePath, err)
		return
	}
	defer f.Close()

	ext := path.Ext(imagePath)
	var img image.Image
	switch strings.ToLower(ext) {
	case ".jpeg", ".jpg":
		img, err = jpeg.Decode(f)
	case ".png":
		img, err = png.Decode(f)
	case ".gif":
		img, err = gif.Decode(f)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in converting to image: %v\n", err)
		return
	}

	b := img.Bounds()
	newImg := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(newImg, newImg.Bounds(), img, b.Min, draw.Src)
	X, Y := b.Dx(), b.Dy()

	blocks := GetBlocks(X, Y, x, y)
	rand.Shuffle(len(blocks), func(i, j int) {
		blocks[i], blocks[j] = blocks[j], blocks[i]
	})

	for _, block := range blocks {
		m1 := newImg.SubImage(block)
		var buf bytes.Buffer
		jpeg.Encode(&buf, m1, nil)

		m.Lock()
		numWorkers := len(m.workers)
		clientID := fmt.Sprintf("%d", rand.Intn(numWorkers))
		randomAddress := m.workers[rand.Intn(numWorkers)]
		m.Unlock()

		args := RelayArgs{
			ID:  clientID,
			X:   block.Min.X,
			Y:   block.Min.Y,
			Pic: buf.Bytes()}

		hops := append(args.Hop, randomAddress)
		args.Hop = hops

		// Sending the image fragment to random worker. See the doc on Worker.Relay method, how it handles it
		Call(randomAddress, "Worker.Relay", args, new(struct{}))

	}

	fmt.Println("End of sendFilefragments...")
	m.Done()
}

func (m *Master) Random(_ struct{}, a *Address) error {
	m.Lock()
	defer m.Unlock()

	numWorkers := len(m.workers)
	randomAddress := m.workers[rand.Intn(numWorkers)]

	a.Addr = randomAddress
	return nil
}

func (m *Master) StopWorkers() {
	fmt.Println("In Master KillWorkers...")

	for _, clientPort := range m.workers {
		Call(clientPort, "Worker.Shutdown", new(struct{}), new(struct{}))
	}
}

func (m *Master) StopMaster() {
	fmt.Println("In Master StopMaster...")

	Call(m.ServerAddress, "Master.Shutdown", new(struct{}), new(struct{}))
}

func (m *Master) Shutdown(_, _ *struct{}) error {
	fmt.Println("In Master Shutdown method...")

	close(m.shutdown)
	if err := m.l.Close(); err != nil { // Close() will cause the Accept() to return with the error
		fmt.Printf("Shutdown: Error in closing listener in Master: %v", err)
	}
	return nil
}

// func (m *Master) StopMaster() {
// 	fmt.Println("In Master StopMaster...server port was ", m.l.Addr())
// 	close(m.shutdown)
// 	if err := m.l.Close(); err != nil {
// 		fmt.Printf("StopMaster: Error in closing listener in Master: %v", err)
// 	}
// }
