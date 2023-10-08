package rpc

import (
	"context"

	"github.com/ItsNotGoodName/radiomux/internal/android"
	"github.com/ItsNotGoodName/radiomux/internal/core"
	"github.com/ItsNotGoodName/radiomux/internal/webrpc"
)

func NewStateService(androidBus android.BusCommand, presetStore core.PresetStore) *StateService {
	return &StateService{
		androidBus:  androidBus,
		presetStore: presetStore,
	}
}

type StateService struct {
	androidBus  android.BusCommand
	presetStore core.PresetStore
}

func (s *StateService) StateMediaSet(ctx context.Context, req *webrpc.SetStateMedia) error {
	if req.PresetId != nil {
		presetID := *req.PresetId

		presetModel, err := s.presetStore.Get(ctx, presetID)
		if err != nil {
			return webrpc.ConvertErr(err)
		}

		if err := s.androidBus.Handle(ctx, req.Id, android.CommandSetMediaItem{URI: presetModel.URI()}); err != nil {
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
