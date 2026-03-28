package services

import (
	"encoding/json"
	"sync"

	"github.com/fasthttp/websocket"
)

type RFIDWebSocketHub struct {
	mu      sync.RWMutex
	clients map[*websocket.Conn]struct{}
}

func NewRFIDWebSocketHub() *RFIDWebSocketHub {
	return &RFIDWebSocketHub{
		clients: make(map[*websocket.Conn]struct{}),
	}
}

func (h *RFIDWebSocketHub) AddClient(conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.clients[conn] = struct{}{}
}

func (h *RFIDWebSocketHub) RemoveClient(conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	delete(h.clients, conn)
}

func (h *RFIDWebSocketHub) BroadcastRFID(rfid string) error {
	payload, err := json.Marshal(map[string]string{
		"rfid": rfid,
	})
	if err != nil {
		return err
	}

	h.mu.RLock()
	clients := make([]*websocket.Conn, 0, len(h.clients))
	for conn := range h.clients {
		clients = append(clients, conn)
	}
	h.mu.RUnlock()

	for _, conn := range clients {
		if err := conn.WriteMessage(websocket.TextMessage, payload); err != nil {
			h.RemoveClient(conn)
			_ = conn.Close()
		}
	}

	return nil
}
