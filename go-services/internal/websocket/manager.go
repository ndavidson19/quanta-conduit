package websocket

import (
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type Manager struct {
	clients        map[*Client]bool
	broadcast      chan []byte
	register       chan *Client
	unregister     chan *Client
	logger         *zap.Logger
	channelManager *ChannelManager
}

type Client struct {
	conn *websocket.Conn
	send chan []byte
}

func NewManager(logger *zap.Logger) *Manager {
	return &Manager{
		clients:        make(map[*Client]bool),
		broadcast:      make(chan []byte),
		register:       make(chan *Client),
		unregister:     make(chan *Client),
		logger:         logger,
		channelManager: NewChannelManager(),
	}
}

func (m *Manager) Run() {
	for {
		select {
		case client := <-m.register:
			m.clients[client] = true
		case client := <-m.unregister:
			if _, ok := m.clients[client]; ok {
				delete(m.clients, client)
				close(client.send)
			}
		case message := <-m.broadcast:
			for client := range m.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(m.clients, client)
				}
			}
		}
	}
}

func (m *Manager) ServeWS(conn *websocket.Conn) {
	client := &Client{conn: conn, send: make(chan []byte, 256)}
	m.register <- client

	go client.writePump()
	go client.readPump(m)
}

func (c *Client) writePump() {
	defer func() {
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)
			if err := w.Close(); err != nil {
				return
			}
		}
	}
}

func (c *Client) readPump(m *Manager) {
	defer func() {
		m.unregister <- c
		c.conn.Close()
	}()
	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				m.logger.Error("unexpected close error", zap.Error(err))
			}
			break
		}
	}
}

func (m *Manager) BroadcastMessage(message []byte) {
	m.broadcast <- message
}
