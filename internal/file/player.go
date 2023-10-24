package file

import (
	"context"

	"github.com/ItsNotGoodName/radiomux/internal"
	"github.com/ItsNotGoodName/radiomux/internal/core"
	"github.com/samber/lo"
)

func NewPlayerStore(store *Store, bus core.Bus) PlayerStore {
	return PlayerStore{
		store: store,
		bus:   bus,
	}
}

type PlayerStore struct {
	store *Store
	bus   core.Bus
}

// Drop implements models.PlayerStore.
func (s PlayerStore) Drop(ctx context.Context) ([]core.Player, error) {
	var oldPlayers []playerModel
	err := s.store.Update(func(db *db) error {
		oldPlayers = db.Players
		db.Players = []playerModel{}
		return nil
	})
	if err != nil {
		return nil, err
	}

	for _, p := range oldPlayers {
		s.bus.PlayerDeleted(ctx, core.EventPlayerDeleted{ID: p.ID})
	}

	return s.List(ctx)
}

// List implements models.PlayerStore.
func (s PlayerStore) List(ctx context.Context) ([]core.Player, error) {
	db, err := s.store.Read()
	if err != nil {
		return nil, err
	}

	return lo.Map(db.Players, func(f playerModel, _ int) core.Player { return convertPlayer(f) }), nil
}

// Delete implements models.PlayerStore.
func (s PlayerStore) Delete(ctx context.Context, id int64) error {
	err := s.store.Update(func(db *db) error {
		found := false
		db.Players = lo.Filter(db.Players, func(item playerModel, index int) bool {
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
	if err != nil {
		return err
	}

	s.bus.PlayerDeleted(ctx, core.EventPlayerDeleted{ID: id})

	return nil
}

// Update implements models.PlayerStore.
func (s PlayerStore) Update(ctx context.Context, player core.Player) (core.Player, error) {
	var tokenUpdated bool
	err := s.store.Update(func(db *db) error {
		for i := range db.Players {
			if db.Players[i].ID == player.ID {
				if db.Players[i].Token != player.Token {
					tokenUpdated = true
				}
				db.Players[i] = unconvertPlayer(player)
				return nil
			}
		}

		return internal.ErrNotFound
	})
	if err != nil {
		return core.Player{}, err
	}

	if tokenUpdated {
		s.bus.PlayerTokenUpdated(ctx, core.EventPlayerTokenUpdated{ID: player.ID})
	}

	return player, nil
}

// Create implements models.PlayerStore.
func (s PlayerStore) Create(ctx context.Context, player core.Player) (core.Player, error) {
	err := s.store.Update(func(db *db) error {
		last, err := lo.Last(db.Players)
		if err != nil {
			player.ID = 1
		} else {
			player.ID = last.ID + 1
		}

		db.Players = append(db.Players, unconvertPlayer(player))

		return nil
	})
	if err != nil {
		return core.Player{}, err
	}

	s.bus.PlayerCreated(ctx, core.EventPlayerCreated{ID: player.ID})

	return player, nil
}

func (s PlayerStore) Get(ctx context.Context, id int64) (core.Player, error) {
	db, err := s.store.Read()
	if err != nil {
		return core.Player{}, err
	}

	player, ok := lo.Find(db.Players, func(player playerModel) bool { return player.ID == id })
	if !ok {
		return core.Player{}, internal.ErrNotFound
	}

	return convertPlayer(player), nil
}

var _ core.PlayerStore = PlayerStore{}
