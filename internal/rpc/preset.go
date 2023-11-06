package rpc

import (
	"context"
	"errors"

	"github.com/ItsNotGoodName/radiomux/internal"
	"github.com/ItsNotGoodName/radiomux/internal/core"
	"github.com/ItsNotGoodName/radiomux/internal/webrpc"
)

type PresetStore interface {
	Create(ctx context.Context, req core.Preset) (core.Preset, error)
	Get(ctx context.Context, id int64) (core.Preset, error)
	List(ctx context.Context) ([]core.Preset, error)
	Update(ctx context.Context, req core.Preset) (core.Preset, error)
	Delete(ctx context.Context, id int64) error
	Drop(ctx context.Context) ([]core.Preset, error)
}

func NewPresetService(presetStore PresetStore) *PresetService {
	return &PresetService{
		presetStore: presetStore,
	}
}

type PresetService struct {
	presetStore PresetStore
}

func (s *PresetService) PresetCreate(ctx context.Context, req *webrpc.CreatePreset) (int64, error) {
	preset, err := s.presetStore.Create(ctx, core.Preset{Name: req.Name})
	if err != nil {
		return 0, webrpc.ConvertErr(err)
	}

	return preset.ID, nil
}

func (s *PresetService) PresetGet(ctx context.Context, id int64) (*webrpc.Preset, error) {
	preset, err := s.presetStore.Get(ctx, id)
	if err != nil {
		return nil, webrpc.ConvertErr(err)
	}

	res, err := webrpc.ConvertPreset(preset)
	if err != nil {
		return nil, webrpc.ConvertErr(err)
	}

	return res, nil
}

func (s *PresetService) PresetList(ctx context.Context) ([]*webrpc.Preset, error) {
	presets, err := s.presetStore.List(ctx)
	if err != nil {
		return nil, webrpc.ConvertErr(err)
	}

	res, err := webrpc.ConvertPresets(presets)
	if err != nil {
		return nil, webrpc.ConvertErr(err)
	}

	return res, nil
}

func (s *PresetService) PresetUpdate(ctx context.Context, req *webrpc.UpdatePreset) error {
	preset, err := s.presetStore.Get(ctx, req.Id)
	if err != nil {
		return webrpc.ConvertErr(err)
	}
	if req.Name != nil {
		preset.Name = *req.Name
	}
	if req.Url != nil {
		return webrpc.ErrNotImplemented
	}
	preset, err = s.presetStore.Update(ctx, preset)
	if err != nil {
		return webrpc.ConvertErr(err)
	}

	return nil
}

func (s *PresetService) PresetDelete(ctx context.Context, ids []int64) error {
	for _, id := range ids {
		err := s.presetStore.Delete(ctx, id)
		if err != nil && !errors.Is(err, internal.ErrNotFound) {
			return webrpc.ConvertErr(err)
		}
	}

	return nil
}

var _ webrpc.PresetService = (*PresetService)(nil)
