package packet

import (
	"fmt"
	"reflect"
)

type Type string

type ErrInvalidType struct {
	expected string
	got      string
}

const (
	Float64 Type = "float64"
	Float32 Type = "float32"
	Int64   Type = "int64"
	Int32   Type = "int32"
	Int     Type = "int"
	Int16   Type = "int16"
	Int8    Type = "int8"
	String  Type = "string"
	Bytes   Type = "[]byte"
	Struct  Type = "struct"
	Bool    Type = "bool"
)

func (t Type) ConvertToType() reflect.Type {
	switch t {
	case Float64:
		return reflect.TypeOf(float64(1))
	case Float32:
		return reflect.TypeOf(float32(1))
	case Int64:
		return reflect.TypeOf(int64(1))
	case Int32:
		return reflect.TypeOf(int32(1))
	case Int:
		return reflect.TypeOf(int(1))
	case Int16:
		return reflect.TypeOf(int16(1))
	case Int8:
		return reflect.TypeOf(int8(1))
	case String:
		return reflect.TypeOf("")
	case Bytes:
		return reflect.TypeOf([]byte{})
	case Bool:
		return reflect.TypeOf(bool(false))
	}
	return nil
}

func NewInvalidTypeErr(expected, got string) *ErrInvalidType {
	return &ErrInvalidType{
		expected: expected,
		got:      got,
	}
}

func (e *ErrInvalidType) Error() string {
	return fmt.Sprintf("invalid type, expected: %s, got: %s", e.expected, e.got)
}
