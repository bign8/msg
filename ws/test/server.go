// +build ignore

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var (
	port     = flag.Int("port", 3001, "port to serve on")
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func check(err error, name string) {
	if err != nil {
		log.Printf("About to panic: %q", name)
		panic(err)
	}
}

func main() {
	flag.Parse()

	// Socket to test against
	http.HandleFunc("/ws/", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		check(err, "upgrader.Upgrade")
		switch r.URL.Path {
		case "/ws/immediate-close":
			check(conn.Close(), "immediate-close: Close")
		case "/ws/binary-static":
			err := conn.WriteMessage(websocket.BinaryMessage, []byte{0x00, 0x01, 0x02, 0x03, 0x04})
			check(err, "conn.WriteMessage")
		case "/ws/wait-30s":
			<-time.After(30 * time.Second)
		default: // echo server
			for {
				typ, p, err := conn.ReadMessage()
				check(err, "conn.ReadMessage")
				check(conn.WriteMessage(typ, p), "conn.WriteMessage")
			}
		}
		conn.Close()
	})

	log.Printf("Server.go: Listening on :%d\n", *port)
	check(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil), "ListenAndServe")
}
