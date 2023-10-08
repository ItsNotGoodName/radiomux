package rpc

import (
	"context"
	"errors"

	"github.com/ItsNotGoodName/radiomux/internal"
	"github.com/ItsNotGoodName/radiomux/internal/core"
	"github.com/ItsNotGoodName/radiomux/internal/webrpc"
)

func NewPlayerService(playerStore core.PlayerStore) *PlayerService {
	return &PlayerService{
		playerStore: playerStore,
	}
}

type PlayerService struct {
	playerStore core.PlayerStore
}

func (s *PlayerService) PlayerCreate(ctx context.Context, req *webrpc.CreatePlayer) (int64, error) {
	token, err := core.GenerateToken()
	if err != nil {
		return 0, webrpc.ConvertErr(err)
	}

	player, err := s.playerStore.Create(ctx, core.Player{
		Name:  req.Name,
		Token: token,
	})
	if err != nil {
		return 0, webrpc.ConvertErr(err)
	}

	return player.ID, nil
}

func (s *PlayerService) PlayerGet(ctx context.Context, id int64) (*webrpc.Player, error) {
	player, err := s.playerStore.Get(ctx, id)
	if err != nil {
		return nil, webrpc.ConvertErr(err)
	}

	return webrpc.ConvertPlayer(player), nil
}

func (s *PlayerService) PlayerList(ctx context.Context) (*webrpc.PlayerListResult, error) {
	players, err := s.playerStore.List(ctx)
	if err != nil {
		return nil, webrpc.ConvertErr(err)
	}

	return &webrpc.PlayerListResult{
		Players: webrpc.ConvertPlayers(players),
		Count:   len(players),
	}, nil
}

func (s *PlayerService) PlayerUpdate(ctx context.Context, req *webrpc.UpdatePlayer) error {
	player, err := s.playerStore.Get(ctx, req.Id)
	if err != nil {
		return webrpc.ConvertErr(err)
	}
	if req.Name != nil {
		player.Name = *req.Name
	}
	player, err = s.playerStore.Update(ctx, player)
	if err != nil {
		return webrpc.ConvertErr(err)
	}

	return nil
}

func (s *PlayerService) PlayerTokenRegenerate(ctx context.Context, id int64) error {
	player, err := s.playerStore.Get(ctx, id)
	if err != nil {
		return webrpc.ConvertErr(err)
	}

	player.Token, err = core.GenerateToken()
	if err != nil {
		return webrpc.ConvertErr(err)
	}

	player, err = s.playerStore.Update(ctx, player)
	if err != nil {
		return webrpc.ConvertErr(err)
	}

	return nil
}

func (s *PlayerService) PlayerDelete(ctx context.Context, ids []int64) error {
	for _, id := range ids {
		err := s.playerStore.Delete(ctx, id)
		if err != nil && !errors.Is(err, internal.ErrNotFound) {
			return webrpc.ConvertErr(err)
		}
	}

	return nil
}

var _ webrpc.PlayerService = (*PlayerService)(nil)
