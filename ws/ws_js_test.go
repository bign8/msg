package ws

import (
	"runtime"
	"testing"

	"github.com/gopherjs/gopherjs/js"
	"github.com/stretchr/testify/assert"
)

func jsSetupFakeWS(t *testing.T) {
	if js.Global.Get("WebSocket") == js.Undefined {
		js.Global.Set("WebSocket", func() *js.Object {
			obj := js.Global.Get("Object").New()
			obj.Set("addEventListener", func() {})
			obj.Set("close", func() {})
			return obj
		})
	}
}

func TestTransport(t *testing.T) {
	tr := New("fake")
	assert.False(t, tr.Able())

	if runtime.Compiler == "gopherjs" {
		jsSetupFakeWS(t)
		assert.NoError(t, tr.Open())
		assert.True(t, tr.Able())
	}

	assert.NoError(t, tr.Kill())
}
