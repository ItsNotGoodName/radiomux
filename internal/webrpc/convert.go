package webrpc

import (
	"errors"

	"github.com/ItsNotGoodName/radiomux/internal"
	"github.com/ItsNotGoodName/radiomux/internal/android"
	"github.com/ItsNotGoodName/radiomux/internal/core"
	"github.com/samber/lo"
)

func ConvertErr(err error) error {
	if errors.Is(err, android.ErrPlayerNotConnected) {
		return ErrNotFound.WithCause(err)
	}
	if errors.Is(err, internal.ErrNotFound) {
		return ErrNotFound.WithCause(err)
	}
	if errors.Is(err, android.ErrPlayerNotReady) {
		return ErrWebrpcInternalError.WithCause(err)
	}
	return ErrWebrpcInternalError.WithCause(err)
}

func ConvertPlayers(p []core.Player) []*Player {
	return lo.Map(p, func(p core.Player, _ int) *Player { return ConvertPlayer(p) })
}

func ConvertPlayer(p core.Player) *Player {
	return &Player{
		Id:    p.ID,
		Name:  p.Name,
		Token: p.Token,
	}
}

func ConvertPresets(p []core.Preset) []*Preset {
	return lo.Map(p, func(p core.Preset, _ int) *Preset { return ConvertPreset(p) })
}

func ConvertPreset(p core.Preset) *Preset {
	return &Preset{
		Id:   p.ID,
		Name: p.Name,
		Url:  p.URL,
	}
}
