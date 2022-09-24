package packet

import (
	"encoding/json"
	"reflect"
)

type Arg struct {
	Type  Type
	Value any
}

type Task struct {
	Signature string
	Args      []*Arg
}

func NewTask(signature string, args ...*Arg) (*Task, error) {
	for _, arg := range args {
		if err := arg.assert(); err != nil {
			return nil, err
		}
	}
	return &Task{
		Signature: signature,
		Args:      args,
	}, nil
}

func (a *Arg) assert() error {
	val := reflect.ValueOf(a.Value)
	switch a.Type {
	case Bytes:
		bytes, ok := a.Value.([]byte)
		if !ok {
			return NewInvalidTypeErr(string(a.Type), val.Kind().String())
		}
		a.Value = bytes
		return nil
	case Struct:
		if val.Kind() != reflect.Struct {
			return NewInvalidTypeErr(string(a.Type), val.Kind().String())
		}
		bytes, err := json.Marshal(a.Value)
		if err != nil {
			return err
		}
		a.Value = json.RawMessage(bytes)
		return nil
	default:
		if val.Kind().String() != string(a.Type) {
			return NewInvalidTypeErr(string(a.Type), val.Kind().String())
		}
		return nil
	}
}

func (t *Task) Serialize() ([]byte, error) {
	return json.Marshal(t)
}
