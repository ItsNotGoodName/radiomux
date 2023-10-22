package androidmock

import (
	"context"
	"errors"
	"sync"

	"github.com/ItsNotGoodName/radiomux/internal/android"
)

type state struct {
	Volume int
	Muted  bool
}

type stateStore struct {
	mu    sync.Mutex
	state state
}

func (s *stateStore) get() state {
	s.mu.Lock()
	state := s.state
	s.mu.Unlock()
	return state
}

func (s *stateStore) update(fn func(s *state)) state {
	s.mu.Lock()
	fn(&s.state)
	state := s.state
	s.mu.Unlock()
	return state
}

type Mock struct {
	id         int64
	controller *android.Controller
	busEvent   android.BusEvent
}

func NewMock(id int64, controller *android.Controller, busEvent android.BusEvent) Mock {
	return Mock{
		id:         id,
		controller: controller,
		busEvent:   busEvent,
	}
}

func (m Mock) Serve(ctx context.Context) error {
	disconnect, err := m.controller.PlayerConnect(m.id)
	if err != nil {
		return err
	}
	defer disconnect()

	s := &stateStore{}

	err = m.busEvent.DeviceInfoChanged(ctx, m.id, android.EventDeviceInfoChanged{
		DeviceInfo: android.DeviceInfo{
			MinVolume: 0,
			MaxVolume: 15,
		},
	})
	if err != nil {
		return err
	}

	err = m.controller.PlayerReady(m.id, func(ctx context.Context, cmd android.Command) error {
		switch t := cmd.(type) {
		case android.CommandStop:
			m.busEvent.IsPlayingChanged(ctx, m.id, android.EventIsPlayingChanged{IsPlaying: false})
			return nil
		case android.CommandPlay:
			return m.busEvent.IsPlayingChanged(ctx, m.id, android.EventIsPlayingChanged{IsPlaying: true})
		case android.CommandPause:
			return m.busEvent.IsPlayingChanged(ctx, m.id, android.EventIsPlayingChanged{IsPlaying: false})
		case android.CommandPrepare:
		case android.CommandSetPlayWhenReady:
		case android.CommandSetMediaItem:
			err := m.busEvent.CurrentURIChanged(ctx, m.id, android.EventCurrentURIChanged{
				URI: t.URI,
			})
			if err != nil {
				return err
			}
			return m.busEvent.MediaMetadataChanged(ctx, m.id, android.EventMediaMetadataChanged{
				MediaMetadata: &android.MediaMetadata{
					Title: t.URI,
				},
			})
		case android.CommandSetVolume:
		case android.CommandSetDeviceVolume:
		case android.CommandIncreaseDeviceVolume:
			s := s.update(func(s *state) {
				s.Volume += 1
				s.Muted = false
			})
			return m.busEvent.DeviceVolumeChanged(ctx, m.id, android.EventDeviceVolumeChanged{Volume: s.Volume, Muted: s.Muted})
		case android.CommandDecreaseDeviceVolume:
			s := s.update(func(s *state) {
				s.Volume -= 1
				s.Muted = false
			})
			return m.busEvent.DeviceVolumeChanged(ctx, m.id, android.EventDeviceVolumeChanged{Volume: s.Volume, Muted: s.Muted})
		case android.CommandSetDeviceMuted:
			s := s.update(func(s *state) {
				s.Muted = t.Muted
				if s.Muted {
					s.Volume = 0
				}
			})
			return m.busEvent.DeviceVolumeChanged(ctx, m.id, android.EventDeviceVolumeChanged{Volume: s.Volume, Muted: s.Muted})
		case android.CommandSyncState:
		case android.CommandSeekToDefaultPosition:
		default:
		}
		return errors.ErrUnsupported
	})
	if err != nil {
		return err
	}

	<-ctx.Done()

	return nil
}
