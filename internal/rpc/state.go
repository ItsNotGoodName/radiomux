package rpc

import (
	"context"
	"errors"
	"fmt"

	"github.com/ItsNotGoodName/radiomux/internal/android"
	"github.com/ItsNotGoodName/radiomux/internal/core"
	"github.com/ItsNotGoodName/radiomux/internal/webrpc"
)

func NewStateService(androidBus android.BusCommand, presetStore PresetStore) *StateService {
	return &StateService{
		androidBus:  androidBus,
		presetStore: presetStore,
	}
}

type StateService struct {
	androidBus  android.BusCommand
	presetStore PresetStore
}

func (s *StateService) StateMediaSet(ctx context.Context, req *webrpc.SetStateMedia) error {
	if req.PresetId != nil {
		presetID := *req.PresetId

		presetModel, err := s.presetStore.Get(ctx, presetID)
		if err != nil {
			return webrpc.ConvertErr(err)
		}

		slug, err := core.PresetSlugParse(presetModel.Slug)
		if err != nil {
			return webrpc.ConvertErr(err)
		}

		var uri string
		switch slug := slug.(type) {
		case core.PresetFile:
			uri = core.Settings.FileURL(slug.SourceID, slug.Path)
		case core.PresetSubsonic:
			return webrpc.ErrNotImplemented.WithCause(errors.New("subsonic preset not implemented"))
		case core.PresetURL:
			uri = string(slug)
		default:
			return webrpc.ErrWebrpcInternalError.WithCause(fmt.Errorf("invalid preset type: %T", slug))
		}

		if err := s.androidBus.Handle(ctx, req.Id, android.CommandSetMediaItem{URI: uri}); err != nil {
			return webrpc.ConvertErr(err)
		}
	} else if req.Uri != nil {
		uri := *req.Uri

		err := s.androidBus.Handle(ctx, req.Id, android.CommandSetMediaItem{
			URI: uri,
		})
		if err != nil {
			return webrpc.ConvertErr(err)
		}
	}

	err := s.androidBus.Handle(ctx, req.Id, android.CommandPlay{})
	if err != nil {
		return webrpc.ConvertErr(err)
	}

	return nil
}

func (s *StateService) StateActionSet(ctx context.Context, req *webrpc.SetStateAction) error {
	if req.Action == nil {
		return webrpc.ErrWebrpcBadRequest
	}

	switch *req.Action {
	case webrpc.StateAction_PLAY:
		err := s.androidBus.Handle(ctx, req.Id, android.CommandPlay{})
		if err != nil {
			return webrpc.ConvertErr(err)
		}
	case webrpc.StateAction_PUASE:
		err := s.androidBus.Handle(ctx, req.Id, android.CommandPause{})
		if err != nil {
			return webrpc.ConvertErr(err)
		}
	case webrpc.StateAction_STOP:
		err := s.androidBus.Handle(ctx, req.Id, android.CommandStop{})
		if err != nil {
			return webrpc.ConvertErr(err)
		}
	case webrpc.StateAction_SEEK:
		err := s.androidBus.Handle(ctx, req.Id, android.CommandSeekToDefaultPosition{})
		if err != nil {
			return webrpc.ConvertErr(err)
		}
	default:
		return webrpc.ErrWebrpcBadRequest
	}

	return nil
}

func (s *StateService) StateVolumeSet(ctx context.Context, req *webrpc.SetStateVolume) error {
	if req.Delta != nil {
		delta := *req.Delta
		if delta > 0 {
			err := s.androidBus.Handle(ctx, req.Id, android.CommandIncreaseDeviceVolume{})
			if err != nil {
				return webrpc.ConvertErr(err)
			}
		} else if delta < 0 {
			err := s.androidBus.Handle(ctx, req.Id, android.CommandDecreaseDeviceVolume{})
			if err != nil {
				return webrpc.ConvertErr(err)
			}
		}
	}
	if req.Volume != nil {
		volume := *req.Volume
		err := s.androidBus.Handle(ctx, req.Id, android.CommandSetDeviceVolume{
			Volume: volume,
		})
		if err != nil {
			return webrpc.ConvertErr(err)
		}
	}
	if req.Mute != nil {
		mute := *req.Mute
		err := s.androidBus.Handle(ctx, req.Id, android.CommandSetDeviceMuted{
			Muted: mute,
		})
		if err != nil {
			return webrpc.ConvertErr(err)
		}
	}

	return nil
}

var _ webrpc.StateService = (*StateService)(nil)
