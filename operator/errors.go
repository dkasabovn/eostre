package operator

import (
	"errors"
	"fmt"
)

type ErrArgumentNumberMismatch struct {
	expected int
	got      int
}

var (
	ErrOperandIsNotFunction  = errors.New("operand is not a function")
	ErrTaskPanicked          = errors.New("task panicked")
	ErrFirstArgNotContext    = errors.New("operand first argument must be of type context.Context or interface")
	ErrParsingType           = errors.New("arg type and value does not match")
	ErrOperandResultMismatch = errors.New("operands must only return error")
)

func NewArgumentNumberMismatchErr(expected, got int) error {
	return &ErrArgumentNumberMismatch{
		expected: expected,
		got:      got,
	}
}

func (e *ErrArgumentNumberMismatch) Error() string {
	return fmt.Sprintf("expected %d arguments to the function, got %d", e.expected, e.got)
}
