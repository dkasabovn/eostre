package worker

import (
	"context"
	"os"
	"os/signal"
	"time"

	"eostre/operator"
	"eostre/packet"
)

type WorkerConfigFunc func(*worker)

type Worker interface {
	Start(context.Context)
	RunOnce(context.Context, *packet.Task)
}

func WithParallelism(a uint) WorkerConfigFunc {
	return func(w *worker) {
		w.paralellism = a
		w.semaphore = make(chan struct{}, a)
	}
}

func WithPollTimeout(t time.Duration) WorkerConfigFunc {
	return func(w *worker) {
		w.pollTimeout = t
	}
}

func WithOSNotify(sigs ...os.Signal) WorkerConfigFunc {
	return func(w *worker) {
		signal.Notify(w.osNotify, sigs...)
	}
}

type worker struct {
	manager     TaskManager
	paralellism uint
	queue       SimpleQueue
	pollTimeout time.Duration
	taskStream  chan *packet.Task
	semaphore   chan struct{}
	osNotify    chan os.Signal
	errorChan   chan error
}

func NewWorker(queue SimpleQueue, opts ...WorkerConfigFunc) Worker {
	w := &worker{
		manager:     *NewTaskManager(),
		paralellism: 1,
		queue:       queue,
		pollTimeout: time.Second * 5,
		taskStream:  make(chan *packet.Task),
		semaphore:   make(chan struct{}, 1),
		osNotify:    make(chan os.Signal),
		errorChan:   make(chan error),
	}

	for _, opt := range opts {
		opt(w)
	}

	return w
}

func (w *worker) read(ctx context.Context) {
	ticker := time.NewTicker(w.pollTimeout)
	for {
		select {
		case <-w.osNotify:
			return
		case <-ticker.C:
			tasks, err := w.queue.Read(ctx)
			if err != nil {
				w.errorChan <- err
				continue
			}

			for _, task := range tasks {
				w.taskStream <- task
			}
		}
	}
}

func (w *worker) RunOnce(ctx context.Context, task *packet.Task) {
	defer func() {
		<-w.semaphore
	}()

	taskFunction, ok := w.manager.GetTask(task.Signature)
	if !ok {
		panic("Could not find task " + task.Signature)
	}

	err := operator.NewOperator(taskFunction, task).Call()
	if err != nil {
		w.errorChan <- err
		return
	}
}

func (w *worker) Start(ctx context.Context) {
	go w.read(ctx)
	for {
		select {
		case task := <-w.taskStream:
			w.semaphore <- struct{}{}
			go w.RunOnce(ctx, task)
		case <-w.osNotify:
			close(w.semaphore)
			for {
				_, more := <-w.semaphore
				if !more {
					break
				}
			}
			return
		}
	}
}
