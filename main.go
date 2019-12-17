package main

import (
	"tour/fractal"
	"tour/wc"
	"flag"
	"strconv"
	"tour/pic"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"log"
	"math"
	"math/cmplx"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
	"tour/ajax"
	"tour/distributed/nodes"
	"tour/distro"
	"tour/tree"
)

var pixelFuns = map[string]func(x, y uint8) uint8{
	"X+Y": func(x, y uint8) uint8 { return x + y },
	"X*Y": func(x, y uint8) uint8 { return x * y },
	"(X+Y)/2": func(x, y uint8) uint8 { return (x + y)/2 },
	"X^Y": func(x, y uint8) uint8 { return x ^ y }, 
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
		http.HandleFunc("/",  func (w http.ResponseWriter, req *http.Request) {
			name := req.FormValue("name")
			fmt.Fprintf(w, "Hello %s", name)
		} )

		http.Handle("/string", String("I’m a frayed knot."))
		http.Handle("/struct", &Struct{"Hello", ":", "USENIX!"})

		log.Fatal(http.ListenAndServe(":4000", nil))
	case *methods:
		http.Handle("/string", String("I’m a frayed knot."))
		http.Handle("/struct", &Struct{"Hello", ":", "USENIX!"})

		log.Fatal(http.ListenAndServe(":4000", nil))
	case *pngencode:
		http.Handle("/", Image1{})
		log.Fatal(http.ListenAndServe(":4000", nil))
	case *mandelbrot:
		http.Handle("/", fractal.MainPage)
		http.Handle("/mandelbrot",
			MyImage{
				requestParserFunc: mandelbrotReqParser,
			})
		log.Fatal(http.ListenAndServe(":4000", nil))


	






	default:

	}



	// 

	//fmt.Println(Sqrt(68))

	//fmt.Println(Same(tree.New(1), tree.New(2)))

	//go webLogging()

	// filePath := string(`C:\development\workspace\Fleeetnet_main\WebContent\images\add.png`)
	// sendFilefragments(filePath)

	// r1 := image.Rect(0, 0, 10, 5)
	// r2 := image.Rect(0, 0, 5, 10)

	//fmt.Println(r1.Intersect(r2))

	//testSplits()

	//startChat()
	//startServer()

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
	
		for i:=0; i<dx; i++ {
			var y = make([]uint8, dy)
			for j:=0; j<dy; j++ {
				y[j] = pixelFunc(uint8(i), uint8(j))
			}
			grid[i] = y
		}
	
		return grid
	}

}

type Image1 struct {
	x, y int
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
	
		img := Image1{x:dx, y:dy}
		img.pixelColors = make([][]color.Color, dx)
	
		for i:=0; i<dx; i++ {
			var y = make([]color.Color, dy)
			for j:=0; j<dy; j++ {
				v := pixelFunc(uint8(i), uint8(j))
				y[j] = color.RGBA{v, v, v, 255}
			}
			img.pixelColors[i] = y
		}
	
		return img
	}

}


