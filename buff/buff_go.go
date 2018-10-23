// +build !js

package buff

func newInternal(b []byte) Buffer {
	return &buff{
		b: b,
		s: make([]string, 0, 8),
	}
}

type buff struct {
	b []byte
	e error
	s []string
}

func (b *buff) Bytes() []byte { return b.b }
func (b *buff) Err() error { return b.e }
func (b *buff) Push(s string) { b.s = append(b.s, s) }
func (b *buff) Pop() { b.s = b.s[:len(b.s)-2] }


func (b *buff) ReadBool() bool { return false }
func (b *buff) ReadByte() byte { return 0x00 }
func (b *buff) ReadInt() int { return 0 }
func (b *buff) ReadFloat() float64 { return 0 }
func (b *buff) ReadStr() string { return "" }
func (b *buff) ReadBytes() []byte { return nil }
func (b *buff) ReadType(Type) {}
func (b *buff) ReadRepeated() int { return 0 }

func (b *buff) WriteBool(bool) {}
func (b *buff) WriteByte(byte) {}
func (b *buff) WriteInt(int) {}
func (b *buff) WriteFloat(float64) {}
func (b *buff) WriteStr(string) {}
func (b *buff) WriteBytes([]byte) {}
func (b *buff) WriteType(Type) {}
func (b *buff) WriteRepeated(int) {}
