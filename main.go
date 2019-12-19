package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"log"
	"math"
	"math/cmplx"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"
	"tour/ajax"
	"tour/chat"
	"tour/distributed/nodes"
	"tour/distro"
	"tour/fractal"
	"tour/pic"
	"tour/tree"
	"tour/wc"
)

var pixelFuns = map[string]func(x, y uint8) uint8{
	"X+Y":     func(x, y uint8) uint8 { return x + y },
	"X*Y":     func(x, y uint8) uint8 { return x * y },
	"(X+Y)/2": func(x, y uint8) uint8 { return (x + y) / 2 },
	"X^Y":     func(x, y uint8) uint8 { return x ^ y },
}

func main() {
	var sqrt = flag.Float64("sqrt", 1.0, "input for sqrt function exercise")
	var wcf = flag.Bool("wcf", false, "flag to run WordCount exercise")
	var slices = flag.String("slices", "", "flag to input function def to run Slices exercise. Enclose the value in double quotes. for e.g. tour -slices=\"X^Y\"")
	var hw2 = flag.Bool("helloworld2", false, "flag to run Hello World 2 exercise")
	var methods = flag.Bool("methods", false, "flag to run Methods exercise")
	var imgver = flag.String("imgvwer", "", "flag to input function def to run Image Viewer exercise. Enclose the value in double quotes. for e.g. tour -imgvwer=\"X^Y\"")
	var pngencode = flag.Bool("pngencode", false, "flag to run PNG Encoding exercise")
	var mandelbrot = flag.Bool("mandelbrot", false, "flag to run Mandelbrot Set Viewer exercise")
	var julia = flag.Bool("julia", false, "flag to run Julia Set Viewer exercise")
	var webLog = flag.Bool("weblog", false, "flag to run Web Logging exercise")
	var chatex = flag.Bool("chatex", false, "flag to run All the Chat related exercises")
	var rpc = flag.Bool("rpc", false, "flag to run RPC related exercises")
	var image = flag.String("image", "", "flag to input image path")
	var x = flag.Int("x", 50, "flag to input width of each block of image")
	var y = flag.Int("y", 50, "flag to input height of each block of image")
	var rand = flag.Bool("random", false, "flag to build the image in random order. All the blocks will appear in random order")
	var mode = flag.String("mode", "s", "possible mode flags to run the file download exercise. for sequential: s, for distributed: d")
	var workers = flag.Int("workers", 5, "No of workers")

	fmt.Println(os.Args)

	flag.Parse()

	switch {
	case *sqrt > 1.0:
		fmt.Println("Running Sqrt")
		fmt.Println(Sqrt(*sqrt))
	case *wcf:
		fmt.Println("Running WordCount", *slices)
		wc.Serve(WordCount)
	case *slices != "":
		fmt.Println("Running Slices", *slices)
		if pixelFunc, ok := pixelFuns[*slices]; ok {
			pic.Serve(picFunUint(pixelFunc))
		} else {
			log.Fatal("Unknown function provided.")
		}
	case *imgver != "":
		fmt.Println("Running Image Viewer", *imgver)
		if pixelFunc, ok := pixelFuns[*imgver]; ok {
			pic.ServeImage(picFunImage(pixelFunc))
		} else {
			log.Fatal("Unknown function provided.")
		}
	case *hw2:
		http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
			name := req.FormValue("name")
			fmt.Fprintf(w, "Hello, %s", name)
		})
	case *methods:
		http.Handle("/string", String("I’m a frayed knot."))
		http.Handle("/struct", &Struct{"Hello", ":", "USENIX!"})
	case *pngencode:
		http.Handle("/", Image1{})
	case *mandelbrot:
		http.Handle("/", fractal.MainPage) // Different color schemes...
		http.Handle("/mandelbrot",
			MyImage{
				requestParserFunc: mandelbrotReqParser,
			})
	case *julia:
		http.Handle("/", fractal.MainPage) // Different color schemes...  Go to http://localhost:4000/juliaviewer
		http.Handle("/julia",
			MyImage{
				requestParserFunc: juliaReqParser,
			})
	case *webLog:
		http.Handle("/", ajax.LogPage)
		go webLogging()
	case *chatex:
		http.Handle("/", chat.ChatPage)
		http.HandleFunc("/join", joinHandler)
		http.HandleFunc("/say", sayHandler)
		http.HandleFunc("/exit", exitHandler)
		http.Handle("/julia",
			MyImage{
				requestParserFunc: juliaReqParser,
			}) // for the julia bot.  Will respond with a julia image for messages like julia: 0.285+0.01i

		startChat() // goroutines started by this method does not stop once the program is stopped. Might be an exercise
	case *rpc:
		imagePath := *image
		f, err := os.Open(imagePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error in opening file %s: %v\n", imagePath, err)
			return
		}
		f.Close()

		http.Handle("/", distro.MainPage)
		if *mode == "s" {
			http.HandleFunc("/start",
				func(c http.ResponseWriter, req *http.Request) {
					sendFilefragments(imagePath, *x, *y, *rand)
				})
		} else if *mode == "d" {
			http.HandleFunc("/start",
				func(c http.ResponseWriter, req *http.Request) {
					startMasterAndWorkerNodes(imagePath, *workers, *x, *y)
				})
		}

	}

	if *hw2 || *methods || *pngencode || *mandelbrot || *julia || *webLog || *chatex || *rpc {
		fmt.Println("Server running at 127.0.0.1:4000...")
		log.Fatal(http.ListenAndServe(":4000", nil))
	}

}

