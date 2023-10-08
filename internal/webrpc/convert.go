package webrpc

import (
	"errors"

	"github.com/ItsNotGoodName/radiomux/internal"
	"github.com/ItsNotGoodName/radiomux/internal/android"
	"github.com/ItsNotGoodName/radiomux/internal/core"
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

func ConvertPlayers(req []core.Player) []*Player {
	res := make([]*Player, 0, len(req))
	for _, r := range req {
		res = append(res, ConvertPlayer(r))
	}
	return res
}

func ConvertPlayer(p core.Player) *Player {
	return &Player{
		Id:    p.ID,
		Name:  p.Name,
		Token: p.Token,
	}
}

func ConvertPresets(req []core.Preset) []*Preset {
	res := make([]*Preset, 0, len(req))
	for _, r := range req {
		res = append(res, ConvertPreset(r))
	}
	return res
}

func ConvertPreset(p core.Preset) *Preset {
	return &Preset{
		Id:   p.ID,
		Name: p.Name,
		Url:  p.URL,
	}
}
