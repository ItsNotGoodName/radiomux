package jsondb

import (
	"encoding/json"
	"net/url"
	"os"
	"sync"

	"github.com/ItsNotGoodName/radiomux/internal/core"
	"github.com/rs/zerolog/log"
)

type db struct {
	Players []playerModel `json:"players"`
	Presets []presetModel `json:"presets"`
	Sources []sourceModel `json:"sources"`
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
	Slug string `json:"url"`
}

func convertPreset(p presetModel) core.Preset {
	slug, err := url.Parse(p.Slug)
	if err != nil {
		log.Err(err).Caller().Send()
	}
	return core.Preset{
		ID:   p.ID,
		Name: p.Name,
		Slug: slug,
	}
}

func unconvertPreset(p core.Preset) presetModel {
	return presetModel{
		ID:   p.ID,
		Name: p.Name,
		Slug: p.Slug.String(),
	}
}

type sourceModel struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
	File struct {
		Path     string `json:"path"`
		Readonly bool   `json:"readonly"`
	} `json:"file,omitempty"`
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
	db := &db{}

	b, err := os.ReadFile(s.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return db, nil
		}
		return nil, err
	}

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

	b, err := json.MarshalIndent(db, "", "  ")
	if err != nil {
		return err
	}

	filePathTmp := s.filePath + ".tmp"
	err = os.WriteFile(filePathTmp, b, 0600)
	if err != nil {
		return err
	}

	return os.Rename(filePathTmp, s.filePath)
}
