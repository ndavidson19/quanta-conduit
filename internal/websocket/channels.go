package websocket

import (
	"sync"
)

type ChannelManager struct {
	channels map[string]map[*AuthenticatedClient]bool
	mu       sync.RWMutex
}

func NewChannelManager() *ChannelManager {
	return &ChannelManager{
		channels: make(map[string]map[*AuthenticatedClient]bool),
	}
}

func (cm *ChannelManager) Subscribe(client *AuthenticatedClient, channel string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if _, ok := cm.channels[channel]; !ok {
		cm.channels[channel] = make(map[*AuthenticatedClient]bool)
	}
	cm.channels[channel][client] = true
}

func (cm *ChannelManager) Unsubscribe(client *AuthenticatedClient, channel string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if clients, ok := cm.channels[channel]; ok {
		delete(clients, client)
		if len(clients) == 0 {
			delete(cm.channels, channel)
		}
	}
}

func (cm *ChannelManager) Broadcast(channel string, message []byte) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	if clients, ok := cm.channels[channel]; ok {
		for client := range clients {
			select {
			case client.send <- message:
			default:
				close(client.send)
				delete(clients, client)
			}
		}
	}
}

// Add this method to the Manager struct in manager.go
func (m *Manager) SubscribeClientToChannel(client *AuthenticatedClient, channel string) {
	m.channelManager.Subscribe(client, channel)
	m.logger.Info("Client subscribed to channel", zap.Uint("userID", client.UserID), zap.String("channel", channel))
}