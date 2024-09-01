package websocket

import (
	"encoding/json"
	"net/http"

	"conduit/internal/auth"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type AuthenticatedClient struct {
	Client
	UserID uint
}

func (m *Manager) AuthenticateAndServeWS(w http.ResponseWriter, r *http.Request, logger *zap.Logger) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Missing authentication token", http.StatusUnauthorized)
		return
	}

	claims, err := auth.ValidateToken(token)
	if err != nil {
		logger.Error("Invalid token", zap.Error(err))
		http.Error(w, "Invalid authentication token", http.StatusUnauthorized)
		return
	}

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // You might want to implement a more secure check
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error("Failed to upgrade to WebSocket", zap.Error(err))
		return
	}

	client := &AuthenticatedClient{
		Client: Client{conn: conn, send: make(chan []byte, 256)},
		UserID: claims.UserID,
	}

	m.register <- &client.Client
	go client.writePump()
	go client.readPump(m)
}

func (c *AuthenticatedClient) readPump(m *Manager) {
	defer func() {
		m.unregister <- &c.Client
		c.conn.Close()
	}()
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				m.logger.Error("unexpected close error", zap.Error(err))
			}
			break
		}

		// Process incoming messages (e.g., subscribe to specific channels)
		var msg map[string]interface{}
		if err := json.Unmarshal(message, &msg); err != nil {
			m.logger.Error("Failed to unmarshal message", zap.Error(err))
			continue
		}

		// Handle subscription requests
		if action, ok := msg["action"].(string); ok && action == "subscribe" {
			if channel, ok := msg["channel"].(string); ok {
				m.SubscribeClientToChannel(c, channel)
			}
		}
	}
}
