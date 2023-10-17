// androidws handle android player connections.
package androidws

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/ItsNotGoodName/radiomux/internal"
	"github.com/ItsNotGoodName/radiomux/internal/android"
	"github.com/ItsNotGoodName/radiomux/internal/core"
	"github.com/ItsNotGoodName/radiomux/internal/protos"
	"github.com/labstack/echo/v4"
)

func auth(ctx context.Context, c echo.Context, playerStore core.PlayerStore) (int64, error) {
	id, err := strconv.ParseInt(c.QueryParam("id"), 10, 64)
	if err != nil {
		return 0, echo.ErrBadRequest.WithInternal(err)
	}

	token := c.QueryParam("token")

	p, err := playerStore.Get(ctx, id)
	if err != nil {
		if errors.Is(err, internal.ErrNotFound) {
			return 0, echo.ErrNotFound.WithInternal(err)
		}
		return 0, err
	}
	if p.CompareToken(token) {
		return 0, echo.ErrUnauthorized.WithInternal(err)
	}

	return id, nil
}

func handleEvent(ctx context.Context, id int64, busEvent android.BusEvent, event *protos.Event) error {
	switch m := event.Payload.(type) {
	case *protos.Event_OnMediaMetadataChanged:
		return busEvent.MediaMetadataChanged(ctx, id, android.EventMediaMetadataChanged{
			MediaMetadata: convertMediaMetadata(m.OnMediaMetadataChanged.GetMediaMetadata()),
		})
	case *protos.Event_OnPlaylistMetadataChanged:
		return busEvent.PlaylistMetadataChanged(ctx, id, android.EventPlaylistMetadataChanged{
			MediaMetadata: convertMediaMetadata(m.OnPlaylistMetadataChanged.GetMediaMetadata()),
		})
	case *protos.Event_OnIsLoadingChanged:
		return busEvent.IsLoadingChanged(ctx, id, android.EventIsLoadingChanged{
			IsLoading: m.OnIsLoadingChanged.GetIsLoading(),
		})
	case *protos.Event_OnPlaybackStateChanged:
		return busEvent.PlaybackStateChanged(ctx, id, android.EventPlaybackStateChanged{
			PlaybackState: android.PlaybackState(m.OnPlaybackStateChanged.GetPlaybackState()),
		})
	case *protos.Event_OnIsPlayingChanged:
		return busEvent.IsPlayingChanged(ctx, id, android.EventIsPlayingChanged{
			IsPlaying: m.OnIsPlayingChanged.GetIsPlaying(),
		})
	case *protos.Event_OnPlayerError:
		error := m.OnPlayerError.GetError()
		if error == nil {
			return nil
		}
		return busEvent.PlayerError(ctx, id, android.EventPlayerError{
			PlaybackError: android.PlaybackError{
				Code:        android.PlaybackCode(error.GetErrorCode()),
				TimestampMs: error.GetTimestampMs(),
			},
		})
	case *protos.Event_OnPlaybackParametersChanged:
		return busEvent.PlaybackParametersChanged(ctx, id, android.EventPlaybackParametersChanged{
			PlaybackParameters: android.PlaybackParameters{
				Speed: int(m.OnPlaybackParametersChanged.GetPlaybackParameters().GetSpeed()),
				Pitch: int(m.OnPlaybackParametersChanged.GetPlaybackParameters().GetPitch()),
			},
		})
	case *protos.Event_OnVolumeChanged:
		return busEvent.VolumeChanged(ctx, id, android.EventVolumeChanged{
			Volume: float64(m.OnVolumeChanged.GetVolume()),
		})
	case *protos.Event_OnDeviceInfoChanged:
		return busEvent.DeviceInfoChanged(ctx, id, android.EventDeviceInfoChanged{
			DeviceInfo: android.DeviceInfo{
				MinVolume: int(m.OnDeviceInfoChanged.GetDeviceInfo().GetMinVolume()),
				MaxVolume: int(m.OnDeviceInfoChanged.GetDeviceInfo().GetMaxVolume()),
			},
		})
	case *protos.Event_OnDeviceVolumeChanged:
		return busEvent.DeviceVolumeChanged(ctx, id, android.EventDeviceVolumeChanged{
			Volume: int(m.OnDeviceVolumeChanged.GetVolume()),
			Muted:  m.OnDeviceVolumeChanged.GetMuted(),
		})
	case *protos.Event_OnCurrentUriChanged:
		return busEvent.CurrentURIChanged(ctx, id, android.EventCurrentURIChanged{
			URI: m.OnCurrentUriChanged.GetUri(),
		})
	case *protos.Event_OnTimelineChanged:
		window := m.OnTimelineChanged.GetWindow()
		if window == nil {
			return nil
		}
		return busEvent.TimelineWindowChanged(ctx, id, android.EventTimelineWindowChanged{
			Window: android.TimelineWindow{
				IsSeekable:      window.IsSeekable,
				IsDynamic:       window.IsDynamic,
				IsLive:          window.IsLive,
				IsPlaceholder:   window.IsPlaceholder,
				DefaultPosition: time.Duration(window.DefaultPositionMs) * time.Millisecond,
				Duration:        time.Duration(window.DurationMs) * time.Millisecond,
			},
			TimeUnset: m.OnTimelineChanged.TimeUnset,
		})
	case *protos.Event_OnPositionChanged:
		timestamp, err := time.Parse(time.RFC3339, m.OnPositionChanged.GetTimestamp())
		if err != nil {
			return err
		}
		return busEvent.PositionChanged(ctx, id, android.EventPositionChanged{
			OldPositionInfo: android.PositionInfo{
				Position: time.Duration(m.OnPositionChanged.GetOldPosition().PositionMs) * time.Millisecond,
			},
			NewPositionInfo: android.PositionInfo{
				Position: time.Duration(m.OnPositionChanged.GetNewPosition().PositionMs) * time.Millisecond,
			},
			Time: timestamp,
		})
	default:
		return fmt.Errorf("received invalid command: %T", m)
	}
}

