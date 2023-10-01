package apiws

import (
	"encoding/json"
	"errors"

	"github.com/ItsNotGoodName/radiomux/internal"
	"github.com/ItsNotGoodName/radiomux/internal/android"
	"github.com/ItsNotGoodName/radiomux/internal/openapi"
	"github.com/ItsNotGoodName/radiomux/pkg/diff"
	"github.com/rs/zerolog/log"
)

type playerStateChange struct {
	id      int64
	changed diff.Changed
}

type PlayerStateVisitor struct {
	stateService *android.StateService

	Empty              bool
	refresh            bool
	playerStateChanges []playerStateChange
}

func NewPlayerState(stateService *android.StateService) *PlayerStateVisitor {
	return &PlayerStateVisitor{
		stateService:       stateService,
		Empty:              false,
		refresh:            true,
		playerStateChanges: []playerStateChange{},
	}
}

func (ps *PlayerStateVisitor) popEmpty() bool {
	if ps.Empty {
		return true
	}
	ps.Empty = true
	return false
}

func (ps *PlayerStateVisitor) popRefresh() bool {
	if ps.refresh {
		ps.refresh = false
		return true
	}
	return false
}

func (ps *PlayerStateVisitor) queue(refresh bool) {
	ps.Empty = false
	if !ps.refresh {
		ps.refresh = refresh
	}
}

func (ps *PlayerStateVisitor) StateChange(msg android.StateChange) {
	log.Debug().Msgf("StateChange: Start: %+v", ps)
	for i := range ps.playerStateChanges {
		if ps.playerStateChanges[i].id == msg.ID {
			// Old player update
			ps.playerStateChanges[i].changed = ps.playerStateChanges[i].changed.Merge(msg.Changed)
			ps.queue(false)
			log.Debug().Msgf("StateChange: End: %+v", ps)
			return
		}
	}

	// New player was added
	ps.playerStateChanges = append(ps.playerStateChanges, playerStateChange{id: msg.ID, changed: diff.ChangedAll})
	ps.queue(true)
	log.Debug().Msgf("StateChange: End: %+v", ps)
}

func (ps *PlayerStateVisitor) Visit() ([]byte, error) {
	log.Debug().Msgf("Visit: Start: %+v", ps)
	if ps.popEmpty() {
		// No state needs updating
		log.Debug().Msgf("Visit: End: %+v", ps)
		return nil, ErrVisitorEmpty
	}

	evt := openapi.Event{}

	if ps.popRefresh() {
		// State needs a full refresh
		playerStates := ps.stateService.List()

		err := evt.MergeEventDataPlayerState(openapi.EventDataPlayerState{
			Data: openapi.ConvertPlayerStates(playerStates),
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

	log.Debug().Msgf("Visit: End: %+v", ps)
	return json.Marshal(evt)
}
