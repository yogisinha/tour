package distro

import (
	"fmt"
	
	"bytes"
	"net/http"
	"strconv"
	"sync"
)

func init() {
	fmt.Println("In init of fragment")
	
	http.HandleFunc("/fs/", poster)
	fmt.Println("Before registering pngHandler...")
	http.HandleFunc("/png/", pngHandler)
	fmt.Println("After registering pngHandler...")
}

var fsLock sync.Mutex
var fsMap = make(map[string][]byte)
var ch = make(chan []byte)	

func poster(conn http.ResponseWriter, req *http.Request) {
	name := req.URL.Path[4:]
	fsLock.Lock()
	data, ok := fsMap[name]
	fsLock.Unlock()
	if !ok {
		conn.WriteHeader(404)
		return
	}
	//conn.Header().Set("Content-Type", "image/png")
	conn.Header().Set("Content-Type", "image/jpeg")
	conn.Write(data)
}

func pngHandler(conn http.ResponseWriter, req *http.Request) {	
	//conn.Header().Set("Content-Type", "image/png")
	//fmt.Println("In pngHandler...")
	//time.Sleep(1*time.Second)
	conn.Header().Set("Content-Type", "image/jpeg")
	conn.Write(<-ch)
}

// Post saves data and returns a URL that can be used to load it.
// The data is assumed to be a PNG image.
func Post(data []byte) string {
	fsLock.Lock()
	n := strconv.Itoa(len(fsMap)) + ".jpeg"
	fsMap[n] = bytes.NewBuffer(data).Bytes()
	fsLock.Unlock()
	return "/fs/" + n
}

var c int
func MyPost(data []byte) string {	
	fsLock.Lock()
	c++
	fsLock.Unlock()

	go func() {
		ch <- data
	}()
	return "/png/" + strconv.Itoa(c) + ".jpeg"
}

// A Fragment is the message to send on ajax.Chan to display
// an image fragment.
type Fragment struct {
	X, Y int
	URL  string
}
