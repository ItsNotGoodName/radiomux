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

func New() *Bus {
	return &Bus{}
}

type Bus struct {
	mu                   sync.Mutex
	onPlayerTokenChanged []func(ctx context.Context, evt core.EventPlayerTokenUpdated) error
	onPlayerCreated      []func(ctx context.Context, evt core.EventPlayerCreated) error
	onPlayerDeleted      []func(ctx context.Context, evt core.EventPlayerDeleted) error
}

// OnPlayerTokenUpdated implements core.Bus.
func (b *Bus) OnPlayerTokenUpdated(h func(ctx context.Context, evt core.EventPlayerTokenUpdated) error) {
	b.mu.Lock()
	b.onPlayerTokenChanged = append(b.onPlayerTokenChanged, h)
	b.mu.Unlock()
}

// PlayerTokenUpdated implements core.Bus.
func (b *Bus) PlayerTokenUpdated(ctx context.Context, evt core.EventPlayerTokenUpdated) {
	b.mu.Lock()
	defer b.mu.Unlock()

	for _, v := range b.onPlayerTokenChanged {
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
func (b *Bus) PlayerCreated(ctx context.Context, evt core.EventPlayerCreated) {
	b.mu.Lock()
	defer b.mu.Unlock()

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
func (b *Bus) PlayerDeleted(ctx context.Context, evt core.EventPlayerDeleted) {
	b.mu.Lock()
	defer b.mu.Unlock()

	for _, v := range b.onPlayerDeleted {
		logErr(v(ctx, evt))
	}
}

var _ core.Bus = (*Bus)(nil)
