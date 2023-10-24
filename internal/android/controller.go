package android

import (
	"context"
	"errors"
	"sync"

	"github.com/ItsNotGoodName/radiomux/internal/core"
	"github.com/rs/zerolog/log"
)

var (
	ErrPlayerNotConnected error = errors.New("player not connected")
	ErrPlayerConflict     error = errors.New("player conflict")
	ErrPlayerNotReady     error = errors.New("player not ready")
)

// ControllerHook must not have a reference to the Controller to prevent deadlocks.
type ControllerHook interface {
	PlayerConnectChanged(id int64) error
	PlayerDisconnectChanged(id int64) error
	PlayerReadyChanged(id int64) error
}

// Controller handles player connections.
type Controller struct {
	mu       sync.Mutex
	hook     ControllerHook
	players  map[int64]struct{}
	handlers map[int64]func(ctx context.Context, cmd Command) error
}

func NewController(hooks ControllerHook, bus core.Bus) *Controller {
	c := &Controller{
		mu:       sync.Mutex{},
		hook:     hooks,
		players:  make(map[int64]struct{}),
		handlers: make(map[int64]func(ctx context.Context, cmd Command) error),
	}

	bus.OnPlayerTokenUpdated(c.onPlayerTokenUpdated)
	bus.OnPlayerDeleted(c.onPlayerDeleted)

	return c
}

func (c *Controller) onPlayerTokenUpdated(ctx context.Context, evt core.EventPlayerTokenUpdated) error {
	return c.onPlayerDeleted(ctx, core.EventPlayerDeleted{ID: evt.ID})
}

func (c *Controller) onPlayerDeleted(ctx context.Context, evt core.EventPlayerDeleted) error {
	err := c.Handle(ctx, evt.ID, CommandDisconnect{})
	if !errors.Is(err, ErrPlayerNotConnected) {
		return err
	}
	return nil
}

func (s *Controller) Handle(ctx context.Context, id int64, cmd Command) error {
	s.mu.Lock()
	handler, ok := s.handlers[id]
	if !ok {
		_, ok := s.players[id]
		s.mu.Unlock()
		if ok {
			return ErrPlayerNotReady
		}
		return ErrPlayerNotConnected
	}
	s.mu.Unlock()

	return handler(ctx, cmd)
}

func (s *Controller) PlayerConnect(id int64) (func(), error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	log.Debug().Str("package", "android").Int64("id", id).Msg("Connect to controller")
	_, ok := s.players[id]
	if ok {
		return nil, ErrPlayerConflict
	}

	err := s.hook.PlayerConnectChanged(id)
	if err != nil {
		return nil, err
	}

	s.players[id] = struct{}{}
	return func() {
		s.mu.Lock()
		defer s.mu.Unlock()

		log.Debug().Str("package", "android").Int64("id", id).Msg("Disconnect from controller")
		if _, found := s.players[id]; !found {
			return
		}
		delete(s.players, id)
		delete(s.handlers, id)

		_ = s.hook.PlayerDisconnectChanged(id)
	}, nil
}

func (s *Controller) PlayerReady(id int64, handler func(ctx context.Context, cmd Command) error) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, found := s.handlers[id]; found {
		return errors.New("duplicate player ready handler")
	}
	err := s.hook.PlayerReadyChanged(id)
	if err != nil {
		return err
	}

	s.handlers[id] = handler
	return nil
}

var _ BusCommand = (*Controller)(nil)
