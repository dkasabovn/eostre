package operator_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"eostre/operator"
	"eostre/packet"
)

// Test basic assertions

func TestOperatorNotAFunction(t *testing.T) {
	t.Parallel()

	task, err := packet.NewTask("asdf", &packet.Arg{
		Type:  packet.Bool,
		Value: false,
	})

	assert.NoError(t, err)

	op := operator.NewOperator("string", task)

	err = op.Call()
	assert.Error(t, err)
	assert.ErrorIs(t, err, operator.ErrOperandIsNotFunction)
}

func TestOperatorArgCountMismatch(t *testing.T) {
	t.Parallel()

	task, err := packet.NewTask("asdf", &packet.Arg{
		Type:  packet.Bool,
		Value: false,
	})

	assert.NoError(t, err)

	op := operator.NewOperator(func(bool) {}, task)

	err = op.Call()
	assert.Error(t, err)
	mismatchErr := &operator.ErrArgumentNumberMismatch{}
	assert.ErrorAs(t, err, &mismatchErr)
}

func TestOperatorFirstArgNotContext(t *testing.T) {
	t.Parallel()

	task, err := packet.NewTask("asdf", &packet.Arg{
		Type:  packet.Bool,
		Value: false,
	})

	assert.NoError(t, err)

	op := operator.NewOperator(func(interface{}, bool) {}, task)

	err = op.Call()
	assert.Error(t, err)
	assert.ErrorIs(t, err, operator.ErrFirstArgNotContext)
}

func TestOperatorNoReturn(t *testing.T) {
	t.Parallel()

	task, err := packet.NewTask("asdf", &packet.Arg{
		Type:  packet.Bool,
		Value: false,
	})

	assert.NoError(t, err)

	op := operator.NewOperator(func(context.Context, bool) {}, task)

	err = op.Call()
	assert.Error(t, err)
	assert.ErrorIs(t, err, operator.ErrOperandResultMismatch)
}

// Test Various Types

func TestOperatorBoolValid(t *testing.T) {
	t.Parallel()

	task, err := packet.NewTask("asdf", &packet.Arg{
		Type:  packet.Bool,
		Value: false,
	})

	assert.NoError(t, err)

	op := operator.NewOperator(func(context.Context, bool) error { return nil }, task)

	err = op.Call()
	assert.NoError(t, err)
}

func TestOperatorBoolInvalid(t *testing.T) {
	t.Parallel()

	task, err := packet.NewTask("asdf", &packet.Arg{
		Type:  packet.Bool,
		Value: false,
	})

	assert.NoError(t, err)

	task.Args[0] = &packet.Arg{
		Type:  packet.Bool,
		Value: "asdf",
	}

	op := operator.NewOperator(func(context.Context, bool) error { return nil }, task)

	err = op.Call()
	assert.Error(t, err)
	assert.ErrorIs(t, err, operator.ErrParsingType)
}

func TestOperatorIntValid(t *testing.T) {
	t.Parallel()

	task, err := packet.NewTask("asdf", &packet.Arg{
		Type:  packet.Int,
		Value: 1,
	})

	assert.NoError(t, err)

	op := operator.NewOperator(func(context.Context, int) error { return nil }, task)

	err = op.Call()
	assert.NoError(t, err)
}

func TestOperatorIntInvalid(t *testing.T) {
	t.Parallel()

	task, err := packet.NewTask("asdf", &packet.Arg{
		Type:  packet.Int,
		Value: 1,
	})

	assert.NoError(t, err)

	task.Args[0] = &packet.Arg{
		Type:  packet.Int,
		Value: 0.0,
	}

	op := operator.NewOperator(func(context.Context, int) error { return nil }, task)

	err = op.Call()
	assert.Error(t, err)
	assert.ErrorIs(t, err, operator.ErrParsingType)
}

func TestOperatorFloatValid(t *testing.T) {
	t.Parallel()

	task, err := packet.NewTask("asdf", &packet.Arg{
		Type:  packet.Float32,
		Value: float32(1.0),
	})

	assert.NoError(t, err)

	op := operator.NewOperator(func(context.Context, float32) error { return nil }, task)

	err = op.Call()
	assert.NoError(t, err)
}

func TestOperatorFloatInvalid(t *testing.T) {
	t.Parallel()

	task, err := packet.NewTask("asdf", &packet.Arg{
		Type:  packet.Float32,
		Value: float32(1.0),
	})

	assert.NoError(t, err)

	task.Args[0] = &packet.Arg{
		Type:  packet.Float32,
		Value: 123123,
	}

	op := operator.NewOperator(func(context.Context, float32) error { return nil }, task)

	err = op.Call()
	assert.Error(t, err)
	assert.ErrorIs(t, err, operator.ErrParsingType)
}

// TODO: Test all other packet.Types and test byte + struct mismatch errors
