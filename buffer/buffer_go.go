package buffer

import (
	"bytes"
)

var _ Buffer = (*buffer)(nil)

type buffer struct {
	buff bytes.Buffer
	err  error
}

func (b *buffer) Set(e error) { b.err = e }
func (b *buffer) Err() error  { return b.err }
func (b *buffer) Len() int    { return b.buff.Len() }
func (b *buffer) Cap() int    { return b.buff.Cap() }
func (b *buffer) Grow(n int)  { b.buff.Grow(n) }

func (b *buffer) WriteBool(bool) {}
func (b *buffer) ReadBool() bool { return false }

func (b *buffer) WriteByte(c byte) {
	if b.err == nil {
		b.err = b.buff.WriteByte(c)
	}
}

func (b *buffer) ReadByte() (c byte) {
	if b.err == nil {
		c, b.err = b.buff.ReadByte()
	}
	return c
}

func (b *buffer) WriteI16(int16) {}
func (b *buffer) ReadI16() int16 { return 0 }

func (b *buffer) WriteI32(int32) {}
func (b *buffer) ReadI32() int32 { return 0 }

func (b *buffer) WriteI64(int64) {}
func (b *buffer) ReadI64() int64 { return 0 }

func (b *buffer) WriteDouble(d float64) {}
func (b *buffer) ReadDouble() float64   { return 0 }

func (b *buffer) WriteString(s string) {}
func (b *buffer) ReadString() string   { return "" }

func (b *buffer) WriteBinary([]byte) {}
func (b *buffer) ReadBinary() []byte { return nil }
