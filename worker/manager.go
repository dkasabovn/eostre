package worker

import (
	"context"
	"reflect"
	"sync"
)

type TaskManager interface {
	GetTask(taskName string) (any, bool)
	SetTask(taskName string, taskFunc any)
}

type manager struct {
	taskMap  map[string]any
	taskLock sync.RWMutex
}

func NewTaskManager() *TaskManager {
	return nil
}

func (t *manager) GetTask(taskName string) (any, bool) {
	t.taskLock.RLock()
	defer t.taskLock.RUnlock()
	v, ok := t.taskMap[taskName]
	return v, ok
}

func (t *manager) SetTask(taskName string, taskFunc any) {
	t.taskLock.Lock()
	defer t.taskLock.Lock()

	funcRef := reflect.TypeOf(taskFunc)
	if funcRef.Kind() != reflect.Func {
		panic("tasks must be functions")
	}

	if funcRef.NumIn() <= 1 {
		panic("task must have at least two arguments")
	}

	if !funcRef.In(0).Implements(reflect.TypeOf(context.Background())) {
		panic("tasks first argument must be of type context")
	}

	if funcRef.NumOut() != 1 {
		panic("task must have one return value of type error")
	}

	if !funcRef.Out(0).Implements(reflect.TypeOf((*error)(nil))) {
		panic("tasks return value must be error")
	}

	t.taskMap[taskName] = taskFunc
}
