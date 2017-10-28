// +build js

//go:generate gopherjs build -m client.go

package main

import (
	"flag"
	"testing"
	"time"

	"github.com/bign8/msg/ws"
	"github.com/gopherjs/gopherjs/js"
)

func TestRealThing(t *testing.T) {
	t.Log("Opening Socket 1")
	print("Opening Socket 2")
	loc := js.Global.Get("document").Get("location").Get("href").Call("replace", "http", "ws", 1).String() + "ws"
	sock := ws.New(loc)
	time.Sleep(time.Second)
	print("Verifying Socket is Able")
	if sock.Able() {
		panic("Should have to open this bad boy first")
	}
}

func main() {
	flag.Set("test.v", "true")
	testing.Main(func(pat string, str string) (bool, error) {
		return true, nil
	}, []testing.InternalTest{
		{
			Name: "TestRealThing",
			F:    TestRealThing,
		},
	}, nil, nil)
}
