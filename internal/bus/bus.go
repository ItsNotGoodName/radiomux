package bus

import (
	"context"
	"sync"

	"github.com/ItsNotGoodName/radiomux/internal/core"
	"github.com/rs/zerolog/log"
)

func logErr(err error) {
	if err != nil {
		log.Err(err).Msg("Failed to handle bus event")
	}
}

func (*Bus) context() context.Context {
	return context.TODO()
}

func New() *Bus {
	return &Bus{}
}

type Bus struct {
	mu                   sync.Mutex
	onPlayerDisconnected []func(ctx context.Context, evt core.EventPlayerDisconnected) error
	onPlayerConnected    []func(ctx context.Context, evt core.EventPlayerConnected) error
	onPlayerCreated      []func(ctx context.Context, evt core.EventPlayerCreated) error
	onPlayerTokenUpdated []func(ctx context.Context, evt core.EventPlayerTokenUpdated) error
	onPlayerDeleted      []func(ctx context.Context, evt core.EventPlayerDeleted) error
}

// OnPlayerDisconnected implements core.Bus.
func (b *Bus) OnPlayerDisconnected(h func(ctx context.Context, evt core.EventPlayerDisconnected) error) {
	b.mu.Lock()
	b.onPlayerDisconnected = append(b.onPlayerDisconnected, h)
	b.mu.Unlock()
}

// PlayerDisconnected implements core.Bus.
func (b *Bus) PlayerDisconnected(evt core.EventPlayerDisconnected) {
	b.mu.Lock()
	defer b.mu.Unlock()

	ctx := b.context()
	for _, v := range b.onPlayerDisconnected {
		logErr(v(ctx, evt))
	}
}

// OnPlayerConnected implements core.Bus.
func (b *Bus) OnPlayerConnected(h func(ctx context.Context, evt core.EventPlayerConnected) error) {
	b.mu.Lock()
	b.onPlayerConnected = append(b.onPlayerConnected, h)
	b.mu.Unlock()
}

// PlayerConnected implements core.Bus.
func (b *Bus) PlayerConnected(evt core.EventPlayerConnected) {
	b.mu.Lock()
	defer b.mu.Unlock()

	ctx := b.context()
	for _, v := range b.onPlayerConnected {
		logErr(v(ctx, evt))
	}
}

// OnPlayerTokenUpdated implements core.Bus.
func (b *Bus) OnPlayerTokenUpdated(h func(ctx context.Context, evt core.EventPlayerTokenUpdated) error) {
	b.mu.Lock()
	b.onPlayerTokenUpdated = append(b.onPlayerTokenUpdated, h)
	b.mu.Unlock()
}

// PlayerTokenUpdated implements core.Bus.
func (b *Bus) PlayerTokenUpdated(evt core.EventPlayerTokenUpdated) {
	b.mu.Lock()
	defer b.mu.Unlock()

	ctx := b.context()
	for _, v := range b.onPlayerTokenUpdated {
		logErr(v(ctx, evt))
	}
}

// OnPlayerCreated implements core.Bus.
func (b *Bus) OnPlayerCreated(h func(ctx context.Context, evt core.EventPlayerCreated) error) {
	b.mu.Lock()
	b.onPlayerCreated = append(b.onPlayerCreated, h)
	b.mu.Unlock()
}

// PlayerCreated implements core.Bus.
func (b *Bus) PlayerCreated(evt core.EventPlayerCreated) {
	b.mu.Lock()
	defer b.mu.Unlock()

	ctx := b.context()
	for _, v := range b.onPlayerCreated {
		logErr(v(ctx, evt))
	}
}

// OnPlayerDeleted implements models.Bus.
func (b *Bus) OnPlayerDeleted(h func(ctx context.Context, evt core.EventPlayerDeleted) error) {
	b.mu.Lock()
	b.onPlayerDeleted = append(b.onPlayerDeleted, h)
	b.mu.Unlock()
}

// PlayerDeleted implements models.Bus.
func (b *Bus) PlayerDeleted(evt core.EventPlayerDeleted) {
	b.mu.Lock()
	defer b.mu.Unlock()

	ctx := b.context()
	for _, v := range b.onPlayerDeleted {
		logErr(v(ctx, evt))
	}
}

var _ core.Bus = (*Bus)(nil)
