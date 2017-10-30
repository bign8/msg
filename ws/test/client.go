// +build js

package main

import (
	"flag"
	"testing"
	"time"

	"github.com/bign8/msg/ws"
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/websocket/websocketjs"
)

var loc = js.Global.Get("document").Get("location").Get("href").Call("replace", "http", "ws", 1).String() + "ws/"

func TestRealThing(t *testing.T) {
	t.Log("Opening Socket 1")
	print("Opening Socket 2")
	sock := ws.New(loc)
	time.Sleep(time.Second)
	print("Verifying Socket is Able")
	if sock.Able() {
		panic("Should have to open this bad boy first")
	}
}

func main() {
	flag.Set("test.v", "true")

	// Wait for a refresh
	sock, err := websocketjs.New(loc + "watcher")
	if err == nil {
		die := func(msg *js.Object) {
			js.Global.Get("document").Get("location").Call("reload")
		}
		sock.AddEventListener("message", false, die)
		sock.AddEventListener("close", false, die)
		sock.AddEventListener("error", false, die)
	} else {
		print("Could not open watcher socket")
	}

	// Run the actual tests
	testing.Main(func(pat string, str string) (bool, error) {
		return true, nil
	}, []testing.InternalTest{
		{
			Name: "TestRealThing",
			F:    TestRealThing,
		},
	}, nil, nil)
}
