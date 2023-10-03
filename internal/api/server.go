package api

import (
	"net/http"

	"github.com/ItsNotGoodName/radiomux/internal/android"
	"github.com/ItsNotGoodName/radiomux/internal/core"
	"github.com/ItsNotGoodName/radiomux/internal/openapi"
	echo "github.com/labstack/echo/v4"
)

func MountServer(e *echo.Echo, s openapi.ServerInterface, path string) {
	g := e.Group("", middlewareError)
	openapi.RegisterHandlersWithBaseURL(g, s, path)
}

func middlewareError(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := next(c); err != nil {
			return openapi.ConvertErr(err)
		}
		return nil
	}
}

func NewServer(playerStore core.PlayerStore, presetSTore core.PresetStore, bus android.BusCommand) *Server {
	return &Server{
		playerStore: playerStore,
		presetStore: presetSTore,
		androidBus:  bus,
	}
}

type Server struct {
	playerStore core.PlayerStore
	presetStore core.PresetStore
	androidBus  android.BusCommand
}

// PostPlayersIdPreset implements ServerInterface.
func (s *Server) PostPlayersIdPreset(c echo.Context, id int64, params openapi.PostPlayersIdPresetParams) error {
	ctx := c.Request().Context()

	presetModel, err := s.presetStore.Get(ctx, params.Preset)
	if err != nil {
		return err
	}

	if err := s.androidBus.Handle(ctx, id, android.CommandSetMediaItem{URI: presetModel.URI()}); err != nil {
		return err
	}

	if err := s.androidBus.Handle(ctx, id, android.CommandPlay{}); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

// DeletePresets implements ServerInterface.
func (s *Server) DeletePresets(c echo.Context) error {
	ctx := c.Request().Context()

	presetsModel, err := s.presetStore.Drop(ctx)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, openapi.ConvertPresets(presetsModel))
}

// DeletePresetsId implements ServerInterface.
func (s *Server) DeletePresetsId(c echo.Context, id int64) error {
	ctx := c.Request().Context()

	err := s.presetStore.Delete(ctx, id)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

// GetPresets implements ServerInterface.
func (s *Server) GetPresets(c echo.Context) error {
	ctx := c.Request().Context()

	presetsModel, err := s.presetStore.List(ctx)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, openapi.ConvertPresets(presetsModel))
}

// GetPresetsId implements ServerInterface.
func (s *Server) GetPresetsId(c echo.Context, id int64) error {
	ctx := c.Request().Context()

	presetModel, err := s.presetStore.Get(ctx, id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, openapi.ConvertPreset(presetModel))
}

// PostPresets implements ServerInterface.
func (s *Server) PostPresets(c echo.Context) error {
	ctx := c.Request().Context()

	body := openapi.PostPresetsJSONBody{}
	if err := c.Bind(&body); err != nil {
		return err
	}

	presetModel, err := s.presetStore.Create(ctx, core.Preset{Name: body.Name, URL: body.Url})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, openapi.ConvertPreset(presetModel))
}

// PostPresetsId implements ServerInterface.
func (s *Server) PostPresetsId(c echo.Context, id int64) error {
	ctx := c.Request().Context()

	body := openapi.PostPresetsIdJSONBody{}
	if err := c.Bind(&body); err != nil {
		return err
	}

	presetModel, err := s.presetStore.Get(ctx, id)
	if err != nil {
		return err
	}
	if body.Name != nil {
		presetModel.Name = *body.Name
	}
	if body.Url != nil {
		presetModel.URL = *body.Url
	}
	presetModel, err = s.presetStore.Update(ctx, presetModel)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, openapi.ConvertPreset(presetModel))
}

// DeletePlayers implements ServerInterface.
func (s *Server) DeletePlayers(c echo.Context) error {
	ctx := c.Request().Context()

	playersModel, err := s.playerStore.Drop(ctx)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, openapi.ConvertPlayers(playersModel))
}

