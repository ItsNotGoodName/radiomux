package androidws

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ItsNotGoodName/radiomux/internal/android"
	"github.com/ItsNotGoodName/radiomux/internal/protos"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"
)

var ErrClosed = errors.New("closed")

type connection struct {
	rpcMailbox *rpcMailbox
	writeC     chan []byte
	closeC     chan struct{}
	doneC      chan struct{}
}

func newConnection() connection {
	return connection{
		rpcMailbox: newRPCMailbox(),
		writeC:     make(chan []byte),
		closeC:     make(chan struct{}, 1),
		doneC:      make(chan struct{}),
	}
}

func (c connection) handle(ctx context.Context, id int64, conn *websocket.Conn, controller *android.Controller, busEvent android.BusEvent) {
	defer close(c.doneC)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	log := log.With().
		Str("package", "androidws").
		Str("remote", conn.RemoteAddr().String()).
		Logger()

	// Read pump
	go c.reader(ctx, cancel, conn, log, busEvent, id)
	// Write pump
	go wsWriter(ctx, cancel, conn, log, c.writeC)

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

	select {
	case <-c.closeC:
		conn.Close()
	case <-ctx.Done():
	}
}

func (c connection) reader(ctx context.Context, cancel context.CancelFunc, conn *websocket.Conn, log zerolog.Logger, busEvent android.BusEvent, id int64) {
	defer cancel()

	conn.SetReadLimit(wsMaxMessageSize)
	conn.SetReadDeadline(time.Now().Add(wsPongWait))
	conn.SetPongHandler(func(string) error { conn.SetReadDeadline(time.Now().Add(wsPongWait)); return nil })

	for {
		// Read command or end on error
		_, data, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure, websocket.CloseNormalClosure) {
				log.Err(err).Msg("Failed to read from WebSocket")
			}
			return
		}

		// Parse event
		var event protos.Event
		if err := proto.Unmarshal(data, &event); err != nil {
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

func (c connection) handleCommand(ctx context.Context, cmd android.Command) error {
	switch cmd.(type) {
	case android.CommandDisconnect:
		select {
		case c.closeC <- struct{}{}:
		default:
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-c.doneC:
			return nil
		}
	default:
		return c.handleCommandRPCWithReply(ctx, cmd)
	}
}

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

	// Setup mailbox
	replyC, remove := c.rpcMailbox.Create(rpcID)
	defer remove()

	// Write
	select {
	case c.writeC <- b:
	case <-c.doneC:
		return ErrClosed
	case <-ctx.Done():
		return ctx.Err()
	}

	// Wait for reply in mailbox
	return rpcWait(ctx, replyC)
}

func (c connection) commandRPCNoReply(ctx context.Context, cmd android.Command) error {
	rpc := protos.Rpc{}

	err := convertCommand(cmd, &rpc)
	if err != nil {
		return err
	}

	// Marshal
	b, err := proto.Marshal(&rpc)
	if err != nil {
		return fmt.Errorf("failed to marshal command: %w", err)
	}

	// Write
	select {
	case c.writeC <- b:
	case <-c.doneC:
		return ErrClosed
	case <-ctx.Done():
		return ctx.Err()
	}

	return nil
}
