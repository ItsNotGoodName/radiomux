package androidws

import (
	"context"
	"fmt"

	"github.com/ItsNotGoodName/radiomux/internal/android"
	"github.com/ItsNotGoodName/radiomux/internal/protos"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"
	"nhooyr.io/websocket"
)

type connection struct {
	rpcMailbox *rpcMailbox
	conn       *websocket.Conn
}

func newConnection(conn *websocket.Conn) connection {
	return connection{
		rpcMailbox: newRPCMailbox(),
		conn:       conn,
	}
}

func (c connection) handle(ctx context.Context, controller *android.Controller, busEvent android.BusEvent, id int64) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go c.reader(ctx, cancel, busEvent, id)

	err := c.handleCommand(ctx, android.CommandSyncState{})
	if err != nil {
		log.Err(err).Msg("Failed to sync state")
		return
	}

	err = controller.PlayerReady(id, c.handleCommand)
	if err != nil {
		log.Err(err).Msg("Failed to ready the player")
		return
	}

	<-ctx.Done()
}

func (c connection) reader(ctx context.Context, cancel context.CancelFunc, busEvent android.BusEvent, id int64) {
	defer cancel()

	for {
		// Read event
		_, message, err := c.conn.Read(ctx)
		if err != nil {
			log.Err(err).Msg("Failed to read websocket")
			return
		}

		// Parse event
		var event protos.Event
		if err := proto.Unmarshal(message, &event); err != nil {
			log.Err(err).Msg("Failed to unmarshal")
			return
		}

		// Route event
		switch m := event.Payload.(type) {
		case *protos.Event_RpcReply:
			if replyC, found := c.rpcMailbox.Get(m.RpcReply.Id); found {
				select {
				case replyC <- struct{}{}:
				default:
					// Default because we do not trust the client to only send the rpc reply once
				}
			}
		default:
			err := handleEvent(ctx, id, busEvent, &event)
			if err != nil {
				log.Err(err).Msg("Failed to handle event")
				return
			}
		}
	}
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
	default:
		return fmt.Errorf("received invalid command: %T", m)
	}
}

func (c connection) handleCommand(ctx context.Context, cmd android.Command) error {
	switch cmd.(type) {
	case android.CommandDisconnect:
		return c.conn.Close(websocket.StatusNormalClosure, "CommandDisconnect")
	default:
		return c.handleCommandRPCWithReply(ctx, cmd)
	}
}

// func (c connection) commandRPCNoReply(ctx context.Context, cmd android.Command) error {
// 	rpc := protos.Rpc{}
//
// 	err := convertCommand(cmd, &rpc)
// 	if err != nil {
// 		return err
// 	}
//
// 	// Marshal
// 	b, err := proto.Marshal(&rpc)
// 	if err != nil {
// 		return fmt.Errorf("failed to marshal command: %w", err)
// 	}
//
// 	// Write
// 	err = c.conn.Write(ctx, websocket.MessageBinary, b)
// 	if err != nil {
// 		return fmt.Errorf("failed to write command: %w", err)
// 	}
//
// 	return nil
// }

func (c connection) handleCommandRPCWithReply(ctx context.Context, cmd android.Command) error {
	rpc := protos.Rpc{}

	if err := convertCommand(cmd, &rpc); err != nil {
		return err
	}

	rpcID := c.rpcMailbox.NextID()
	rpc.Id = rpcID

	// Marshal
	b, err := proto.Marshal(&rpc)
	if err != nil {
		return fmt.Errorf("failed to marshal command: %w", err)
	}

	replyC, remove := c.rpcMailbox.Create(rpcID)
	defer remove()

	// Write
	if err := c.conn.Write(ctx, websocket.MessageBinary, b); err != nil {
		return fmt.Errorf("failed to write command: %w", err)
	}

	// Wait for reply
	return rpcWait(ctx, replyC)
}
