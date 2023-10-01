package android

import (
	"context"
)

// BusEvent is used by the player to emit events.
type BusEvent interface {
	MediaMetadataChanged(ctx context.Context, id int64, event EventMediaMetadataChanged) error
	PlaylistMetadataChanged(ctx context.Context, id int64, event EventPlaylistMetadataChanged) error
	IsLoadingChanged(ctx context.Context, id int64, event EventIsLoadingChanged) error
	PlaybackStateChanged(ctx context.Context, id int64, event EventPlaybackStateChanged) error
	IsPlayingChanged(ctx context.Context, id int64, event EventIsPlayingChanged) error
	PlayerError(ctx context.Context, id int64, event EventPlayerError) error
	PlaybackParametersChanged(ctx context.Context, id int64, event EventPlaybackParametersChanged) error
	VolumeChanged(ctx context.Context, id int64, event EventVolumeChanged) error
	DeviceInfoChanged(ctx context.Context, id int64, event EventDeviceInfoChanged) error
	DeviceVolumeChanged(ctx context.Context, id int64, event EventDeviceVolumeChanged) error
	CurrentURIChanged(ctx context.Context, id int64, event EventCurrentURIChanged) error
}

type EventMediaMetadataChanged struct {
	MediaMetadata *MediaMetadata
}

type EventPlaylistMetadataChanged struct {
	MediaMetadata *MediaMetadata
}

type EventIsLoadingChanged struct {
	IsLoading bool
}

type EventPlaybackStateChanged struct {
	PlaybackState PlaybackState
}

type EventIsPlayingChanged struct {
	IsPlaying bool
}

type EventPlayerError struct {
	PlaybackError PlaybackError
}

type EventPlaybackParametersChanged struct {
	PlaybackParameters PlaybackParameters
}

type EventVolumeChanged struct {
	Volume float64
	Muted  bool
}

type EventDeviceInfoChanged struct {
	DeviceInfo DeviceInfo
}

type EventDeviceVolumeChanged struct {
	Volume int
	Muted  bool
}

type EventCurrentURIChanged struct {
	URI string
}