func Sqrt(x float64) float64 {
	var threshold = 0.01
	z := float64(1)
	var old = z - (z*z-x)/2*z
	var new float64
	for {
		new = old - (old*old-x)/(2*old)
		if math.Abs(new-old) < threshold {
			return new
		}

		old = new
	}

}

func WordCount(s string) map[string]int {
	wordScanner := bufio.NewScanner(strings.NewReader(s))
	wordScanner.Split(bufio.ScanWords)

	var mp = map[string]int{}
	for wordScanner.Scan() {
		word := wordScanner.Text()
		mp[word] = mp[word] + 1
	}

	return mp
}

var pixelFunc func(x, y uint8) uint8

func picFunUint(pixelFunc func(x, y uint8) uint8) func(x, y int) [][]uint8 {

	return func(dx, dy int) [][]uint8 {
		var grid = make([][]uint8, dx)

		for i := 0; i < dx; i++ {
			var y = make([]uint8, dy)
			for j := 0; j < dy; j++ {
				y[j] = pixelFunc(uint8(i), uint8(j))
			}
			grid[i] = y
		}

		return grid
	}

}

type Image1 struct {
	x, y        int
	pixelColors [][]color.Color
}

func (i Image1) ColorModel() color.Model {
	return color.RGBAModel
}

func (i Image1) Bounds() image.Rectangle {
	return image.Rect(0, 0, i.x, i.y)
}

func (i Image1) At(x, y int) color.Color {
	return i.pixelColors[x][y]
}

func picFunImage(pixelFunc func(x, y uint8) uint8) func(x, y int) image.Image {

	return func(dx, dy int) image.Image {

		img := Image1{x: dx, y: dy}
		img.pixelColors = make([][]color.Color, dx)

		for i := 0; i < dx; i++ {
			var y = make([]color.Color, dy)
			for j := 0; j < dy; j++ {
				v := pixelFunc(uint8(i), uint8(j))
				y[j] = color.RGBA{v, v, v, 255}
			}
			img.pixelColors[i] = y
		}

		return img
	}

}

func PicImage(dx, dy int) image.Image {

	img := Image1{x: dx, y: dy}
	img.pixelColors = make([][]color.Color, dx)

	for i := 0; i < dx; i++ {
		var y = make([]color.Color, dy)
		for j := 0; j < dy; j++ {
			v := (uint8(i) * uint8(j))
			y[j] = color.RGBA{v, v, v, 255}
		}
		img.pixelColors[i] = y
	}

	return img
}

func (i Image1) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Println("In ServeHTTP method of Image1...")
	x, _ := strconv.Atoi(req.FormValue("x"))
	y, _ := strconv.Atoi(req.FormValue("y"))

	img := PicImage(x, y)

	var buf bytes.Buffer
	png.Encode(&buf, img)

	w.Header().Set("Content-Type", "image/png")
	w.Write(buf.Bytes())
}

type String string

type Struct struct {
	Greeting string
	Punct    string
	Who      string
}

func (s String) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte(s))
}

func (s *Struct) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	b, _ := json.Marshal(*s)
	w.Write(b)
}

// Structs and methods for Mandelbrot and Julia Set viewer exercise -- Start

type Colorizer interface {
	At(x, y int) color.Color
}

type mandelbrotColorizer struct {
	props ColorProps
}

type juliaColorizer struct {
	props ColorProps
}

type MyImage struct {
	requestParserFunc func(string) (ColorProps, Colorizer, error)
	colorProps        ColorProps
	Colorizer
}

