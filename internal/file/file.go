package file

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/ItsNotGoodName/radiomux/internal/core"
)

type db struct {
	Players []playerModel `json:"players"`
	Presets []presetModel `json:"presets"`
}

type playerModel struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Token string `json:"token"`
}

func convertPlayer(p playerModel) core.Player {
	return core.Player{
		ID:    p.ID,
		Name:  p.Name,
		Token: p.Token,
	}
}

func unconvertPlayer(p core.Player) playerModel {
	return playerModel{
		ID:    p.ID,
		Name:  p.Name,
		Token: p.Token,
	}
}

type presetModel struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

func convertPreset(p presetModel) core.Preset {
	return core.Preset{
		ID:   p.ID,
		Name: p.Name,
		URL:  p.URL,
	}
}

func unconvertPreset(p core.Preset) presetModel {
	return presetModel{
		ID:   p.ID,
		Name: p.Name,
		URL:  p.URL,
	}
}

type Store struct {
	filePath string
	updateMu sync.Mutex
}

func NewStore(filePath string) *Store {
	return &Store{
		filePath: filePath,
		updateMu: sync.Mutex{},
	}
}

func (s *Store) Read() (*db, error) {
	b, err := os.ReadFile(s.filePath)
	if err != nil {
		return nil, err
	}

	db := &db{}
	if err := json.Unmarshal(b, db); err != nil {
		return nil, err
	}

	return db, nil
}

func (s *Store) Update(fn func(db *db) error) error {
	s.updateMu.Lock()
	defer s.updateMu.Unlock()

	db, err := s.Read()
	if err != nil {
		return err
	}

	if err := fn(db); err != nil {
		return err
	}

	b, err := json.Marshal(db)
	if err != nil {
		return err
	}

	err = os.WriteFile(s.filePath, b, 0600)
	if err != nil {
		return err
	}

	return nil
}
