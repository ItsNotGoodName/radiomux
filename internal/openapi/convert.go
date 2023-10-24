package openapi

import (
	"encoding/json"
	"errors"

	"github.com/ItsNotGoodName/radiomux/internal"
	"github.com/ItsNotGoodName/radiomux/internal/android"
	"github.com/ItsNotGoodName/radiomux/internal/core"
	"github.com/ItsNotGoodName/radiomux/pkg/diff"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func ConvertNotification(notification Notification) ([]byte, error) {
	evt := Event{}
	err := evt.MergeEventDataNotification(EventDataNotification{
		Data: notification,
	})
	if err != nil {
		return nil, err
	}
	return json.Marshal(evt)
}

func ConvertErr(err error) error {
	if errors.Is(err, android.ErrPlayerNotConnected) {
		return echo.ErrNotFound.WithInternal(err)
	}
	if errors.Is(err, internal.ErrNotFound) {
		return echo.ErrNotFound.WithInternal(err)
	}
	if errors.Is(err, errors.ErrUnsupported) {
		return echo.ErrNotImplemented.WithInternal(err)
	}
	if errors.Is(err, android.ErrPlayerNotReady) {
		return echo.ErrTooEarly.WithInternal(err)
	}

	return err
}

func ConvertPlayerPlaybackState(s android.PlaybackState) PlayerPlaybackState {
	switch s {
	case android.PLAYBACK_STATE_IDLE:
		return IDLE
	case android.PLAYBACK_STATE_BUFFERING:
		return BUFFERING
	case android.PLAYBACK_STATE_READY:
		return READY
	case android.PLAYBACK_STATE_ENDED:
		return ENDED
	default:
		log.Error().Msg("invalid playback state")
		return IDLE
	}
}

func ConvertPlaybackError(e android.PlaybackError) string {
	if e.Null {
		return ""
	}
	return e.String()
}

func ConvertPlayerStates(states []android.State, players []core.Player) []PlayerState {
	playersStates := make([]PlayerState, 0, len(states))
	for i, s := range states {
		playersStates = append(playersStates, ConvertPlayerState(&s, players[i]))
	}
	return playersStates
}

func ConvertPlayerState(s *android.State, p core.Player) PlayerState {
	return PlayerState{
		Id:                      s.ID,
		Name:                    p.Name,
		Connected:               s.Connected,
		Ready:                   s.Ready,
		MinVolume:               s.MinVolume,
		MaxVolume:               s.MaxVolume,
		Volume:                  s.Volume,
		Muted:                   s.Muted,
		PlaybackState:           ConvertPlayerPlaybackState(s.PlaybackState),
		PlaybackError:           ConvertPlaybackError(s.PlaybackError),
		Playing:                 s.Playing,
		Loading:                 s.Loading,
		Title:                   s.Title,
		Genre:                   s.Genre,
		Station:                 s.Station,
		Uri:                     s.URI,
		TimelineIsSeekable:      s.TimelineIsSeekable,
		TimelineIsLive:          s.TimelineIsLive,
		TimelineIsPlaceholder:   s.TimelineIsPlaceholder,
		TimelineDefaultPosition: s.TimelineDefaultPosition.Milliseconds(),
		TimelineDuration:        s.TimelineDuration.Milliseconds(),
		Position:                s.Position.Milliseconds(),
		PositionTime:            s.PositionTime,
	}
}

func ConvertPlayerStatePartial(s *android.State, c diff.Changed) PlayerStatePartial {
	p := PlayerStatePartial{Id: s.ID}

	if c.Is(android.StateChangedConnected) {
		p.Connected = &s.Connected
	}
	if c.Is(android.StateChangedReady) {
		p.Ready = &s.Ready
	}
	if c.Is(android.StateChangedMinVolume) {
		p.MinVolume = &s.MinVolume
	}
	if c.Is(android.StateChangedMaxVolume) {
		p.MaxVolume = &s.MaxVolume
	}
	if c.Is(android.StateChangedVolume) {
		p.Volume = &s.Volume
	}
	if c.Is(android.StateChangedMuted) {
		p.Muted = &s.Muted
	}
	if c.Is(android.StateChangedPlabackState) {
		tmp := ConvertPlayerPlaybackState(s.PlaybackState)
		p.PlaybackState = &tmp
	}
	if c.Is(android.StateChangedPlaybackError) {
		tmp := ConvertPlaybackError(s.PlaybackError)
		p.PlaybackError = &tmp
	}
	if c.Is(android.StateChangedPlaying) {
		p.Playing = &s.Playing
	}
	if c.Is(android.StateChangedLoading) {
		p.Loading = &s.Loading
	}
	if c.Is(android.StateChangedTitle) {
		p.Title = &s.Title
	}
	if c.Is(android.StateChangedGenre) {
		p.Genre = &s.Genre
	}
	if c.Is(android.StateChangedStation) {
		p.Station = &s.Station
	}
	if c.Is(android.StateChangedURI) {
		p.Uri = &s.URI
	}
	if c.Is(android.StateTimelineIsSeekable) {
		p.TimelineIsSeekable = &s.TimelineIsSeekable
	}
	if c.Is(android.StateTimelineIsLive) {
		p.TimelineIsLive = &s.TimelineIsLive
	}
	if c.Is(android.StateTimelineIsPlaceholder) {
		p.TimelineIsPlaceholder = &s.TimelineIsPlaceholder
	}
	if c.Is(android.StateTimelineDefaultPosition) {
		val := s.TimelineDefaultPosition.Milliseconds()
		p.TimelineDefaultPosition = &val
	}
	if c.Is(android.StateTimelineDuration) {
		val := s.TimelineDuration.Milliseconds()
		p.TimelineDuration = &val
	}
	if c.Is(android.StatePosition) {
		val := s.Position.Milliseconds()
		p.Position = &val
	}
	if c.Is(android.StatePositionTime) {
		p.PositionTime = &s.PositionTime
	}
	return p
}
