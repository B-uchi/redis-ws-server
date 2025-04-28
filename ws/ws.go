package ws

import (
	"encoding/json"
	"huddle-ws-server/rd"
	"huddle-ws-server/types"
	"sync"

	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
)

type Client struct {
	Connection     *websocket.Conn
	UserId         uuid.UUID
	Channels       map[uuid.UUID]bool
	DirectMessages map[uuid.UUID]bool
}

type Manager struct {
	clients    map[*Client]bool
	userConns  map[uuid.UUID][]*Client
	broadcast  chan types.Message
	register   chan *Client
	unregister chan *Client
	mutex      sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{
		clients:    make(map[*Client]bool),
		userConns:  make(map[uuid.UUID][]*Client),
		broadcast:  make(chan types.Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

var WsManager = NewManager()

func (manager *Manager) Start() {
	for {
		select {
		case client := <-manager.register:
			manager.mutex.Lock()
			manager.clients[client] = true
			manager.userConns[client.UserId] = append(manager.userConns[client.UserId], client)
			manager.mutex.Unlock()

			statusPayload, _ := json.Marshal(map[string]interface{}{
				"userId": client.UserId.String(),
				"status": "online",
			})

			// Broadcast online user
			rd.Publish("user_online_status", statusPayload)

		case client := <-manager.unregister:
			manager.mutex.Lock()
			if _, ok := manager.clients[client]; ok {
				delete(manager.clients, client)
				manager.removeUserConnections(client)
				client.Connection.Close()

				if activeUserConnections, exists := manager.userConns[client.UserId]; !exists || len(activeUserConnections) == 0 {
					statusPayload, _ := json.Marshal(map[string]interface{}{
						"userId": client.UserId.String(),
						"status": "offline",
					})
					rd.Publish("user_online_status", statusPayload)

				}

			}
			manager.mutex.Unlock()
		case message := <-manager.broadcast:
			manager.BroadcastMessage(message)
		}

	}
}

func (manager *Manager) removeUserConnections(client *Client) {
	if connections, ok := manager.userConns[client.UserId]; ok {
		newConnections := make([]*Client, 0)
		for _, connection := range connections {
			if connection != client {
				newConnections = append(newConnections, connection)
			}
		}
		if len(newConnections) == 0 {
			delete(manager.userConns, client.UserId)
		} else {
			manager.userConns[client.UserId] = newConnections
		}
	}
}

func (manager *Manager) BroadcastUserStatus(userID uuid.UUID, status string) {
	manager.mutex.RLock()
	defer manager.mutex.RUnlock()

	statusUpdate := types.Message{
		Type: "user_status",
		Data: map[string]interface{}{
			"userId": userID,
			"status": status,
		},
	}

	// Broadcast to all connected clients
	for client := range manager.clients {
		err := client.Connection.WriteJSON(statusUpdate)
		if err != nil {
			go func(c *Client) {
				manager.unregister <- c
			}(client)
		}
	}
}

func (manager *Manager) SubscribeToChannel(client *Client, channelId uuid.UUID) {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()
	client.Channels[channelId] = true
}

func (manager *Manager) SubscribeToConversation(client *Client, conversationID uuid.UUID) {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()
	client.DirectMessages[conversationID] = true
}

func (manager *Manager) BroadcastReaction(channelID *uuid.UUID, messageID string, reaction types.MessageReactionResponse, action string) {
	manager.broadcast <- types.Message{
		Type:      "reaction",
		ChannelID: channelID,
		Reaction: &types.MessageReactionEvent{
			MessageID: messageID,
			Reaction:  reaction,
			Action:    action,
		},
	}
}

func (m *Manager) BroadcastMessage(msg types.Message) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	for client := range m.clients {
		if client.UserId.String() == msg.Message.SenderID {
			continue
		}

		// For channel messages
		if msg.ChannelID != nil {
			if subscribed := client.Channels[*msg.ChannelID]; !subscribed {
				continue
			}
		}

		// For direct messages
		if msg.ConversationID != nil {
			if subscribed := client.DirectMessages[*msg.ConversationID]; !subscribed {
				continue
			}
		}

		err := client.Connection.WriteJSON(msg)
		if err != nil {
			go func(c *Client) {
				m.unregister <- c
			}(client)
		}
	}
}
