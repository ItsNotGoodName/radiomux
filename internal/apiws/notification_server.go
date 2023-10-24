package apiws

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/ItsNotGoodName/radiomux/internal"
	"github.com/ItsNotGoodName/radiomux/internal/core"
	"github.com/ItsNotGoodName/radiomux/internal/openapi"
)

type NotificationServer struct {
	mu          sync.Mutex
	subs        map[*chan []byte]struct{}
	playerStore core.PlayerStore
}

func NewNotificationServer(bus core.Bus, playerStore core.PlayerStore) *NotificationServer {
	ns := NotificationServer{
		subs:        make(map[*chan []byte]struct{}),
		playerStore: playerStore,
	}

	// Notications to client
	bus.OnPlayerConnected(ns.onPlayerConnected)
	bus.OnPlayerDisconnected(ns.onPlayerDisconnected)

	return &ns
}

func (n *NotificationServer) onPlayerConnected(ctx context.Context, event core.EventPlayerConnected) error {
	player, err := n.playerStore.Get(ctx, event.ID)
	if err != nil {
		if errors.Is(err, internal.ErrNotFound) {
			return nil
		}
		return err
	}

	b, err := openapi.ConvertNotification(openapi.Notification{
		Title:       fmt.Sprintf("Player %d", event.ID),
		Description: fmt.Sprintf("'%s' connected.", player.Name),
	})
	if err != nil {
		return err
	}

	n.Publish(b)

	return nil
}

func (n *NotificationServer) onPlayerDisconnected(ctx context.Context, event core.EventPlayerDisconnected) error {
	player, err := n.playerStore.Get(ctx, event.ID)
	if err != nil {
		if errors.Is(err, internal.ErrNotFound) {
			return nil
		}
		return err
	}

	b, err := openapi.ConvertNotification(openapi.Notification{
		Title:       fmt.Sprintf("Player %d", event.ID),
		Description: fmt.Sprintf("'%s' disconnected.", player.Name),
		Error:       true,
	})
	if err != nil {
		return err
	}

	n.Publish(b)

	return nil
}

func (n *NotificationServer) Publish(b []byte) {
	n.mu.Lock()
	for sub := range n.subs {
		select {
		case *sub <- b:
		}
	}
	n.mu.Unlock()
}

func (n *NotificationServer) Subscribe() (chan []byte, func()) {
	n.mu.Lock()
	c := make(chan []byte)
	n.subs[&c] = struct{}{}
	n.mu.Unlock()

	return c, func() {
		n.mu.Lock()
		delete(n.subs, &c)
		n.mu.Unlock()
	}
}
