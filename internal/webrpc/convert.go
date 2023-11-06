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
	if errors.Is(err, errors.ErrUnsupported) {
		return ErrNotImplemented.WithCause(err)
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

func ConvertPresets(req []core.Preset) ([]*Preset, error) {
	res := make([]*Preset, 0, len(req))
	for _, r := range req {
		p, err := ConvertPreset(r)
		if err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, nil
}

func ConvertPreset(p core.Preset) (*Preset, error) {
	var url string
	{
		slug, err := core.PresetSlugParse(p.Slug)
		if err != nil {
			return nil, err
		}
		switch slug := slug.(type) {
		case core.PresetFile:
			url = core.Settings.FileURL(slug.SourceID, slug.Path)
		case core.PresetSubsonic:
		case core.PresetURL:
			url = string(slug)
		default:
		}
	}

	return &Preset{
		Id:   p.ID,
		Name: p.Name,
		Url:  url,
	}, nil
}
