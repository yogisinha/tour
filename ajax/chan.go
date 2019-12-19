// The ajax package lets Go programs send messages to a web browser using a channel.
package ajax

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"time"
)

// Chan is a channel connected to the client web browsers.
// Messages sent on chan are delivered to JavaScript running in the browser.
var Chan chan<- interface{}

func handler(c http.ResponseWriter, req *http.Request) {
	ch := make(chan []byte, 1)
	Chan <- pending{c, ch, int64(time.Nanosecond + 10e9), req}
	c.Header().Set("Content-Type", "application/json")
	c.Write(<-ch)
}

type pending struct {
	conn    http.ResponseWriter
	ch      chan []byte
	timeout int64
	Req     *http.Request
}

func init() {
	http.HandleFunc("/_ajaxchan_", handler)
	ch := make(chan interface{}, 1000)
	Chan = ch
	go loop(ch)
}

var timeout = []byte(`{"timeout": true}`)

func loop(ch chan interface{}) {
	data := make(map[string][]byte)
	wait := make(map[string]map[chan []byte]pending)
	n := 0
	tick := time.Tick(10e9)

	for {
		select {
		case <-tick:
			now := int64(time.Now().Nanosecond())
			for _, w := range wait {
				for ch, p := range w {
					if p.timeout < now {
						ch <- timeout
						delete(w, ch)
					}
				}
			}

		case r := <-ch:
			switch m := r.(type) {
			case pending:
				if m.Req.FormValue("poll") != "" {
					b := []byte(fmt.Sprintf(`{"N": %d}`, n))
					m.ch <- b
					break
				}
				name := m.Req.FormValue("n")
				if name == "" {
					m.ch <- nil
					break
				}
				b := data[name]
				if b != nil {
					m.ch <- b
					break
				}
				w := wait[name]
				if w == nil {
					w = make(map[chan []byte]pending)
					wait[name] = w
				}
				w[m.ch] = m

			default:
				typ := reflect.TypeOf(m).Name()
				if typ == "" {
					typ = "?"
				}
				b, err := json.Marshal(map[string]interface{}{typ: m})
				if err != nil {
					log.Fatal("json.Marshal: ", err)
					break
				}
				name := strconv.Itoa(n)
				n++
				data[name] = b
				w := wait[name]
				if w != nil {
					delete(wait, name)
					for ch := range w {
						ch <- b
					}
				}
			}
		}
	}
}
