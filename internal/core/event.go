package core

import "context"

// Bus is a synchronous event bus.
type Bus interface {
	PlayerDeleted(ctx context.Context, evt EventPlayerDeleted)
	OnPlayerDeleted(func(ctx context.Context, evt EventPlayerDeleted) error) func()
}

type EventPlayerDeleted struct {
	ID int64
}
