// +build !js

package buff

import "errors"

func newInternal(b []byte) Buffer {
	return &buff{
		arr: b,
		ctx: make([]string, 0, 8),
	}
}

type buff struct {
	arr []byte
	err error
	ctx []string
	off int
}

func (b *buff) Bytes() []byte    { return b.arr }
func (b *buff) Err() error       { return b.err }
func (b *buff) Push(name string) { b.ctx = append(b.ctx, name) }
func (b *buff) Pop()             { b.ctx = b.ctx[:len(b.ctx)-2] }

func (b *buff) ReadBool() bool     { return false }
func (b *buff) ReadByte() byte     { return 0x00 }
func (b *buff) ReadInt() int       { return 0 }
func (b *buff) ReadFloat() float64 { return 0 }
func (b *buff) ReadStr() string    { return "" }
func (b *buff) ReadBytes() (a []byte) {
	if b.err != nil {
		return nil
	}
	b.Push("ReadBytes")
	n := b.ReadInt()
	if b.err == nil && b.off+n > len(b.arr) {
		b.fail("long read")
	}
	if b.err == nil {
		a = b.arr[b.off : b.off+n]
		b.off += n
	}
	b.Pop()
	return a
}
func (b *buff) ReadType(Type)     {}
func (b *buff) ReadRepeated() int { return 0 }

func (b *buff) WriteBool(bool)     {}
func (b *buff) WriteByte(byte)     {}
func (b *buff) WriteInt(int)       {}
func (b *buff) WriteFloat(float64) {}
func (b *buff) WriteStr(string)    {}
func (b *buff) WriteBytes(a []byte) {
	var l int
	if b.err == nil {
		l = len(a)
		b.WriteInt(l)
	}
	if b.err == nil {
		n := b.grow(l)
		copy(b.arr[n:], a)
	}
}
func (b *buff) WriteType(Type)    {}
func (b *buff) WriteRepeated(int) {}

func (b *buff) grow(n int) int {
	// try growing by reslicing
	if l := len(b.arr); n <= cap(b.arr)-l {
		b.arr = b.arr[:l+n]
		return l
	}

	panic("TODO: figure out a better way to grow")
}

func (b *buff) fail(msg string) {
	if b.err != nil {
		panic("already failed with: " + b.err.Error())
	}
	prefix := b.ctx[0]
	for _, s := range b.ctx[1:] {
		prefix += "." + s
	}
	b.err = errors.New(prefix + ": " + msg)
}