type ColorProps struct {
	width, height, iterations int
	origin, crange, c         complex128
}

func (m MyImage) ColorModel() color.Model {
	return color.RGBAModel
}

func (m MyImage) Bounds() image.Rectangle {
	return image.Rect(0, 0, m.colorProps.width, m.colorProps.height)
}

func mandelbrotReqParser(p string) (ColorProps, Colorizer, error) {
	var (
		x, y, iterations int
		origin, crange   complex128
	)

	r := strings.NewReader(p)
	_, err := fmt.Fscanf(r, "%d %d %g %g %d", &x, &y, &origin, &crange, &iterations)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fscanf: %v\n", err)
		return ColorProps{}, mandelbrotColorizer{}, err
	}
	fmt.Println(x, y, origin, crange, iterations)

	return ColorProps{x, y, iterations, origin, crange, 0},
		mandelbrotColorizer{
			ColorProps{x, y, iterations, origin, crange, 0}}, nil
}

func (m mandelbrotColorizer) At(x, y int) color.Color {
	fx := float64(x) / float64(m.props.width)
	fy := float64(y) / float64(m.props.height)
	c := m.props.origin + complex(real(m.props.crange)*fx, imag(m.props.crange)*fy)

	return mandelbrotColor(c, m.props.iterations)
}

func mandelbrotColor(c complex128, iterations int) color.Color {
	const contrast = 15
	var z complex128
	for i := 0; i < iterations; i++ {
		z = z*z + c
		if cmplx.Abs(z) > 2 {
			return color.Gray{255 - uint8(contrast*i)}
			//return fractal.Cycle(i, iterations)
			//return fractal.Ramp(i, iterations)
		}
	}

	return color.Black
}

func juliaReqParser(p string) (ColorProps, Colorizer, error) {
	var (
		x, y, iterations  int
		origin, crange, c complex128
	)

	r := strings.NewReader(p)
	_, err := fmt.Fscanf(r, "%d %d %g %g %g %d", &x, &y, &origin, &crange, &c, &iterations)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fscanf: %v, parsing the input : %s\n", err, p)
		return ColorProps{}, juliaColorizer{}, err
	}
	fmt.Println(x, y, origin, crange, c, iterations)

	return ColorProps{x, y, iterations, origin, crange, c},
		juliaColorizer{
			ColorProps{x, y, iterations, origin, crange, c}}, nil
}

func (m juliaColorizer) At(x, y int) color.Color {
	fx := float64(x) / float64(m.props.width)
	fy := float64(y) / float64(m.props.height)
	z := m.props.origin + complex(real(m.props.crange)*fx, imag(m.props.crange)*fy)

	return juliaColor(z, m.props.c, m.props.iterations)
}

func juliaColor(z, c complex128, iterations int) color.Color {
	const contrast = 15
	for i := 0; i < iterations; i++ {
		z = z*z + c
		if cmplx.Abs(z) > 2 {
			//return color.Gray{255 - uint8(contrast*i)}
			//return fractal.Cycle(i, iterations)
			return fractal.Ramp(i, iterations)
		}
	}
	return color.Black
}

func (m MyImage) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	colorProps, colorizer, _ := m.requestParserFunc(req.FormValue("p"))
	m.colorProps = colorProps
	m.Colorizer = colorizer

	var buf bytes.Buffer
	png.Encode(&buf, m)

	w.Header().Set("Content-Type", "image/png")
	w.Write(buf.Bytes())
}

// Structs and methods for Mandelbrot and Julia Set viewer exercise -- End

// Same checks whether 2 Binary search Trees are equivalent or not.
// Exercise using goroutines and channels
func Same(t1, t2 *tree.Tree) bool {
	var ch1 = make(chan int)
	var ch2 = make(chan int)

	startWalk(t1, ch1)
	startWalk(t2, ch2)

	for {
		v1, ok1 := <-ch1
		v2, ok2 := <-ch2

		if ok1 && ok2 {
			if v1 == v2 {
				continue
			} else {
				drainChannel(ch1)
				drainChannel(ch2)
				return false
			}
		} else if ok1 && !ok2 {
			drainChannel(ch1)
			return false
		} else if !ok1 && ok2 {
			drainChannel(ch2)
			return false
		} else {
			return true
		}

	}

}

// draining the channels to avoid the goroutine leak.
func drainChannel(ch <-chan int) {
	for range ch {
	}
}

func startWalk(t *tree.Tree, ch chan<- int) {
	var wg sync.WaitGroup
	wg.Add(1)
	go Walk(t, ch, &wg)

	go func() {
		wg.Wait()
		close(ch)
	}()
}

