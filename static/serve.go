package static

import (
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path"
	"unicode/utf8"
)

type fileServer string

func (fsrv fileServer) ServeHTTP(conn http.ResponseWriter, req *http.Request) {
	name := string(fsrv)
	f, err := os.Open(name)
	if err != nil {
		log.Fatal(err)
		conn.WriteHeader(404)
		return
	}

	if ctype := mime.TypeByExtension(path.Ext(name)); ctype != "" {
		conn.Header().Set("Content-Type", ctype)
	} else {
		// read first chunk to decide between utf-8 text and binary
		var buf [1024]byte
		n, _ := io.ReadFull(f, buf[0:])
		b := buf[0:n]
		if isText(b) {
			conn.Header().Set("Content-Type", "text-plain; charset=utf-8")
		} else {
			conn.Header().Set("Content-Type", "application/octet-stream")
		}
		conn.Write(b)
	}
	io.Copy(conn, f)
}

// Serve returns an http.Handler that serves the given static file,
// loaded from src/pkg/tour.
func Serve(file string) http.Handler {
	// goroot := os.Getenv("GOROOT")
	// name := goroot + "/src/pkg/tour/" + file

	goroot := os.Getenv("GOPATH")
	name := goroot[:len(goroot)] + "\\src\\tour\\" + file

	fmt.Println(name)
	//name = "C:\\Yogesh\\development\\goworkspace\\gopl.io\\src\\tour\\pic\\pic.html"

	if access(name) {
		return fileServer(name)
	}
	log.Fatal("Cannot find tour/" + file)
	return nil
}

func access(name string) bool {
	f, err := os.Open(name)
	if err != nil {
		return false
	}
	f.Close()
	return true
}

func isText(b []byte) bool {
	for len(b) > 0 && utf8.FullRune(b) {
		rune, size := utf8.DecodeRune(b)
		if size == 1 && rune == utf8.RuneError {

			return false
		}
		if 0x80 <= rune && rune <= 0x9F {
			return false
		}
		if rune < ' ' {
			switch rune {
			case '\n', '\r', '\t':

			default:

				return false
			}
		}
		b = b[size:]
	}
	return true
}
