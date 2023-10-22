package demo

import (
	"context"
	"errors"

	"github.com/ItsNotGoodName/radiomux/internal"
	"github.com/ItsNotGoodName/radiomux/internal/core"
)

var MockPlayers []core.Player = []core.Player{
	{
		ID:   1,
		Name: "Mock Player 1",
	},
	{
		ID:   2,
		Name: "Mock Player 2",
	},
	{
		ID:   3,
		Name: "Mock Player 3",
	},
}

func NewPlayerStore() PlayerStore {
	return PlayerStore{}
}

type PlayerStore struct {
}

// Create implements core.PlayerStore.
func (PlayerStore) Create(ctx context.Context, req core.Player) (core.Player, error) {
	return core.Player{}, errors.ErrUnsupported
}

// Delete implements core.PlayerStore.
func (PlayerStore) Delete(ctx context.Context, id int64) error {
	return errors.ErrUnsupported
}

// Drop implements core.PlayerStore.
func (PlayerStore) Drop(ctx context.Context) ([]core.Player, error) {
	return nil, errors.ErrUnsupported
}

// Get implements core.PlayerStore.
func (PlayerStore) Get(ctx context.Context, id int64) (core.Player, error) {
	for _, p := range MockPlayers {
		if p.ID == id {
			return p, nil
		}
	}
	return core.Player{}, internal.ErrNotFound
}

// List implements core.PlayerStore.
func (PlayerStore) List(ctx context.Context) ([]core.Player, error) {
	return MockPlayers, nil
}

// Update implements core.PlayerStore.
func (PlayerStore) Update(ctx context.Context, req core.Player) (core.Player, error) {
	return core.Player{}, errors.ErrUnsupported
}

var _ core.PlayerStore = PlayerStore{}

var MockPresets []core.Preset = []core.Preset{
	{
		ID:   1,
		Name: "Mock Preset 1",
		URL:  "https://example.com/mock-preset-1",
	},
	{
		ID:   2,
		Name: "Mock Preset 2",
		URL:  "https://example.com/mock-preset-2",
	},
	{
		ID:   3,
		Name: "Mock Preset 3",
		URL:  "https://example.com/mock-preset-3",
	},
	{
		ID:   4,
		Name: "Mock Preset 3 Duplicate",
		URL:  "https://example.com/mock-preset-3",
	},
}

func NewPresetStore() PresetStore {
	return PresetStore{}
}

type PresetStore struct {
}

// Create implements core.PresetStore.
func (PresetStore) Create(ctx context.Context, req core.Preset) (core.Preset, error) {
	return core.Preset{}, errors.ErrUnsupported
}

// Delete implements core.PresetStore.
func (PresetStore) Delete(ctx context.Context, id int64) error {
	return errors.ErrUnsupported
}

// Drop implements core.PresetStore.
func (PresetStore) Drop(ctx context.Context) ([]core.Preset, error) {
	return nil, errors.ErrUnsupported
}

// Get implements core.PresetStore.
func (PresetStore) Get(ctx context.Context, id int64) (core.Preset, error) {
	for _, p := range MockPresets {
		if p.ID == id {
			return p, nil
		}
	}
	return core.Preset{}, internal.ErrNotFound
}

// List implements core.PresetStore.
func (PresetStore) List(ctx context.Context) ([]core.Preset, error) {
	return MockPresets, nil
}

// Update implements core.PresetStore.
func (PresetStore) Update(ctx context.Context, req core.Preset) (core.Preset, error) {
	return core.Preset{}, errors.ErrUnsupported
}

var _ core.PresetStore = PresetStore{}
