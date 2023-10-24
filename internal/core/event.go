package core

import "context"

// Bus is a synchronous event bus.
type Bus interface {
	PlayerCreated(ctx context.Context, evt EventPlayerCreated)
	OnPlayerCreated(func(ctx context.Context, evt EventPlayerCreated) error)
	PlayerTokenUpdated(ctx context.Context, evt EventPlayerTokenUpdated)
	OnPlayerTokenUpdated(func(ctx context.Context, evt EventPlayerTokenUpdated) error)
	PlayerDeleted(ctx context.Context, evt EventPlayerDeleted)
	OnPlayerDeleted(func(ctx context.Context, evt EventPlayerDeleted) error)
}

type EventPlayerCreated struct {
	ID int64
}

type EventPlayerTokenUpdated struct {
	ID int64
}

type EventPlayerDeleted struct {
	ID int64
}
