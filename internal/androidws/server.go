package androidws

import (
	"errors"
	"fmt"

	"github.com/ItsNotGoodName/radiomux/internal"
	"github.com/ItsNotGoodName/radiomux/internal/android"
	"github.com/ItsNotGoodName/radiomux/internal/core"
	"github.com/labstack/echo/v4"
)

type Server struct {
	playerStore core.PlayerStore
	controller  *android.Controller
	busEvent    android.BusEvent
	wsURL       string
}

func NewServer(playerStore core.PlayerStore, controller *android.Controller, busEvent android.BusEvent, wsURL string) Server {
	return Server{
		playerStore: playerStore,
		controller:  controller,
		busEvent:    busEvent,
		wsURL:       wsURL,
	}
}

const Path = "/api/android/ws"

func (s Server) PlayerWSURL(p core.Player) string {
	return fmt.Sprintf("%s%s?id=%d&token=%s", s.wsURL, Path, p.ID, p.Token)
}

func (s Server) ServeEcho(c echo.Context) error {
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
	conn, err := wsUpgrade(w, r)
	if err != nil {
		return nil
	}

	newConnection().handle(ctx, id, conn, s.controller, s.busEvent)

	return nil
}
