package apiws

import (
	"context"
	"encoding/json"
	"errors"

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

func (ps *playerStateVisitor) Visit() ([]byte, error) {
	if ps.popEmpty() {
		// No state needs updating
		return nil, errVisitorEmpty
	}

	evt := openapi.Event{}

	if ps.popRefresh() {
		// State needs a full refresh
		playerStates := ps.stateService.List()
		names := make([]string, len(playerStates))
		for i := range playerStates {
			f, err := ps.playerStore.Get(context.TODO(), playerStates[i].ID)
			if err == nil {
				names[i] = f.Name
			}
		}

		err := evt.MergeEventDataPlayerState(openapi.EventDataPlayerState{
			Data: openapi.ConvertPlayerStates(playerStates, names),
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
		// We just need to send partial state
		partials := []openapi.PlayerStatePartial{}
		for i := range ps.playerStateChanges {
			if ps.playerStateChanges[i].changed == diff.ChangedNone {
				// We don't need to update
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