package api

import (
	"fmt"
	"net/http"

	"github.com/ItsNotGoodName/radiomux/internal/core"
	"github.com/ItsNotGoodName/radiomux/internal/openapi"
	echo "github.com/labstack/echo/v4"
	qrcode "github.com/skip2/go-qrcode"
)

func MountServer(e *echo.Echo, s openapi.ServerInterface) {
	g := e.Group("/api", middlewareError)
	openapi.RegisterHandlers(g, s)
}

func middlewareError(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := next(c); err != nil {
			return openapi.ConvertErr(err)
		}
		return nil
	}
}

func NewServer(playerStore core.PlayerStore) *Server {
	return &Server{
		playerStore: playerStore,
	}
}

type Server struct {
	playerStore core.PlayerStore
}

func (s *Server) GetPlayersIdQr(c echo.Context, id int64) error {
	ctx := c.Request().Context()

	player, err := s.playerStore.Get(ctx, id)
	if err != nil {
		return err
	}

	var png []byte
	// TODO: use real url
	png, err = qrcode.Encode(fmt.Sprintf("https://example.com/ws?id=%d&token=%s", player.ID, player.Token), qrcode.Medium, 256)
	if err != nil {
		return err
	}

	c.Response().Header().Set("Cache-Control", "no-store")
	return c.Blob(http.StatusOK, "image/png", png)
}

var _ openapi.ServerInterface = (*Server)(nil)
