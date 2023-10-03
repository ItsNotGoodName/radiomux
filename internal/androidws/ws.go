package androidws

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

func wsWriter(ctx context.Context, cancel context.CancelFunc, conn *websocket.Conn, log zerolog.Logger, writeC chan []byte) {
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
		case data := <-writeC:
			conn.SetWriteDeadline(time.Now().Add(wsWriteWait))

			// Send data or end on error
			if err := conn.WriteMessage(websocket.BinaryMessage, data); err != nil {
				log.Err(err).Msg("Failed to write to WebSocket")
				return
			}
		}
	}
}
