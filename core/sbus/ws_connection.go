package sbus

import (
	"encoding/hex"
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/wwengg/threego/core/slog"
)

type WsConnection struct {
	Connection
	// conn is the current connection's WebSocket socket TCP socket. (当前连接的socket TCP套接字)
	conn *websocket.Conn
}

func NewWsConnection(cID uint64, taskHandler STaskHandler, conn *websocket.Conn, onConnStart, onConnStop func(conn SConnection), datapack SDataPack) SConnection {
	return &WsConnection{
		Connection: Connection{
			ConnID:         cID,
			ConnIdStr:      fmt.Sprintf("%d", cID),
			TaskHandler:    taskHandler,
			OnConnStart:    onConnStart,
			OnConnStop:     onConnStop,
			Datapack:       datapack,
			Property:       nil,
			IOReadBuffSize: 0,
		},
		conn: conn,
	}
}

func wsSendFunc(wsConn *websocket.Conn) func([]byte) error {
	return func(data []byte) error {
		return wsConn.WriteMessage(websocket.BinaryMessage, data)
	}
}

func wsReadFunc(wsConn *websocket.Conn) func(conn SConnection, buffer []byte) (n int, err error) {
	return func(conn SConnection, buffer []byte) (n int, err error) {
		if len(buffer) != 0 {
			buffer = buffer[len(buffer):]
		}
		messageType, buffer, err := wsConn.ReadMessage()
		if err != nil {
			conn.Stop()
			return 0, fmt.Errorf("websocket Read Error: %s", err.Error())
		}
		if messageType == websocket.PingMessage {
			//conn.updateActivity()
			return 0, nil
		}
		n = len(buffer)
		if err != nil {
			return 0, fmt.Errorf("read msg head [read datalen=%d], error = %s", n, err.Error())
		}
		slog.Ins().Debugf("read buffer %s \n", hex.EncodeToString(buffer[0:n]))

		return n, nil
	}
}