func Walk(t *tree.Tree, ch chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	if t == nil {
		return
	}

	wg.Add(1)
	Walk(t.Left, ch, wg)
	ch <- t.Value

	wg.Add(1)
	Walk(t.Right, ch, wg)
}

type Log string

func webLogging() {
	t := time.NewTicker(time.Second)
	abort := make(chan struct{})
	go func() {
		os.Stdin.Read(make([]byte, 1))
		t.Stop()
		abort <- struct{}{}
	}()

	var i int
loop:
	for {
		select {
		case <-t.C:
			ajax.Chan <- Log(fmt.Sprintf("%d", i))
			i++
		case <-abort:
			fmt.Println("In default...")
			break loop
		}
	}

	fmt.Println("End")
}

// Structs, channel definitions and methods for Chat related exercises. -- Start

var chatRoom = make(chan interface{})
var juliaChan = make(chan Say)

func startChat() {
	fmt.Println("Starting the chat manager...")
	go chatManager()
	go receiveAllMsgs()
}

func receiveAllMsgs() {
	const imgurl = "512 512 -2-1i 3+2i %s 100"

	for {
		s := <-juliaChan
		msg := s.Msg
		chatRoom <- s // send the normal message
		// Send the img tag if message is for Julia Bot (ex message for julia bot, "julia: 0.285+0.01i")
		if strings.HasPrefix(msg, "julia:") {
			parts := strings.Split(msg, " ")
			c := strings.Trim(parts[1], "")
			imgurl := "<img src=\"/julia?p=" + url.QueryEscape(fmt.Sprintf(imgurl, c)) + "\">"
			chatRoom <- Say{"julia", s.Who + ": " + imgurl}
		}
	}

}

type Status string

func chatManager() {
	members := map[string]bool{}

	for {
		switch m := (<-chatRoom).(type) {
		case Join:
			ajax.Chan <- Chat{m.Who, m.Who + " joined the room."}
			ajax.Chan <- Status(joinStatus(m.Who, members))
		case Say:
			ajax.Chan <- Chat{m.Who, m.Msg}
		case Exit:
			ajax.Chan <- Chat{m.Who, m.Who + " left the chat."}
			ajax.Chan <- Status(exitStatus(m.Who, members))
		}
	}
}

func sliceToString(l []string) string {
	var buf bytes.Buffer
	buf.WriteByte('[')
	for i, n := range l {
		if i > 0 {
			buf.WriteString(", ")
		}
		fmt.Fprintf(&buf, "%s", n)
	}
	buf.WriteByte(']')
	return buf.String()
}

func exitStatus(n string, members map[string]bool) string {
	l := []string{"-" + n}
	delete(members, n)
	for m := range members {
		l = append(l, m)
	}

	return sliceToString(l)
}

func joinStatus(n string, members map[string]bool) string {
	l := []string{"+" + n}
	for m := range members {
		l = append(l, m)
	}
	members[n] = true

	return sliceToString(l)
}

type Join struct {
	Who string
}

type Exit struct {
	Who string
}

type Say struct {
	Who string
	Msg string
}

type Chat struct {
	Who string
	Msg string
}

func joinHandler(c http.ResponseWriter, req *http.Request) {
	chatRoom <- Join{req.FormValue("id")}
}

func sayHandler(c http.ResponseWriter, req *http.Request) {
	juliaChan <- Say{req.FormValue("id"), req.FormValue("msg")}
}

func exitHandler(c http.ResponseWriter, req *http.Request) {
	chatRoom <- Exit{req.FormValue("id")}
}

// Structs, channel definitions and methods for Chat related exercises. -- End

func sendFilefragments(imagePath string, m, n int, random bool) {
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

	blocks := nodes.GetBlocks(X, Y, m, n)
	if random {
		rand.Shuffle(len(blocks), func(i, j int) {
			blocks[i], blocks[j] = blocks[j], blocks[i]
		})
	}

	for _, block := range blocks {
		m1 := newImg.SubImage(block)
		var buf bytes.Buffer
		png.Encode(&buf, m1)
		ajax.Chan <- distro.Fragment{X: block.Min.X, Y: block.Min.Y, URL: distro.Post(buf.Bytes())}
	}

}

func startMasterAndWorkerNodes(imagePath string, workers, m, n int) {

	mr := nodes.StartMaster(imagePath, m, n)

	nodes.StartWorkers(workers, mr.ServerAddress)

	mr.Wait()

	mr.StopWorkers()
	mr.StopMaster()

}