func PicImage(dx, dy int) image.Image {	

	img := Image1{x:dx, y:dy}
	img.pixelColors = make([][]color.Color, dx)

	for i:=0; i<dx; i++ {
		var y = make([]color.Color, dy)
		for j:=0; j<dy; j++ {
			v := (uint8(i)*uint8(j))
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


func testSplits() {
	X, Y := 247, 165
	m, n := 5, 5

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

	var x0, y0, i int
	i = 1
	for y := 0; y < ysteps; y++ {

		for x := 0; x < xsteps; x++ {
			frstRect := image.Rect(x0, y0, x0+m, y0+n)
			secondRect := image.Rect(x0, y0,
				int(math.Min(float64(x0+m), float64(X))),
				int(math.Min(float64(y0+n), float64(Y))))

			fmt.Printf("%d. %v\n", i, frstRect.Intersect(secondRect))

			x0 += m
			i++
		}

		x0 = 0
		y0 += n

	}

}

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





type Colorizer interface {
	At(x, y int) color.Color
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

type mandelbrotColorizer struct {
	props ColorProps
}

type juliaColorizer struct {
	props ColorProps
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
			//return fractal.Ramp(i, iterations)
			return color.Gray{255 - uint8(contrast*i)}
		}
	}

	return color.Black
}

func (m MyImage) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	colorProps, colorizer, _ := m.requestParserFunc(req.FormValue("p"))
	fmt.Println("in ServeHTTP", colorizer)
	m.colorProps = colorProps
	m.Colorizer = colorizer

	var buf bytes.Buffer
	png.Encode(&buf, m)

	w.Header().Set("Content-Type", "image/png")
	w.Write(buf.Bytes())
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
			return color.Gray{255 - uint8(contrast*i)}
		}
	}
	return color.Black
}





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

func sendFilefragments(filePath string) {
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
		ajax.Chan <- distro.Fragment{X: block.Min.X, Y: block.Min.Y, URL: distro.Post(buf.Bytes())}
	}

	fmt.Println("End of sendFilefragments...")
}

// func sequentialFileDownloadHandler(c http.ResponseWriter, req *http.Request) {
// 	//filePath := string(`C:\development\workspace\Fleeetnet_main\WebContent\images\add.png`)
// 	filePath := string(`C:\Users\600131444\Pictures\vehicleImage.jpg`)
// 	go sendFilefragments(filePath) // to Client
// }

func sequentialFileDownloadHandler(c http.ResponseWriter, req *http.Request) {
	nodes.StartServer()

	ports := []string{"5000", "5001", "5002", "5003", "5004", "5005"}
	nodes.StartWorkers(ports)

	// go func() {
	// 	var res nodes.StartResult
	// 	nodes.Call(nodes.ServerPort, "Master.Start", nodes.StartArgs{true}, &res)
	// 	fmt.Println("Result ", res.ID)

	// }()

}

func distributedFileDownloadHandler(c http.ResponseWriter, req *http.Request) {
	ports := []string{"5000", "5001", "5002", "5003", "5004", "5005"}
	nodes.StartWorkers(ports)
	time.Sleep(time.Second)

	// go func() {
	// 	var res nodes.StartResult
	// 	nodes.Call(nodes.ServerPort, "Master.Start", nodes.StartArgs{true}, &res)
	// 	fmt.Println("Result ", res.ID)

	// }()

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

//var ChatPage = static.Serve("chat/chat.html")

func startServer() {
	fmt.Println("in startServer... before registering /")
	//http.Handle("/", fractal.MainPage)
	//http.Handle("/", ajax.LogPage)
	//http.Handle("/", chat.ChatPage)
	http.Handle("/", distro.MainPage)

	http.Handle("/string", String("Hello World String"))
	http.Handle("/struct", &Struct{"Hello", ":", "Go Course..."})

	http.Handle("/mandelbrot",
		MyImage{
			requestParserFunc: mandelbrotReqParser,
		})
	http.Handle("/julia",
		MyImage{
			requestParserFunc: juliaReqParser,
		})

	http.HandleFunc("/join", joinHandler)
	http.HandleFunc("/say", sayHandler)
	http.HandleFunc("/exit", exitHandler)

	// m := nodes.StartServer()
	// time.Sleep(time.Second)
	//ports := []string{"5000", "5001"}

	// m.KillWorkers()
	// time.Sleep(time.Second)
	// m.StopMaster()

	// nodes.Ch = make(chan struct{})
	// go func() {
	// 	os.Stdin.Read(make([]byte, 1))
	// 	close(nodes.Ch)
	// }()

	http.HandleFunc("/start", sequentialFileDownloadHandler)

	log.Fatal(http.ListenAndServe(":4000", nil))
}

type Master struct {
}

type Node struct {
}





