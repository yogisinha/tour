// Copyright 2010 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fractal

import (
	"net/http"
	"tour/static"
)

var extra = map[string]http.Handler{
	"/":             static.Serve("fractal/mandelbrot.html"),
	"/juliaviewer":  static.Serve("fractal/julia.html"),
	"/newtonviewer": static.Serve("fractal/newton.html"),
	"/behavior.js":  static.Serve("fractal/behavior.js"),
	"/center.gif":   static.Serve("fractal/center.gif"),
	"/gsv.js":       static.Serve("fractal/gsv.js"),
	"/none.png":     static.Serve("fractal/none.png"),
}

// MainPage serves the main fractal viewing HTML page.
var MainPage http.Handler = http.HandlerFunc(mainPage)

func mainPage(w http.ResponseWriter, req *http.Request) {
	h := extra[req.URL.Path]
	if h != nil {
		h.ServeHTTP(w, req)
		return
	}
	http.NotFound(w, req)
	w.WriteHeader(http.StatusNotFound)
}
