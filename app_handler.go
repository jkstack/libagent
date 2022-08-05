package agent

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jkstack/agent/internal/utils"
	"github.com/jkstack/anet"
	"github.com/jkstack/jkframe/logging"
)

func (app *app) read(ctx context.Context, cancel context.CancelFunc, conn *websocket.Conn) {
	defer utils.Recover("read")
	defer cancel()
	defer conn.Close()
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		_, data, err := conn.ReadMessage()
		if err != nil {
			logging.Error("read message: %v", err)
			return
		}

		app.incInPackets()
		app.incInBytes(uint64(len(data)))

		var msg anet.Msg
		err = json.Unmarshal(data, &msg)
		if err != nil {
			logging.Error("json unmarshal: %v", err)
			return
		}

		app.chRead <- &msg
	}
}

func (app *app) write(ctx context.Context, cancel context.CancelFunc, conn *websocket.Conn) {
	defer utils.Recover("write")
	defer cancel()
	defer conn.Close()

	for {
		select {
		case msg := <-app.chWrite:
			data, err := json.Marshal(msg)
			if err != nil {
				logging.Error("json marshal: %v", err)
				return
			}

			err = conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				logging.Error("write message: %v", err)
				return
			}

			app.incOutPackets()
			app.incOutBytes(uint64(len(data)))
		case <-ctx.Done():
			return
		}
	}
}

func (app *app) readCallback() {
	for {
		app.a.OnMessage(<-app.chRead)
	}
}

func (app *app) keepalive(ctx context.Context, conn *websocket.Conn) {
	tk := time.NewTicker(10 * time.Second)
	defer tk.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-tk.C:
			conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(2*time.Second))
		}
	}
}
