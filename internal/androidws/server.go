package androidws

import (
	"errors"

	"github.com/ItsNotGoodName/radiomux/internal"
	"github.com/ItsNotGoodName/radiomux/internal/android"
	"github.com/ItsNotGoodName/radiomux/internal/core"
	"github.com/labstack/echo/v4"
	"nhooyr.io/websocket"
)

type Server struct {
	playerStore core.PlayerStore
	controller  *android.Controller
	busEvent    android.BusEvent
}

func NewServer(playerStore core.PlayerStore, controller *android.Controller, busEvent android.BusEvent) Server {
	return Server{
		playerStore: playerStore,
		controller:  controller,
		busEvent:    busEvent,
	}
}

func (s Server) Handle(c echo.Context) error {
	ctx := c.Request().Context()

	// Auth
	id, err := auth(ctx, c, s.playerStore)
	if err != nil {
		return err
	}

	// Connect
	disconnect, err := s.controller.PlayerConnect(id)
	if err != nil {
		if errors.Is(err, internal.ErrNotFound) {
			return echo.ErrNotFound.WithInternal(err)
		}
		return echo.ErrConflict.WithInternal(err)
	}
	defer disconnect()

	// Setup websocket
	w := c.Response()
	r := c.Request()
	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		return err
	}
	defer conn.Close(websocket.StatusInternalError, "the sky is falling")

	newConnection(conn).handle(ctx, s.controller, s.busEvent, id)

	return nil
}
