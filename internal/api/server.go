package api

import (
	"context"
	"net/http"

	"github.com/ItsNotGoodName/radiomux/internal/core"
	"github.com/cyphar/filepath-securejoin"
	echo "github.com/labstack/echo/v4"
	qrcode "github.com/skip2/go-qrcode"
)

type SourceStore interface {
	GetFileSource(ctx context.Context, id int64) (core.FileSource, error)
}

func NewServer(playerStore core.PlayerStore, sourceStore SourceStore) *Server {
	return &Server{
		playerStore: playerStore,
		sourceStore: sourceStore,
	}
}

type Server struct {
	playerStore core.PlayerStore
	sourceStore SourceStore
}

func (s *Server) GetSourcesIdSlug(c echo.Context, id int64, slug string) error {
	ctx := c.Request().Context()

	source, err := s.sourceStore.GetFileSource(ctx, id)
	if err != nil {
		return err
	}

	securePath, err := securejoin.SecureJoin(source.Path, slug)
	if err != nil {
		return err
	}

	return c.File(securePath)
}

func (s *Server) GetPlayersIdQr(c echo.Context, id int64) error {
	ctx := c.Request().Context()

	player, err := s.playerStore.Get(ctx, id)
	if err != nil {
		return err
	}

	var png []byte
	url := core.Settings.PlayerWSURL(player)
	png, err = qrcode.Encode(url, qrcode.Medium, 256)
	if err != nil {
		return err
	}

	c.Response().Header().Set("Cache-Control", "no-store")
	return c.Blob(http.StatusOK, "image/png", png)
}
