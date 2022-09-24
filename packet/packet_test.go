package packet_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"eostre/packet"
)

func TestArgFloatValid(t *testing.T) {
	t.Parallel()

	arg := &packet.Arg{
		Type:  packet.Float64,
		Value: 1.0,
	}

	task, err := packet.NewTask("asdf", arg)
	assert.NoError(t, err)
	assert.Equal(t, task.Signature, "asdf")
	assert.Len(t, task.Args, 1)
}

func TestArgFloatInvalid(t *testing.T) {
	t.Parallel()

	arg := &packet.Arg{
		Type:  packet.Float64,
		Value: "",
	}

	_, err := packet.NewTask("asdf", arg)
	assert.Error(t, err)
}

func TestIntValid(t *testing.T) {
	t.Parallel()

	arg := &packet.Arg{
		Type:  packet.Int64,
		Value: int64(1),
	}

	task, err := packet.NewTask("asdf", arg)
	assert.NoError(t, err)
	assert.Equal(t, task.Signature, "asdf")
	assert.Len(t, task.Args, 1)
}

func TestIntInvalid(t *testing.T) {
	t.Parallel()

	arg := &packet.Arg{
		Type:  packet.Int64,
		Value: "howdy",
	}

	_, err := packet.NewTask("asdf", arg)
	assert.Error(t, err)
}

func TestStringValid(t *testing.T) {
	t.Parallel()

	arg := &packet.Arg{
		Type:  packet.String,
		Value: "howdy",
	}

	task, err := packet.NewTask("asdf", arg)
	assert.NoError(t, err)
	assert.Equal(t, task.Signature, "asdf")
	assert.Len(t, task.Args, 1)
}

func TestStringInvalid(t *testing.T) {
	t.Parallel()

	arg := &packet.Arg{
		Type:  packet.String,
		Value: 1,
	}

	_, err := packet.NewTask("asdf", arg)
	assert.Error(t, err)
}

func TestBoolValid(t *testing.T) {
	t.Parallel()

	arg := &packet.Arg{
		Type:  packet.Bool,
		Value: true,
	}

	task, err := packet.NewTask("asdf", arg)
	assert.NoError(t, err)
	assert.Equal(t, task.Signature, "asdf")
	assert.Len(t, task.Args, 1)
}

func TestBoolInvalid(t *testing.T) {
	t.Parallel()

	arg := &packet.Arg{
		Type:  packet.Bool,
		Value: 1,
	}

	_, err := packet.NewTask("asdf", arg)
	assert.Error(t, err)
}

func TestBytesValid(t *testing.T) {
	t.Parallel()

	arg := &packet.Arg{
		Type:  packet.Bytes,
		Value: []byte("howdy"),
	}

	task, err := packet.NewTask("asdf", arg)
	assert.NoError(t, err)
	assert.Equal(t, task.Signature, "asdf")
	assert.Len(t, task.Args, 1)
	assert.Equal(t, string(task.Args[0].Value.([]byte)), "howdy")
}

func TestBytesInvalid(t *testing.T) {
	t.Parallel()

	arg := &packet.Arg{
		Type:  packet.Bytes,
		Value: "howdy",
	}

	_, err := packet.NewTask("asdf", arg)
	assert.Error(t, err)
}

func TestStructValid(t *testing.T) {
	t.Parallel()

	structT := struct {
		Test int
	}{
		Test: 1,
	}

	arg := &packet.Arg{
		Type:  packet.Struct,
		Value: structT,
	}

	marsh, err := json.Marshal(structT)
	assert.NoError(t, err)

	task, err := packet.NewTask("asdf", arg)
	assert.NoError(t, err)
	assert.Equal(t, task.Signature, "asdf")
	assert.Len(t, task.Args, 1)
	assert.Equal(t, json.RawMessage(marsh), task.Args[0].Value.(json.RawMessage))
}

func TestStructInvalid(t *testing.T) {
	t.Parallel()

	arg := &packet.Arg{
		Type:  packet.Struct,
		Value: "asdf",
	}

	_, err := packet.NewTask("asdf", arg)
	assert.Error(t, err)
}
