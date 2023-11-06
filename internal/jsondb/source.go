package jsondb

import (
	"context"
	"fmt"
	"slices"

	"github.com/ItsNotGoodName/radiomux/internal"
	"github.com/ItsNotGoodName/radiomux/internal/core"
)

func NewSourceStore(store *Store) SourceStore {
	return SourceStore{
		store: store,
	}
}

type SourceStore struct {
	store *Store
}

func (s SourceStore) GetFileSource(ctx context.Context, id int64) (core.FileSource, error) {
	db, err := s.store.Read()
	if err != nil {
		return core.FileSource{}, err
	}

	index := slices.IndexFunc(db.Sources, func(sm sourceModel) bool { return sm.ID == id })
	if index == -1 {
		return core.FileSource{}, internal.ErrNotFound
	}

	if db.Sources[index].Type != "file" {
		return core.FileSource{}, fmt.Errorf("invalid file source type: %s", db.Sources[index].Type)
	}

	return core.FileSource{
		ID:       db.Sources[index].ID,
		Name:     db.Sources[index].Name,
		Path:     db.Sources[index].File.Path,
		Readonly: db.Sources[index].File.Readonly,
	}, nil
}
