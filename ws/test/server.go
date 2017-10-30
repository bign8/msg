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
	updates := make(chan time.Time)
	register := make(chan chan time.Time)
	unregister := make(chan chan time.Time)

	// Statis serve everything!
	dir := os.Getenv("GOPATH") + "/src/github.com/bign8/msg/ws/test"
	http.Handle("/", http.FileServer(http.Dir(dir)))

	// Watch client
	go watch("client.go", updates, func() error {
		cmd := exec.Command("go", "generate", "./...")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	})

	// Signal to all sockets
	go func() {
		listeners := make(map[chan<- time.Time]bool)
		for {
			select {
			case client := <-register:
				listeners[client] = true
			case client := <-unregister:
				delete(listeners, client)
			case now := <-updates:
				for client := range listeners {
					client <- now
				}
			}
		}
	}()

	// Watcher socket
	http.HandleFunc("/ws/watcher", func(w http.ResponseWriter, r *http.Request) {
		conn, _ := upgrader.Upgrade(w, r, nil)
		ticker := make(chan time.Time)
		register <- ticker
		defer close(ticker)
		for now := range ticker {
			if conn.WriteJSON(now) != nil {
				break
			}
		}
		conn.Close()
	})

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

func watch(fn string, signal chan<- time.Time, onchange func() error) {
	log.Printf("Watcher: %q", fn)
	full := os.Getenv("GOPATH") + "/src/github.com/bign8/msg/ws/test/" + fn
	var last time.Time
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	for range ticker.C {
		info, err := os.Stat(full)
		check(err, "os.Stat")
		if mod := info.ModTime(); last != mod {
			log.Printf("%s: Starting Rebuild...", fn)
			check(onchange(), "onchange")
			log.Printf("%s: Rebuild Complete.", fn)
			last = mod
			signal <- mod
		}
	}
}
