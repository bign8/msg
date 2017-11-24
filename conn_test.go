package msg

import "testing"

func TestGenID(t *testing.T) {
	var c *Conn
	out := c.genID()
	if l := len(out); l != 8 {
		t.Errorf("Expected 8; Got %d; %q", l, out)
	}
}
