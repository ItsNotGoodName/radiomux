package bus

import (
	"context"

	"github.com/ItsNotGoodName/radiomux/internal/core"
	"github.com/google/uuid"
	"github.com/mustafaturan/bus/v3"
	"github.com/rs/zerolog/log"
)

func logEmitErr(err error) {
	if err != nil {
		log.Err(err).Msg("Failed to emit bus event")
	}
}

type generator struct{}

func (generator) Generate() string {
	return uuid.NewString()
}

func New() (Bus, error) {
	bus, err := bus.NewBus(generator{})
	if err != nil {
		return Bus{}, err
	}

	bus.RegisterTopics(
		TopicPlayerDeleted,
	)

	return Bus{
		bus: bus,
	}, nil
}

const (
	TopicPlayerDeleted = "player.deleted"
)

type Bus struct {
	bus *bus.Bus
}

// OnPlayerDeleted implements models.Bus.
func (b Bus) OnPlayerDeleted(h func(ctx context.Context, evt core.EventPlayerDeleted) error) func() {
	key := uuid.NewString()

	b.bus.RegisterHandler(key, bus.Handler{
		Handle: func(ctx context.Context, evt bus.Event) {
			h(ctx, evt.Data.(core.EventPlayerDeleted))
		},
		Matcher: TopicPlayerDeleted,
	})

	return func() { b.bus.DeregisterHandler(key) }
}

// PlayerDeleted implements models.Bus.
func (b Bus) PlayerDeleted(ctx context.Context, evt core.EventPlayerDeleted) {
	logEmitErr(b.bus.Emit(ctx, TopicPlayerDeleted, evt))
}

var _ core.Bus = Bus{}
