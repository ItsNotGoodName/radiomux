package apiws

import (
	"context"
	"errors"

	"github.com/ItsNotGoodName/radiomux/internal/android"
	"github.com/ItsNotGoodName/radiomux/internal/core"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type Server struct {
	stateService       *android.StateService
	playerStore        core.PlayerStore
	notificationServer *NotificationServer
}

func NewServer(stateService *android.StateService, playerStore core.PlayerStore, notificationServer *NotificationServer) Server {
	return Server{
		stateService:       stateService,
		playerStore:        playerStore,
		notificationServer: notificationServer,
	}
}

func (s Server) ServeEcho(c echo.Context) error {
	w := c.Response()
	r := c.Request()
	conn, err := wsUpgrade(w, r)
	if err != nil {
		return nil
	}

	s.handle(r.Context(), conn)

	return nil
}

func (s Server) handle(ctx context.Context, conn *websocket.Conn) {
	log := log.With().
		Str("package", "apiws").
		Str("remote", conn.RemoteAddr().String()).
		Logger()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Subscribe to all player state changes
	stateC, stateCDisconnect := s.stateService.Subscribe()
	defer stateCDisconnect()

	// Subscribe to all notifcations
	notificationC, notificationCancel := s.notificationServer.Subscribe()
	defer notificationCancel()

	// Visitors allow the client to control when they are ready to receive updates
	playerStateVisitor := newPlayerStateVisitor(s.stateService, s.playerStore)
	bufferVisitor := newBufferVisitor(10)
	visitors := newVisitors(playerStateVisitor, bufferVisitor)

	// Write pump
	sync := newSignal()
	flushC := make(chan chan []byte)
	go wsWriter(ctx, cancel, conn, log, sync, flushC)

	// Read pump
	readC := make(chan []byte)
	go wsReader(ctx, cancel, conn, log, readC)

	for {
		select {
		case <-ctx.Done():
			return
		case data, ok := <-readC:
			// Reader
			if !ok {
				return
			}

			log.Error().Bytes("data", data).Msg("The WebSocket client is not supposed to send data...")
		case dataC := <-flushC:
			// Writer
			ok := func() bool {
				data, err := visitors.Visit(ctx)
				if err != nil {
					if errors.Is(err, errVisitorEmpty) {
						return true
					}
					log.Err(err).Msg("Failed to flush")
					return false
				}

				select {
				case <-ctx.Done():
					return false
				case dataC <- data:
				}

				if visitors.HasMore() {
					sync.Queue()
				}

				return true
			}()
			close(dataC)
			if !ok {
				return
			}
		case state, ok := <-stateC:
			// Player state
			if !ok {
				return
			}

			playerStateVisitor.StateChange(state)

			if playerStateVisitor.HasMore() {
				sync.Queue()
			}
		case notification := <-notificationC:
			// Notification
			err := bufferVisitor.Push(notification)
			if err != nil {
				log.Err(err).Msg("Failed to pus to buffer")
				return
			}

			if bufferVisitor.HasMore() {
				sync.Queue()
			}
		}
	}
}
