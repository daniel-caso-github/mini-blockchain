package api

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/danielcaso/mini-blockchain/internal/blockchain"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Hub struct {
	clients    map[*websocket.Conn]bool
	broadcast  chan []byte
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
	mu         sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*websocket.Conn]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case conn := <-h.register:
			h.mu.Lock()
			h.clients[conn] = true
			h.mu.Unlock()

		case conn := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[conn]; ok {
				delete(h.clients, conn)
				conn.Close()
			}
			h.mu.Unlock()

		case msg := <-h.broadcast:
			h.mu.Lock()
			for conn := range h.clients {
				if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
					conn.Close()
					delete(h.clients, conn)
				}
			}
			h.mu.Unlock()
		}
	}
}

func (h *Hub) HandleWS(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("websocket upgrade error: %v", err)
		return
	}

	h.register <- conn

	// Read loop to detect client disconnect
	go func() {
		defer func() {
			h.unregister <- conn
		}()
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				break
			}
		}
	}()
}

func (h *Hub) BroadcastBlock(block blockchain.Block) {
	msg := WSMessage{
		Type:  "new_block",
		Block: block,
	}
	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("websocket marshal error: %v", err)
		return
	}
	h.broadcast <- data
}

func (h *Hub) BroadcastDifficulty(newDifficulty int) {
	msg := WSDifficultyMessage{
		Type:       "difficulty_adjusted",
		Difficulty: newDifficulty,
	}
	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("websocket marshal error: %v", err)
		return
	}
	h.broadcast <- data
}
