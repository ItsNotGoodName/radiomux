package apiws

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
)

const (
	// Time allowed to write a message to the peer.
	wsWriteWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	wsPongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	wsPingPeriod = (wsPongWait * 9) / 10

	// Maximum message size allowed from peer.
	wsMaxMessageSize = 512
)

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func wsUpgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	return wsUpgrader.Upgrade(w, r, nil)
}

func wsReader(ctx context.Context, cancel context.CancelFunc, conn *websocket.Conn, log zerolog.Logger, readC chan []byte) {
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

		// Send data to handler
		select {
		case readC <- data:
		case <-ctx.Done():
			return
		}
	}
}

func wsWriter(ctx context.Context, cancel context.CancelFunc, conn *websocket.Conn, log zerolog.Logger, sync signal, flushC chan chan []byte) {
	defer cancel()

	ticker := time.NewTicker(wsPingPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			conn.SetWriteDeadline(time.Now().Add(wsWriteWait))

			// Send ping or end on error
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Err(err).Msg("Failed to write ping")
				return
			}
		case <-sync.C:
			dataC := make(chan []byte)
			select {
			case <-ctx.Done():
				return
			case flushC <- dataC:
			}

			data, ok := <-dataC
			if !ok {
				continue
			}

			conn.SetWriteDeadline(time.Now().Add(wsWriteWait))

			// Send data or end on error
			if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Err(err).Msg("Failed to write to WebSocket")
				return
			}
		}
	}
}
