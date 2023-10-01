// androidws handle android player connections.
package androidws

import (
	"context"
	"errors"
	"fmt"
	"strconv"

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
		ArtworkData:     m.ArtworkData,
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
