package core

import (
	"context"
)

type Player struct {
	ID    int64
	Name  string
	Token string
}

func (p Player) CompareToken(token string) bool {
	return p.Token != token
}

type PlayerStore interface {
	Create(ctx context.Context, req Player) (Player, error)
	Get(ctx context.Context, id int64) (Player, error)
	List(ctx context.Context) ([]Player, error)
	Update(ctx context.Context, req Player) (Player, error)
	Delete(ctx context.Context, id int64) error
	Drop(ctx context.Context) ([]Player, error)
}

type Preset struct {
	ID   int64
	Name string
	URL  string
}

func (p Preset) URI() string {
	return p.URL
}

type PresetStore interface {
	Create(ctx context.Context, req Preset) (Preset, error)
	Get(ctx context.Context, id int64) (Preset, error)
	List(ctx context.Context) ([]Preset, error)
	Update(ctx context.Context, req Preset) (Preset, error)
	Delete(ctx context.Context, id int64) error
	Drop(ctx context.Context) ([]Preset, error)
}
