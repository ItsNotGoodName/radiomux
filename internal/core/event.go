package core

import "context"

// Bus is a synchronous event bus.
type Bus interface {
	PlayerDisconnected(evt EventPlayerDisconnected)
	OnPlayerDisconnected(h func(ctx context.Context, evt EventPlayerDisconnected) error)
	PlayerConnected(evt EventPlayerConnected)
	OnPlayerConnected(h func(ctx context.Context, evt EventPlayerConnected) error)
	PlayerCreated(evt EventPlayerCreated)
	OnPlayerCreated(h func(ctx context.Context, evt EventPlayerCreated) error)
	PlayerTokenUpdated(evt EventPlayerTokenUpdated)
	OnPlayerTokenUpdated(h func(ctx context.Context, evt EventPlayerTokenUpdated) error)
	PlayerDeleted(evt EventPlayerDeleted)
	OnPlayerDeleted(h func(ctx context.Context, evt EventPlayerDeleted) error)
}

type EventPlayerDisconnected struct {
	ID int64
}

type EventPlayerConnected struct {
	ID int64
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
