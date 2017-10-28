package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

var (
	port     = flag.Int("port", 3000, "port to serve on")
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
	errc := make(chan error)

	// Static File server
	dir := os.Getenv("GOPATH") + "/src/github.com/bign8/msg/ws/test"
	http.Handle("/", http.FileServer(http.Dir(dir)))

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
	})

	// TODO: add watcher to generate client code
	go func() {
		fmt.Println("TODO: watch generation list")
		// gopherjs build -m client.go
		// or
		// go generate ./...
	}()

	// The full HTTP server
	go func() {
		log.Printf("Listening on :%d\n", *port)
		errc <- http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	}()

	// Wait for something to fail
	check(<-errc, "main")
}
