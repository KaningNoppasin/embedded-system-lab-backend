package services

import (
	"encoding/json"
	"sync"

	"github.com/KaningNoppasin/embedded-system-lab-backend/app/mqtt"
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

	return h.broadcast(payload)
}

func (h *RFIDWebSocketHub) BroadcastTemperature(topic string, temperature float64) error {
	payload, err := json.Marshal(map[string]any{
		"type":        "temperature",
		"topic":       topic,
		"temperature": temperature,
	})
	if err != nil {
		return err
	}

	return h.broadcast(payload)
}

func (h *RFIDWebSocketHub) BroadcastINA219(topic string, ina219Payload mqtt.INA219Payload) error {
	payload, err := json.Marshal(map[string]any{
		"type":           "ina219",
		"topic":          topic,
		"ina219_payload": ina219Payload,
	})
	if err != nil {
		return err
	}

	return h.broadcast(payload)
}

func (h *RFIDWebSocketHub) broadcast(payload []byte) error {
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
