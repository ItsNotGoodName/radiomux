package main

import (
	"context"
	"errors"
	"time"

	"github.com/ItsNotGoodName/radiomux/internal/android"
	"github.com/ItsNotGoodName/radiomux/internal/apiws"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"nhooyr.io/websocket"
)

type Websocket struct {
	stateService *android.StateService
}

func (s Websocket) Handle(c echo.Context) error {
	// Setup websocket
	w := c.Response()
	r := c.Request()
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})
	if err != nil {
		return err
	}
	defer conn.Close(websocket.StatusInternalError, "the sky is falling")

	ctx := r.Context()

	stateC, stateCDisconnect := s.stateService.Subscribe()
	defer stateCDisconnect()

	playerState := apiws.NewPlayerState(s.stateService)
	client := apiws.NewClient(playerState)

	syncC := make(chan struct{}, 1)
	syncC <- struct{}{}
	flushC := make(chan chan []byte)
	defer close(syncC)
	go websocketWriter(ctx, conn, syncC, flushC)

	readC := make(chan []byte)
	go websocketReader(ctx, conn, readC)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case data, ok := <-readC:
			// Reader
			if !ok {
				return nil
			}

			log.Debug().Bytes("data", data).Msg("READ")
		case dataC := <-flushC:
			data, err := client.Flush()
			if err != nil {
				close(dataC)
				if errors.Is(err, apiws.ErrVisitorEmpty) {
					continue
				}
				log.Err(err).Msg("Flush failed")
				return nil
			}

			select {
			case <-ctx.Done():
				close(dataC)
				return ctx.Err()
			case dataC <- data:
			}

			if !playerState.Empty {
				select {
				case syncC <- struct{}{}:
				default:
				}
			}
		case now, ok := <-stateC:
			// State
			if !ok {
				return nil
			}

			playerState.StateChange(now)

			select {
			case syncC <- struct{}{}:
			default:
			}
		}
	}
}

func websocketWriter(ctx context.Context, conn *websocket.Conn, syncC <-chan struct{}, flushC chan chan []byte) {
	for {
		select {
		case _, ok := <-syncC:
			if !ok {
				return
			}

			dataC := make(chan []byte)
			log.Debug().Msg("LOOK START")
			select {
			case flushC <- dataC:
			case <-ctx.Done():
				return
			}
			log.Debug().Msg("LOOK END")

			data, ok := <-dataC
			if !ok {
				continue
			}

			// NOTE: This is emulating a slow websocket connection, it should be removed
			time.Sleep(0 * time.Second)

			err := conn.Write(ctx, websocket.MessageText, data)
			if err != nil {
				log.Err(err).Send()
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

func websocketReader(ctx context.Context, conn *websocket.Conn, readC chan []byte) {
	for {
		_, data, err := conn.Read(ctx)
		if err != nil {
			return
		}

		select {
		case <-ctx.Done():
			return
		case readC <- data:
		}
	}
}
