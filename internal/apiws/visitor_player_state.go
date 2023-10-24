package apiws

import (
	"cmp"
	"context"
	"encoding/json"
	"errors"
	"slices"

	"github.com/ItsNotGoodName/radiomux/internal"
	"github.com/ItsNotGoodName/radiomux/internal/android"
	"github.com/ItsNotGoodName/radiomux/internal/core"
	"github.com/ItsNotGoodName/radiomux/internal/openapi"
	"github.com/ItsNotGoodName/radiomux/pkg/diff"
)

type playerStateChange struct {
	id      int64
	changed diff.Changed
}

type playerStateVisitor struct {
	stateService *android.StateService
	playerStore  core.PlayerStore

	empty              bool
	refresh            bool
	playerStateChanges []playerStateChange
}

func newPlayerStateVisitor(stateService *android.StateService, playerStore core.PlayerStore) *playerStateVisitor {
	return &playerStateVisitor{
		stateService:       stateService,
		playerStore:        playerStore,
		empty:              false,
		refresh:            true,
		playerStateChanges: []playerStateChange{},
	}
}

func (ps *playerStateVisitor) popEmpty() bool {
	if ps.empty {
		return true
	}
	ps.empty = true
	return false
}

func (ps *playerStateVisitor) popRefresh() bool {
	if ps.refresh {
		ps.refresh = false
		return true
	}
	return false
}

func (ps *playerStateVisitor) queue(refresh bool) {
	ps.empty = false
	if !ps.refresh {
		ps.refresh = refresh
	}
}

func (ps *playerStateVisitor) StateChange(msg android.StateChange) {
	for i := range ps.playerStateChanges {
		if ps.playerStateChanges[i].id == msg.ID {
			// Old player update
			ps.playerStateChanges[i].changed = ps.playerStateChanges[i].changed.Merge(msg.Changed)
			ps.queue(false)
			return
		}
	}

	// New player was added
	ps.playerStateChanges = append(ps.playerStateChanges, playerStateChange{id: msg.ID, changed: diff.ChangedAll})
	ps.queue(true)
}

func (ps *playerStateVisitor) Visit(ctx context.Context) ([]byte, error) {
	if ps.popEmpty() {
		// No state needs updating
		return nil, errVisitorEmpty
	}

	evt := openapi.Event{}

	if ps.popRefresh() {
		// State needs a full refresh

		// Get players and states
		playerStates := ps.stateService.List()
		playerStatePlayers := make([]core.Player, len(playerStates))
		{
			players, err := ps.playerStore.List(ctx)
			if err != nil {
				return nil, err
			}
			for i := range playerStates {
				index, found := slices.BinarySearchFunc(players, playerStates[i], func(p core.Player, s android.State) int {
					return cmp.Compare(p.ID, s.ID)
				})
				if found {
					playerStatePlayers[index] = players[index]
				}
			}
		}

		// Create event
		err := evt.MergeEventDataPlayerState(openapi.EventDataPlayerState{
			Data: openapi.ConvertPlayerStates(playerStates, playerStatePlayers),
		})
		if err != nil {
			return nil, err
		}

		// Blow away old state changes
		states := make([]playerStateChange, 0, len(playerStates))
		for _, s := range playerStates {
			states = append(states, playerStateChange{id: s.ID, changed: diff.ChangedNone})
		}
		ps.playerStateChanges = states
	} else {
		// Just need to send partial state

		partials := []openapi.PlayerStatePartial{}
		for i := range ps.playerStateChanges {
			if ps.playerStateChanges[i].changed == diff.ChangedNone {
				// Don't need to update
				continue
			}

			state, err := ps.stateService.Get(ps.playerStateChanges[i].id)
			if err != nil {
				if errors.Is(err, internal.ErrNotFound) {
					// Edge case
					ps.queue(true)
					continue
				}

				return nil, err
			}

			partials = append(partials, openapi.ConvertPlayerStatePartial(&state, ps.playerStateChanges[i].changed))
			ps.playerStateChanges[i].changed = diff.ChangedNone
		}

		err := evt.MergeEventDataPlayerStatePartial(openapi.EventDataPlayerStatePartial{
			Data: partials,
		})
		if err != nil {
			return nil, err
		}
	}

	return json.Marshal(evt)
}

func (ps playerStateVisitor) HasMore() bool {
	return !ps.empty
}
