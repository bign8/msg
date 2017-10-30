package main

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

//go:generate gopherjs build -mv client.go

var (
	port     = flag.Int("port", 3000, "port to serve on")
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func main() {
	updates := make(chan time.Time)
	register := make(chan chan time.Time)

	go watch("client.go", updates, func() error {
		cmd := exec.Command("go", "generate", "./...")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	})
	go watch("server.go", updates, func() error {
		log.Println("TODO: refresh running server")
		return nil
	})
	go func() {
		listeners := make(map[chan<- time.Time]bool)
		for {
			select {
			case client := <-register:
				listeners[client] = true
			case now := <-updates:
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

	// Statis serve everything!
	dir := os.Getenv("GOPATH") + "/src/github.com/bign8/msg/ws/test"
	http.Handle("/", http.FileServer(http.Dir(dir)))

	// Perform updates and client refreshes
	log.Printf("Watcher: Serving on :%d", *port)
	check(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil), "ListenAndServe")
}

func check(err error, name string) {
	if err != nil {
		log.Printf("About to panic: %q", name)
		panic(err)
	}
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
