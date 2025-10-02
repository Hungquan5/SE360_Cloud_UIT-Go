package ws

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Hub struct {
	mu    sync.RWMutex
	rooms map[string]map[*websocket.Conn]struct{}
}

func NewHub() *Hub {
	return &Hub{rooms: make(map[string]map[*websocket.Conn]struct{})}
}

func (h *Hub) Add(room string, c *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	m, ok := h.rooms[room]
	if !ok {
		m = make(map[*websocket.Conn]struct{})
		h.rooms[room] = m
	}
	m[c] = struct{}{}
}

func (h *Hub) Remove(room string, c *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if m, ok := h.rooms[room]; ok {
		delete(m, c)
		if len(m) == 0 {
			delete(h.rooms, room)
		}
	}
}

func (h *Hub) Broadcast(room string, payload []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if m, ok := h.rooms[room]; ok {
		for c := range m {
			_ = c.WriteMessage(websocket.TextMessage, payload)
		}
	}
}

func (h *Hub) BroadcastAll(payload []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for _, m := range h.rooms {
		for c := range m {
			_ = c.WriteMessage(websocket.TextMessage, payload)
		}
	}
}
