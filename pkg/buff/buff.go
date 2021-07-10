package buff

import "github.com/bign8/msg"

// New constructs a Buffer for packing/unpacking info
func New(b []byte) Buffer {
	return newInternal(b)
}

// NewReader ...
func NewReader(msg.Reader) Buffer {
	return nil
}

// NewWriter ...
func NewWriter(msg.Writer) Buffer {
	return nil
}

// Struct is a base serializable type.
type Struct interface {
	Read(Buffer)
	Write(Buffer)
}

// Buffer is a base object for reading and writing msgs to.
type Buffer interface {
	Bytes() []byte // Resulting buffer state
	Err() error    // Any errors encountered when processing
	Push(string)   // adds a info to errors
	Pop()          // remove error info
	Set(error)
	Len() int
	Cap() int
	Grow(int)

	Reader
	Writer
}

type Reader interface {
	ReadBinary() []byte
	ReadBool() bool
	ReadByte() byte
	ReadDouble() float64
	// ReadI16() int16
	// ReadI32() int32
	// ReadI64() int64
	ReadInt() int
	ReadRepeated() (size int) // list or map
	ReadString() string
	ReadStruct(Struct)
}

type Writer interface {
	WriteBinary([]byte)
	WriteBool(bool)
	WriteByte(byte)
	WriteDouble(float64)
	// WriteI16(int16)
	// WriteI32(int32)
	// WriteI64(int64)
	WriteInt(int)
	WriteRepeated(size int) // list or map
	WriteString(string)
	WriteStruct(Struct)
}
