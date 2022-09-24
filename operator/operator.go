package operator

import (
	"context"
	"encoding/json"
	"reflect"

	"eostre/packet"
)

type Operand any

type Operator struct {
	task    *packet.Task
	operand Operand
}

func NewOperator(operand Operand, task *packet.Task) *Operator {
	return &Operator{
		task:    task,
		operand: operand,
	}
}

func (o *Operator) Call() (err error) {
	defer func() {
		if e := recover(); e != nil {
			switch e := e.(type) {
			case error:
				err = e
			default:
				err = ErrTaskPanicked
			}
		}
	}()

	params, err := o.validateArguments()
	if err != nil {
		return err
	}

	callableOperand := reflect.ValueOf(o.operand)

	results := callableOperand.Call(params)

	if len(results) != 1 || !results[0].Type().Implements(reflect.TypeOf((*error)(nil))) {
		return ErrOperandResultMismatch
	}

	err = results[0].Interface().(error)
	return err
}

func (o *Operator) validateArguments() ([]reflect.Value, error) {
	reflectedFunction := reflect.TypeOf(o.operand)
	if reflectedFunction.Kind() != reflect.Func {
		return nil, ErrOperandIsNotFunction
	}

	// context.Context is automatically passed in to operands thus an operand should have one more argument than tasks
	if len(o.task.Args)+1 != reflectedFunction.NumIn() {
		return nil, NewArgumentNumberMismatchErr(len(o.task.Args)+1, reflectedFunction.NumIn())
	}

	// check that the first argument to an operand is context.Context
	if reflectedFunction.In(0).Kind() != reflect.Interface {
		return nil, ErrFirstArgNotContext
	}

	operandArgs := make([]reflect.Value, len(o.task.Args)+1)
	for i, taskArg := range o.task.Args {
		opArg := reflectedFunction.In(i + 1)
		value, err := convertValue(taskArg.Type, opArg, taskArg.Value)
		if err != nil {
			return nil, err
		}
		operandArgs[i+1] = value
	}

	operandArgs[0] = reflect.ValueOf(context.Background())

	return operandArgs, nil
}

func convertValue(argType packet.Type, opType reflect.Type, value any) (reflect.Value, error) {
	switch argType {
	case packet.Float64, packet.Float32:
		float, ok := value.(float64)
		if !ok {
			return reflect.Value{}, ErrParsingType
		}
		v := reflect.New(argType.ConvertToType())
		if v.Type() != opType {
			return reflect.Value{}, ErrParsingType
		}
		v.Elem().SetFloat(float)
		return v.Elem(), nil
	case packet.Int64, packet.Int32, packet.Int, packet.Int16, packet.Int8:
		intv, ok := value.(int64)
		if !ok {
			return reflect.Value{}, ErrParsingType
		}
		v := reflect.New(argType.ConvertToType())
		if v.Type() != opType {
			return reflect.Value{}, ErrParsingType
		}
		v.Elem().SetInt(intv)
		return v.Elem(), nil
	case packet.String:
		stringv, ok := value.(string)
		if !ok {
			return reflect.Value{}, ErrParsingType
		}
		v := reflect.New(argType.ConvertToType())
		if v.Type() != opType {
			return reflect.Value{}, ErrParsingType
		}
		v.Elem().SetString(stringv)
		return v.Elem(), nil
	case packet.Struct:
		bytesv, ok := value.([]byte)
		if !ok {
			return reflect.Value{}, ErrParsingType
		}
		v := reflect.New(opType)
		if err := json.Unmarshal(bytesv, v.Interface()); err != nil {
			return reflect.Value{}, ErrParsingType
		}
		return v.Elem(), nil
	case packet.Bool:
		boolv, ok := value.(bool)
		if !ok {
			return reflect.Value{}, ErrParsingType
		}
		v := reflect.New(argType.ConvertToType())
		if v.Type() != opType {
			return reflect.Value{}, ErrParsingType
		}
		v.Elem().SetBool(boolv)
		return v.Elem(), nil
	}
	return reflect.Value{}, ErrParsingType
}
