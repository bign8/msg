package buffer

// Struct is a recursive deserializer object
type Struct interface {
	Read(Buffer)
	Write(Buffer)

	// TODO: some kind of per-type or per-package version negotiation
}

// Buffer is an error safe buffer
type Buffer interface {
	Set(error)
	Err() error
	Len() int
	Cap() int
	Grow(int)

	WriteBool(bool)
	WriteByte(byte)
	WriteI16(int16)
	WriteI32(int32)
	WriteI64(int64)
	WriteDouble(float64)
	WriteString(string)
	WriteBinary([]byte)

	ReadBool() bool
	ReadByte() byte
	ReadI16() int16
	ReadI32() int32
	ReadI64() int64
	ReadDouble() float64
	ReadString() string
	ReadBinary() []byte
}
