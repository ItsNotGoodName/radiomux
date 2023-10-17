package android

import (
	"context"
	"sync"

	"github.com/ItsNotGoodName/radiomux/pkg/diff"
	"github.com/ItsNotGoodName/radiomux/pkg/jsonext"
	"github.com/rs/zerolog/log"
)

const (
	StateChangedConnected diff.Changed = 1 << iota
	StateChangedReady
	StateChangedMinVolume
	StateChangedMaxVolume
	StateChangedVolume
	StateChangedMuted
	StateChangedPlabackState
	StateChangedPlaybackError
	StateChangedPlaying
	StateChangedLoading
	StateChangedTitle
	StateChangedGenre
	StateChangedStation
	StateChangedURI
)

type State struct {
	ID            int64
	Connected     bool
	Ready         bool
	MinVolume     int
	MaxVolume     int
	Volume        int
	Muted         bool
	PlaybackState PlaybackState
	PlaybackError PlaybackError
	Playing       bool
	Loading       bool
	Title         string
	Genre         string
	Station       string
	URI           string
}

type StateChange struct {
	ID      int64
	Changed diff.Changed
}

type StateStore interface {
	List() []State
	Get(id int64) (State, error)
	Update(id int64, fn func(state State, changed diff.Changed) (State, diff.Changed)) error
}

type StatePubSub interface {
	Broadcast(id int64, changed diff.Changed)
	Subscribe() (<-chan StateChange, func())
}

func NewState(id int64) State {
	return State{
		ID:            id,
		PlaybackState: PLAYBACK_STATE_IDLE,
		PlaybackError: NewPlaybackError(),
	}
}

func NewStateService(statePubSub StatePubSub, stateStore StateStore) *StateService {
	return &StateService{
		StatePubSub: statePubSub,
		StateStore:  stateStore,
		statesMu:    sync.Mutex{},
		states:      []State{},
	}
}

type StateService struct {
	StatePubSub
	StateStore

	statesMu sync.Mutex
	states   []State
}

// PositionChanged implements BusEvent.
func (*StateService) PositionChanged(ctx context.Context, id int64, event EventPositionChanged) error {
	log.Debug().Int64("id", id).Msg(jsonext.String(event))
	return nil
}

// TimelineWindowChanged implements BusEvent.
func (*StateService) TimelineWindowChanged(ctx context.Context, id int64, event EventTimelineWindowChanged) error {
	log.Debug().Int64("id", id).Msg(jsonext.String(event))
	return nil
}

// CurrentURIChanged implements BusEvent.
func (s *StateService) CurrentURIChanged(ctx context.Context, id int64, event EventCurrentURIChanged) error {
	return s.Update(id, func(state State, changed diff.Changed) (State, diff.Changed) {
		if state.URI != event.URI {
			state.URI = event.URI
			changed = changed.Merge(StateChangedURI)
		}
		return state, changed
	})
}

// PlayerConnectChanged implements ControllerHooks.
func (s *StateService) PlayerConnectChanged(id int64) error {
	return s.Update(id, func(state State, changed diff.Changed) (State, diff.Changed) {
		if !state.Connected {
			state.Connected = true
			changed = changed.Merge(StateChangedConnected)
		}
		return state, changed
	})
}

// PlayerDisconnectChanged implements ControllerHooks.
func (s *StateService) PlayerDisconnectChanged(id int64) error {
	return s.Update(id, func(state State, changed diff.Changed) (State, diff.Changed) {
		return NewState(id), diff.ChangedAll
	})
}

// PlayerReadyChanged implements ControllerHooks.
func (s *StateService) PlayerReadyChanged(id int64) error {
	return s.Update(id, func(state State, changed diff.Changed) (State, diff.Changed) {
		if !state.Ready {
			state.Ready = true
			changed = changed.Merge(StateChangedReady)
		}
		return state, changed
	})
}

// DeviceInfoChanged implements BusEvent.
func (s *StateService) DeviceInfoChanged(ctx context.Context, id int64, event EventDeviceInfoChanged) error {
	// log.Debug().Int64("id", id).Msg(jsonext.String(event))
	return s.Update(id, func(state State, changed diff.Changed) (State, diff.Changed) {
		if state.MinVolume != event.DeviceInfo.MinVolume {
			state.MinVolume = event.DeviceInfo.MinVolume
			changed = changed.Merge(StateChangedMinVolume)
		}
		if state.MaxVolume != event.DeviceInfo.MaxVolume {
			state.MaxVolume = event.DeviceInfo.MaxVolume
			changed = changed.Merge(StateChangedMaxVolume)
		}
		return state, changed
	})
}

