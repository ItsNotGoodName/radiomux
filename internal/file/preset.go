package file

import (
	"context"

	"github.com/ItsNotGoodName/radiomux/internal"
	"github.com/ItsNotGoodName/radiomux/internal/core"
	"github.com/samber/lo"
)

func NewPresetStore(store *Store, bus core.Bus) PresetStore {
	return PresetStore{
		store: store,
		bus:   bus,
	}
}

type PresetStore struct {
	store *Store
	bus   core.Bus
}

// Drop implements models.PresetStore.
func (s PresetStore) Drop(ctx context.Context) ([]core.Preset, error) {
	err := s.store.Update(func(db *db) error {
		db.Presets = []presetModel{}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return s.List(ctx)
}

// List implements models.PresetStore.
func (s PresetStore) List(ctx context.Context) ([]core.Preset, error) {
	db, err := s.store.Read()
	if err != nil {
		return nil, err
	}

	return lo.Map(db.Presets, func(f presetModel, _ int) core.Preset { return convertPreset(f) }), nil
}

// Delete implements models.PresetStore.
func (s PresetStore) Delete(ctx context.Context, id int64) error {
	return s.store.Update(func(db *db) error {
		found := false
		db.Presets = lo.Filter(db.Presets, func(item presetModel, index int) bool {
			if item.ID == id {
				found = true
				return false
			}
			return true
		})
		if !found {
			return internal.ErrNotFound
		}
		return nil
	})
}

// Update implements models.PresetStore.
func (s PresetStore) Update(ctx context.Context, preset core.Preset) (core.Preset, error) {
	err := s.store.Update(func(db *db) error {
		for i := range db.Presets {
			if db.Presets[i].ID == preset.ID {
				db.Presets[i] = unconvertPreset(preset)
				return nil
			}
		}

		return internal.ErrNotFound
	})
	if err != nil {
		return core.Preset{}, err
	}

	return preset, nil
}

// Create implements models.PresetStore.
func (s PresetStore) Create(ctx context.Context, preset core.Preset) (core.Preset, error) {
	err := s.store.Update(func(db *db) error {
		last, err := lo.Last(db.Presets)
		if err != nil {
			preset.ID = 1
		} else {
			preset.ID = last.ID + 1
		}

		db.Presets = append(db.Presets, unconvertPreset(preset))

		return nil
	})
	if err != nil {
		return core.Preset{}, err
	}

	return preset, nil
}

func (s PresetStore) Get(ctx context.Context, id int64) (core.Preset, error) {
	db, err := s.store.Read()
	if err != nil {
		return core.Preset{}, err
	}

	preset, ok := lo.Find(db.Presets, func(preset presetModel) bool { return preset.ID == id })
	if !ok {
		return core.Preset{}, internal.ErrNotFound
	}

	return convertPreset(preset), nil
}

var _ core.PresetStore = PresetStore{}
