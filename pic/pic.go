// Copyright 2010 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pic

import (
	"encoding/json"
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"tour/stack"
	"tour/static"
)

const cacheSize = 20

var cache struct {
	sync.Mutex
	n    uint
	data [cacheSize][]byte
}

var byteFunc func(dx, dy int) [][]uint8
var imageFunc func(dx, dy int) image.Image

func runpic(dx, dy int) (pngdata []byte, err string) {
	defer func() {
		if v := recover(); v != nil {
			err = fmt.Sprintln("panic: ", v)
			err += "\nCall stack:\n" + stack.Stack()
		}
	}()

	var m image.Image
	if byteFunc != nil {
		data := byteFunc(dx, dy)
		if len(data) != dy {
			return nil, fmt.Sprintf("len(data) = %d, asked for %d", len(data), dy)
		}
		for i, row := range data {
			if len(row) != dx {
				return nil, fmt.Sprintf("len(data[%d]) = %d, asked for %d", i, len(row), dx)
			}
		}
		m1 := image.NewNRGBA(image.Rect(0, 0, dx, dy))
		for y := 0; y < dy; y++ {
			for x := 0; x < dx; x++ {
				v := data[y][x]
				m1.Set(y, x, color.RGBA{v, v, v, 255})
			}
		}
		m = m1
	} else {
		m = imageFunc(dx, dy)

		if m == nil {
			return nil, "image function returned nil"
		}
		if m.Bounds().Dx() != dx || m.Bounds().Dy() != dy {
			return nil, fmt.Sprintf("image function returned %dx%d, asked for %dx%d", m.Bounds().Dx(), m.Bounds().Dy(), dx, dy)
		}
	}

	var buf bytes.Buffer
	png.Encode(&buf, m)
	return buf.Bytes(), ""
}

func picHandler(c http.ResponseWriter, req *http.Request) {	
	dx, _ := strconv.Atoi(req.FormValue("dx"))
	dy, _ := strconv.Atoi(req.FormValue("dy"))
	if dx <= 0 || dy <= 0 {
		c.WriteHeader(404)
		c.Write([]byte("invalid dx, dy"))
		return
	}

	var result struct {
		Img   string
		Error string
	}
	data, err := runpic(dx, dy)
	if err != "" {		
		result.Error = "<pre>" + err + "</pre>"
	} else {		
		cache.Lock()
		n := cache.n
		cache.data[n%cacheSize] = data
		cache.n++
		cache.Unlock()
		result.Img = fmt.Sprintf(`/cache/%d.png`, n)
	}
	
	b, er := json.Marshal(result)	// only possible error is bad type and result is ok
	if er != nil {
		fmt.Printf("Error %v", er)
	}
	
	c.Header().Set("Content-Type", "application/json")
	c.Write(b)
}

func cacheHandler(c http.ResponseWriter, req *http.Request) {
	p := req.URL.Path[7:]
	if !strings.HasSuffix(p, ".png") {
		http.NotFound(c, req)
		return
	}
	p = p[0 : len(p)-4]
	n, err := strconv.Atoi(p)
	if err != nil {
		http.NotFound(c, req)
		return
	}
	cache.Lock()
	if n+cacheSize <= int(cache.n) {
		cache.Unlock()
		http.NotFound(c, req)
		return
	}
	data := cache.data[n%cacheSize]
	cache.Unlock()

	c.Header().Set("Content-Type", "image/png")
	c.Write(data)
}

// Serve runs a web server on port 4000 showing the
// 8-bit grayscale picture returned by f.
func Serve(f func(dx, dy int) [][]uint8) {
	byteFunc = f
	serve()
}

// ServeImage runs a web server on port 4000 showing the
// picture returned by f.
func ServeImage(f func(dx, dy int) image.Image) {
	imageFunc = f
	serve()
}

func serve() {
	http.Handle("/", static.Serve("pic/pic.html"))
	http.HandleFunc("/pic", picHandler)
	http.HandleFunc("/cache/", cacheHandler)
	fmt.Println("Listening on 127.0.0.1:4000")
	log.Fatal(http.ListenAndServe("127.0.0.1:4000", nil))

}
