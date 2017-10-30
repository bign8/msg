package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"time"

	"github.com/gorilla/websocket"
	"github.com/yhat/wsutil"
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
	unregister := make(chan chan time.Time)
	var server *exec.Cmd
	defer func() {
		server.Process.Kill()
	}()

	// Statis serve everything!
	dir := os.Getenv("GOPATH") + "/src/github.com/bign8/msg/ws/test"
	http.Handle("/", http.FileServer(http.Dir(dir)))

	// Reverse proxy for web-sockets
	serv, _ := url.Parse("http://localhost:3001")
	// rp := httputil.NewSingleHostReverseProxy(serv)
	rp := wsutil.NewSingleHostReverseProxy(serv)
	http.Handle("/ws/", rp)

	go watch("client.go", updates, func() error {
		cmd := exec.Command("go", "generate", "./...")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	})
	go watch("server.go", updates, func() error {
		if server != nil {
			server.Process.Kill()
		}
		server = exec.Command("go", "run", dir+"/server.go")
		server.Stdout = os.Stdout
		server.Stderr = os.Stderr
		return server.Start()
	})
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
