package worker

import (
	"context"

	"eostre/packet"
)

// An interface to describe a queue
type SimpleQueue interface {
	Read(context.Context) ([]*packet.Task, error)
	Send(context.Context, *packet.Task) error
	SendMany(context.Context, []*packet.Task) error
}
