// +build js

package main

import (
	"flag"
	"testing"
	"time"

	"github.com/bign8/msg/pkg/ws"
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/websocket/websocketjs"
)

var loc = js.Global.Get("document").Get("location").Get("href").Call("replace", "http", "ws", 1).String() + "ws/"

func TestImmediateClose(t *testing.T) {
	sock := ws.New(loc + "immediate-close")
	if sock.Able() {
		panic("Should have to open this bad boy first")
	}
	if err := sock.Open(); err != nil {
		t.Errorf("Could not open socket: %s", err)
	}
	if !sock.Able() {
		t.Errorf("Client should be able for a bit")
	}
	select {
	case <-sock.Wait(): // auto-closing socket
	case <-time.After(time.Second):
		t.Fatal("Socket didn't close after 1 second")
	}
	if sock.Able() {
		t.Errorf("Socket should be dead by now")
	}
}

func main() {
	flag.Set("test.v", "true") // verbose unit test logs

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
		{Name: "TestImmediateClose", F: TestImmediateClose},
		{Name: "TestImmediateClose", F: TestImmediateClose},
		{Name: "TestImmediateClose", F: TestImmediateClose},
		{Name: "TestImmediateClose", F: TestImmediateClose},
		{Name: "TestImmediateClose", F: TestImmediateClose},
		{Name: "TestImmediateClose", F: TestImmediateClose},
		{Name: "TestImmediateClose", F: TestImmediateClose},
		{Name: "TestImmediateClose", F: TestImmediateClose},
		{Name: "TestImmediateClose", F: TestImmediateClose},
		{Name: "TestImmediateClose", F: TestImmediateClose},
		{Name: "TestImmediateClose", F: TestImmediateClose},
		{Name: "TestImmediateClose", F: TestImmediateClose},
		{Name: "TestImmediateClose", F: TestImmediateClose},
		{Name: "TestImmediateClose", F: TestImmediateClose},
		{Name: "TestImmediateClose", F: TestImmediateClose},
		{Name: "TestImmediateClose", F: TestImmediateClose},
		{Name: "TestImmediateClose", F: TestImmediateClose},
		{Name: "TestImmediateClose", F: TestImmediateClose},
		{Name: "TestImmediateClose", F: TestImmediateClose},
		{Name: "TestImmediateClose", F: TestImmediateClose},
		{Name: "TestImmediateClose", F: TestImmediateClose},
		{Name: "TestImmediateClose", F: TestImmediateClose},
		{Name: "TestImmediateClose", F: TestImmediateClose},
		{Name: "TestImmediateClose", F: TestImmediateClose},
		{Name: "TestImmediateClose", F: TestImmediateClose},
	}, nil, nil)
}