func convertMediaMetadata(x *protos.MediaMetadata) *android.MediaMetadata {
	if x == nil {
		return &android.MediaMetadata{}
	}
	m := x

	return &android.MediaMetadata{
		Title:        m.Title,
		Artist:       m.Artist,
		AlbumTitle:   m.AlbumTitle,
		AlbumArtist:  m.AlbumArtist,
		DisplayTitle: m.DisplayTitle,
		Subtitle:     m.Subtitle,
		Description:  m.Description,
		UserRating: android.Rating{
			IsRated:    m.UserRating.GetIsRated(),
			RatingType: android.RatingType(m.UserRating.GetRatingType()),
		},
		OverallRating: android.Rating{
			IsRated:    m.OverallRating.GetIsRated(),
			RatingType: android.RatingType(m.OverallRating.GetRatingType()),
		},
		// ArtworkData:     m.ArtworkData,
		ArtworkDataType: android.PictureType(m.ArtworkDataType),
		ArtworkUri:      m.ArtworkUri,
		TrackNumber:     int(m.TrackNumber),
		TotalTrackCount: int(m.TotalTrackCount),
		IsBrowsable:     m.IsBrowsable,
		IsPlayable:      m.IsPlayable,
		RecordingYear:   int(m.RecordingYear),
		RecordingMonth:  int(m.RecordingMonth),
		RecordingDay:    int(m.RecordingDay),
		ReleaseYear:     int(m.ReleaseYear),
		ReleaseMonth:    int(m.ReleaseMonth),
		ReleaseDay:      int(m.ReleaseDay),
		Writer:          m.Writer,
		Composer:        m.Composer,
		Conductor:       m.Conductor,
		DiscNumber:      int(m.DiscNumber),
		TotalDiscCount:  int(m.TotalDiscCount),
		Genre:           m.Genre,
		Compilation:     m.Compilation,
		Station:         m.Station,
		MediaType:       android.MediaType(m.MediaType),
	}
}

func convertCommand(cmd android.Command, rpc *protos.Rpc) error {
	switch t := cmd.(type) {
	case android.CommandStop:
		rpc.Payload = &protos.Rpc_Stop{}
	case android.CommandPlay:
		rpc.Payload = &protos.Rpc_Play{}
	case android.CommandPause:
		rpc.Payload = &protos.Rpc_Pause{}
	case android.CommandPrepare:
		rpc.Payload = &protos.Rpc_Prepare{}
	case android.CommandSetPlayWhenReady:
		rpc.Payload = &protos.Rpc_SetPlayWhenReady{
			SetPlayWhenReady: &protos.SetPlayWhenReady{
				PlayWhenReady: t.PlayWhenReady,
			},
		}
	case android.CommandSetMediaItem:
		rpc.Payload = &protos.Rpc_SetMediaItem{
			SetMediaItem: &protos.SetMediaItem{
				Uri: t.URI,
			},
		}
	case android.CommandSetVolume:
		rpc.Payload = &protos.Rpc_SetVolume{
			SetVolume: &protos.SetVolume{
				Volume: float32(t.Volume),
			},
		}
	case android.CommandSetDeviceVolume:
		rpc.Payload = &protos.Rpc_SetDeviceVolume{
			SetDeviceVolume: &protos.SetDeviceVolume{
				Volume: int32(t.Volume),
			},
		}
	case android.CommandIncreaseDeviceVolume:
		rpc.Payload = &protos.Rpc_IncreaseDeviceVolume{}
	case android.CommandDecreaseDeviceVolume:
		rpc.Payload = &protos.Rpc_DecreaseDeviceVolume{}
	case android.CommandSetDeviceMuted:
		rpc.Payload = &protos.Rpc_SetDeviceMuted{
			SetDeviceMuted: &protos.SetDeviceMuted{
				Muted: t.Muted,
			},
		}
	case android.CommandSyncState:
		rpc.Payload = &protos.Rpc_SyncState{}
	case android.CommandSeekToDefaultPosition:
		rpc.Payload = &protos.Rpc_SeekToDefaultPosition{}
	default:
		return fmt.Errorf("received invalid command: %T", cmd)
	}

	return nil
}