// PostPlayers implements ServerInterface.
func (s *Server) PostPlayers(c echo.Context) error {
	ctx := c.Request().Context()

	body := openapi.PostPlayersJSONBody{}
	if err := c.Bind(&body); err != nil {
		return err
	}

	playerModel, err := s.playerStore.Create(ctx, core.Player{Name: body.Name})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, openapi.ConvertPlayer(playerModel))
}

// PostPlayersId implements ServerInterface.
func (s *Server) PostPlayersId(c echo.Context, id int64) error {
	ctx := c.Request().Context()

	body := openapi.PostPlayersIdJSONBody{}
	if err := c.Bind(&body); err != nil {
		return err
	}

	playerModel, err := s.playerStore.Get(ctx, id)
	if err != nil {
		return err
	}
	if body.Name != nil {
		playerModel.Name = *body.Name
	}
	playerModel, err = s.playerStore.Update(ctx, playerModel)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, openapi.ConvertPlayer(playerModel))
}

// DeletePlayersId implements ServerInterface.
func (s *Server) DeletePlayersId(c echo.Context, id int64) error {
	ctx := c.Request().Context()

	err := s.playerStore.Delete(ctx, id)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

// GetPlayersId implements ServerInterface.
func (s *Server) GetPlayersId(c echo.Context, id int64) error {
	ctx := c.Request().Context()

	playerModel, err := s.playerStore.Get(ctx, id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, openapi.ConvertPlayer(playerModel))
}

// PostPlayersIdSeek implements ServerInterface.
func (s *Server) PostPlayersIdSeek(c echo.Context, id int64) error {
	err := s.androidBus.Handle(c.Request().Context(), id, android.CommandSeekToDefaultPosition{})
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

// PostPlayersIdPause implements ServerInterface.
func (s *Server) PostPlayersIdPause(c echo.Context, id int64) error {
	ctx := c.Request().Context()

	err := s.androidBus.Handle(ctx, id, android.CommandPause{})
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

// PostPlayersIdMedia implements ServerInterface.
func (s *Server) PostPlayersIdMedia(c echo.Context, id int64, params openapi.PostPlayersIdMediaParams) error {
	ctx := c.Request().Context()

	err := s.androidBus.Handle(ctx, id, android.CommandSetMediaItem{
		URI: params.Uri,
	})
	if err != nil {
		return err
	}

	err = s.androidBus.Handle(ctx, id, android.CommandPlay{})
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

// PostPlayersIdVolume implements ServerInterface.
func (s *Server) PostPlayersIdVolume(c echo.Context, id int64, params openapi.PostPlayersIdVolumeParams) error {
	ctx := c.Request().Context()

	if params.Delta != nil {
		if *params.Delta > 0 {
			err := s.androidBus.Handle(ctx, id, android.CommandIncreaseDeviceVolume{})
			if err != nil {
				return err
			}
		} else if *params.Delta < 0 {
			err := s.androidBus.Handle(ctx, id, android.CommandDecreaseDeviceVolume{})
			if err != nil {
				return err
			}
		}
	}

	if params.Mute != nil {
		err := s.androidBus.Handle(ctx, id, android.CommandSetDeviceMuted{Muted: *params.Mute})
		if err != nil {
			return err
		}
	}

	return c.NoContent(http.StatusNoContent)
}

// GetPlayers implements ServerInterface.
func (s *Server) GetPlayers(c echo.Context) error {
	ctx := c.Request().Context()

	playersModel, err := s.playerStore.List(ctx)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, openapi.ConvertPlayers(playersModel))
}

// PostPlayersIdPlay implements ServerInterface.
func (s *Server) PostPlayersIdPlay(c echo.Context, id int64) error {
	ctx := c.Request().Context()

	err := s.androidBus.Handle(ctx, id, android.CommandPlay{})
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

// PostPlayersIdStop implements ServerInterface.
func (s *Server) PostPlayersIdStop(c echo.Context, id int64) error {
	ctx := c.Request().Context()

	err := s.androidBus.Handle(ctx, id, android.CommandStop{})
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

var _ openapi.ServerInterface = (*Server)(nil)
