package android

import (
	"context"
	"errors"
	"sync"

	"github.com/ItsNotGoodName/radiomux/internal/core"
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

func NewController(hooks ControllerHook, bus core.Bus) (*Controller, func()) {
	c := &Controller{
		mu:       sync.Mutex{},
		hook:     hooks,
		players:  make(map[int64]struct{}),
		handlers: make(map[int64]func(ctx context.Context, cmd Command) error),
	}
	return c, bus.OnPlayerDeleted(c.onPlayerDeleted)
}

func (c *Controller) onPlayerDeleted(ctx context.Context, evt core.EventPlayerDeleted) error {
	return c.Handle(ctx, evt.ID, CommandDisconnect{})
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
