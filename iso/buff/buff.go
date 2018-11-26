package buff

// New constructs a Buffer for packing/unpacking info
func New(b []byte) Buffer {
	return newInternal(b)
}

// Type is a base serializable type.
type Type interface {
	Read(Buffer)
	Write(Buffer)
}

// Buffer is a base object for reading and writing msgs to.
type Buffer interface {
	Bytes() []byte // Resulting buffer state
	Err() error    // Any errors encountered when processing
	Push(string)   // adds a info to errors
	Pop()          // remove error info

	// Read Methods
	ReadBool() bool
	ReadByte() byte
	ReadInt() int
	ReadFloat() float64
	ReadStr() string
	ReadBytes() []byte
	ReadType(Type)
	ReadRepeated() (size int) // list or map

	// Write Methods
	WriteBool(bool)
	WriteByte(byte)
	WriteInt(int)
	WriteFloat(float64)
	WriteStr(string)
	WriteBytes([]byte)
	WriteType(Type)
	WriteRepeated(size int) // list or map
}