// DeviceVolumeChanged implements BusEvent.
func (s *StateService) DeviceVolumeChanged(ctx context.Context, id int64, event EventDeviceVolumeChanged) error {
	// log.Debug().Int64("id", id).Msg(jsonext.String(event))
	return s.Update(id, func(state State, changed diff.Changed) (State, diff.Changed) {
		if state.Volume != event.Volume {
			state.Volume = event.Volume
			changed = changed.Merge(StateChangedVolume)
		}
		if state.Muted != event.Muted {
			state.Muted = event.Muted
			changed = changed.Merge(StateChangedMuted)
		}
		return state, changed
	})
}

// IsLoadingChanged implements BusEvent.
func (s *StateService) IsLoadingChanged(ctx context.Context, id int64, event EventIsLoadingChanged) error {
	// log.Debug().Int64("id", id).Msg(jsonext.String(event))
	return s.Update(id, func(state State, changed diff.Changed) (State, diff.Changed) {
		if state.Loading != event.IsLoading {
			state.Loading = event.IsLoading
			changed = changed.Merge(StateChangedLoading)
		}
		// Clear playback error
		if state.Loading && !state.PlaybackError.Null {
			state.PlaybackError = NewPlaybackError()
			changed = changed.Merge(StateChangedPlaybackError)
		}
		return state, changed
	})
}

// IsPlayingChanged implements BusEvent.
func (s *StateService) IsPlayingChanged(ctx context.Context, id int64, event EventIsPlayingChanged) error {
	// log.Debug().Int64("id", id).Msg(jsonext.String(event))
	return s.Update(id, func(state State, changed diff.Changed) (State, diff.Changed) {
		if state.Playing != event.IsPlaying {
			state.Playing = event.IsPlaying
			changed = changed.Merge(StateChangedPlaying)
		}
		return state, changed
	})
}

// MediaMetadataChanged implements BusEvent.
func (s *StateService) MediaMetadataChanged(ctx context.Context, id int64, event EventMediaMetadataChanged) error {
	// log.Debug().Int64("id", id).Msg(jsonext.String(event))
	return s.Update(id, func(state State, changed diff.Changed) (State, diff.Changed) {
		if state.Title != event.MediaMetadata.Title {
			state.Title = event.MediaMetadata.Title
			changed = changed.Merge(StateChangedTitle)
		}
		if state.Genre != event.MediaMetadata.Genre {
			state.Genre = event.MediaMetadata.Genre
			changed = changed.Merge(StateChangedGenre)
		}
		if state.Station != event.MediaMetadata.Station {
			state.Station = event.MediaMetadata.Station
			changed = changed.Merge(StateChangedStation)
		}
		return state, changed
	})
}

// PlaybackParametersChanged implements BusEvent.
func (*StateService) PlaybackParametersChanged(ctx context.Context, id int64, event EventPlaybackParametersChanged) error {
	return nil
}

// PlaybackStateChanged implements BusEvent.
func (s *StateService) PlaybackStateChanged(ctx context.Context, id int64, event EventPlaybackStateChanged) error {
	// log.Debug().Msg(jsonext.String(event))
	return s.Update(id, func(state State, changed diff.Changed) (State, diff.Changed) {
		if state.PlaybackState != event.PlaybackState {
			state.PlaybackState = event.PlaybackState
			changed = StateChangedPlabackState
		}
		return state, changed
	})
}

// PlayerError implements BusEvent.
func (s *StateService) PlayerError(ctx context.Context, id int64, event EventPlayerError) error {
	// log.Debug().Int64("id", id).Msg(jsonext.String(event))
	return s.Update(id, func(state State, changed diff.Changed) (State, diff.Changed) {
		state.PlaybackError = event.PlaybackError
		return state, changed.Merge(StateChangedPlaybackError)
	})
}

// PlaylistMetadataChanged implements BusEvent.
func (*StateService) PlaylistMetadataChanged(ctx context.Context, id int64, event EventPlaylistMetadataChanged) error {
	return nil
}

// VolumeChanged implements BusEvent.
func (s *StateService) VolumeChanged(ctx context.Context, id int64, event EventVolumeChanged) error {
	return nil
}

var _ BusEvent = (*StateService)(nil)
var _ ControllerHook = (*StateService)(nil)
