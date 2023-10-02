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
	stateService *android.StateService
	playerStore  core.PlayerStore
}

func NewServer(stateService *android.StateService, playerStore core.PlayerStore) Server {
	return Server{
		stateService: stateService,
		playerStore:  playerStore,
	}
}

func (s Server) Handle(c echo.Context) error {
	w := c.Response()
	r := c.Request()
	conn, err := wsUpgrade(w, r)
	if err != nil {
		return err
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

	// We subscribe to all player state changes
	stateC, stateCDisconnect := s.stateService.Subscribe()
	defer stateCDisconnect()

	// This allows the client to control when they are ready to receive updates
	playerStateVisitor := newPlayerStateVisitor(s.stateService, s.playerStore)
	visitors := newVisitors(playerStateVisitor)
	sync := newSignal()

	// Write pump
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
				data, err := visitors.Visit()
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
		}
	}
}
