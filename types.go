package msg

import "errors"

var (
	_ Type = Slice{}
	_ Type = Map{}
	_ Type = Ptr{}
	_ Type = Struct{}

	Bool       Type = kind(1)
	Byte       Type = kind(2)
	Complex128 Type = kind(3)
	Complex64  Type = kind(4)
	Error      Type = kind(5)
	Float32    Type = kind(6)
	Float64    Type = kind(7)
	Int        Type = kind(8)
	Int16      Type = kind(9)
	Int32      Type = kind(10)
	Int64      Type = kind(11)
	Int8       Type = kind(012)
	Rune       Type = kind(13)
	String     Type = kind(14)
	Uint       Type = kind(16)
	Uint16     Type = kind(17)
	Uint32     Type = kind(18)
	Uint64     Type = kind(19)
	Uint8      Type = kind(20)
)

type Type interface {
	ReadFrom(in Reader) (o interface{}, err error)
	WriteTo(o interface{}, out Writer) error
}

type Slice struct {
	Value Type
}

func (s Slice) ReadFrom(in Reader) (o interface{}, err error) { return nil, errors.New("TODO") }
func (s Slice) WriteTo(o interface{}, out Writer) error       { return errors.New("TODO") }

type Map struct {
	Key, Value Type
}

func (s Map) ReadFrom(in Reader) (o interface{}, err error) { return nil, errors.New("TODO") }
func (s Map) WriteTo(o interface{}, out Writer) error       { return errors.New("TODO") }

type Ptr struct {
	Value Type
}

func (s Ptr) ReadFrom(in Reader) (o interface{}, err error) { return nil, errors.New("TODO") }
func (s Ptr) WriteTo(o interface{}, out Writer) error       { return errors.New("TODO") }

type Struct struct {
	Fields []Type
}

func (s Struct) ReadFrom(in Reader) (o interface{}, err error) { return nil, errors.New("TODO") }
func (s Struct) WriteTo(o interface{}, out Writer) error       { return errors.New("TODO") }

type kind uint8

func (k kind) ReadFrom(in Reader) (o interface{}, err error) { return nil, errors.New("TODO") }
func (k kind) WriteTo(o interface{}, out Writer) error       { return errors.New("TODO") }
