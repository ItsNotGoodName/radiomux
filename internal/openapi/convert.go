package openapi

import (
	"errors"

	"github.com/ItsNotGoodName/radiomux/internal"
	"github.com/ItsNotGoodName/radiomux/internal/android"
	"github.com/ItsNotGoodName/radiomux/pkg/diff"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func ConvertErr(err error) error {
	if errors.Is(err, android.ErrPlayerNotConnected) {
		return echo.ErrNotFound.WithInternal(err)
	}
	if errors.Is(err, internal.ErrNotFound) {
		return echo.ErrNotFound.WithInternal(err)
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

func ConvertPlayerStates(states []android.State, names []string) []PlayerState {
	players := make([]PlayerState, 0, len(states))
	for i, s := range states {
		players = append(players, ConvertPlayerState(&s, names[i]))
	}
	return players
}

func ConvertPlayerState(s *android.State, name string) PlayerState {
	return PlayerState{
		Id:            s.ID,
		Name:          name,
		Connected:     s.Connected,
		Ready:         s.Ready,
		MinVolume:     s.MinVolume,
		MaxVolume:     s.MaxVolume,
		Volume:        s.Volume,
		Muted:         s.Muted,
		PlaybackState: ConvertPlayerPlaybackState(s.PlaybackState),
		PlaybackError: ConvertPlaybackError(s.PlaybackError),
		Playing:       s.Playing,
		Loading:       s.Loading,
		Title:         s.Title,
		Genre:         s.Genre,
		Station:       s.Station,
		Uri:           s.URI,
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
	return p
}
