// services/gateway/internal/ws/hub.go
package ws

import (
	"sync"

	"github.com/gofiber/contrib/websocket"
)

type Hub struct {
	mu      sync.RWMutex
	clients map[*websocket.Conn]bool
}

func NewHub() *Hub { return &Hub{clients: make(map[*websocket.Conn]bool)} }

func (h *Hub) Register(c *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[c] = true
}

func (h *Hub) Unregister(c *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.clients, c)
}

func (h *Hub) Broadcast(message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for c := range h.clients {
		if err := c.WriteMessage(websocket.TextMessage, message); err != nil {
			go h.Unregister(c)
		}
	}
}