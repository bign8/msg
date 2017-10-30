package main

//go:generate gopherjs build -mv client.go

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
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

	// Socket management signal verification
	refresh := make(chan time.Time)
	register := make(chan chan<- time.Time)
	go func() {
		listeners := make(map[chan<- time.Time]bool)
		for {
			select {
			case client := <-register:
				listeners[client] = true
			case now := <-refresh:
				for client := range listeners {
					select {
					case client <- now:
					default:
						delete(listeners, client)
					}
				}
			}
		}
	}()

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
		case "/ws/watcher":
			// ticker := time.NewTicker(10 * time.Second)
			ticker := make(chan time.Time)
			register <- ticker
			defer close(ticker)
			for now := range ticker {
				if conn.WriteJSON(now) != nil {
					break
				}
			}
		default: // echo server
			for {
				typ, p, err := conn.ReadMessage()
				check(err, "conn.ReadMessage")
				check(conn.WriteMessage(typ, p), "conn.WriteMessage")
			}
		}
		conn.Close()
	})

	// TODO: add watcher to generate client code
	go func() {
		var last time.Time
		ticker := time.NewTicker(200 * time.Millisecond)
		defer ticker.Stop()
		for range ticker.C {
			info, err := os.Stat(dir + "/client.go")
			check(err, "os.Stat")
			if last != info.ModTime() {
				log.Printf("Client.go: Detected change, rebuilding client")
				cmd := exec.Command("go", "generate", "./...")
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				if err := cmd.Run(); err != nil {
					log.Printf("Failed to execute regen: %s", err)
				} else {
					log.Printf("Client.go: Rebuild complete.")
					last = info.ModTime()
					refresh <- last
				}
			}
		}
	}()

	// The full HTTP server
	go func() {
		log.Printf("Server.go: Listening on :%d\n", *port)
		errc <- http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	}()

	// Wait for something to fail
	check(<-errc, "main")
}
