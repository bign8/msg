package main

import (
	"net/http"
	"os"
)

func main() {
	dir := os.Getenv("GOPATH") + "/src/github.com/bign8/msg/ws/test"
	http.Handle("/", http.FileServer(http.Dir(dir)))

	// TODO: add watcher to generate client code
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}
}
