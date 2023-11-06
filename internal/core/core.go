package core

import (
	"context"

	"github.com/ItsNotGoodName/radiomux/pkg/pagination"
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
	// List MUST return a sorted list of players by the ID.
	List(ctx context.Context) ([]Player, error)
	Update(ctx context.Context, req Player) (Player, error)
	Delete(ctx context.Context, id int64) error
	Drop(ctx context.Context) ([]Player, error)
}

type FileSource struct {
	ID       int64
	Name     string
	Path     string
	Readonly bool
}

type File struct {
	SourceID  int64
	Path      string
	Directory bool
}

type FileListRequest struct {
	Page pagination.Page
}

type FileListResponse struct {
	Items      []File
	PageResult pagination.PageResult
}

type SubsonicSource struct {
	ID       int64
	Name     string
	Address  string
	Username string
	Password string
}
