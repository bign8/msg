package main

//go:generate gopherjs build -mv client.go

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"sync"
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

	// Allow broadcasting to multiple routines at once (not really protecting something)
	cond := sync.NewCond((&sync.RWMutex{}).RLocker())

	// Statis serve everything!
	dir := os.Getenv("GOPATH") + "/src/github.com/bign8/msg/ws/test"
	http.Handle("/", http.FileServer(http.Dir(dir)))

	// Watch client
	go func() {
		var last time.Time
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()
		for range ticker.C {
			info, err := os.Stat(dir + "/client.go")
			check(err, "os.Stat")
			if mod := info.ModTime(); last != mod {
				log.Print("Starting Rebuild...")
				cmd := exec.Command("go", "generate", "./...")
				cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
				check(cmd.Run(), "onchange")
				log.Print("Rebuild Complete.")
				last = mod
				cond.Broadcast()
			}
		}
	}()

	// Watcher socket
	http.HandleFunc("/ws/watcher", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		check(err, "upgrader.Upgrade")
		for {
			cond.L.Lock()
			cond.Wait() // https://golang.org/pkg/sync/#Cond.Wait
			cond.L.Unlock()
			if conn.WriteJSON(time.Now()) != nil {
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
				if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
					return
				}
				check(err, "conn.ReadMessage")
				check(conn.WriteMessage(typ, p), "conn.WriteMessage")
			}
		}
		conn.Close()
	})

	log.Printf("Listening on :%d\n", *port)
	check(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil), "ListenAndServe")
}
