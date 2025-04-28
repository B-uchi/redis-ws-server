package ws

import (
	"encoding/json"
	"huddle-ws-server/database"
	"huddle-ws-server/models"
	"huddle-ws-server/types"
	"log"
	"time"

	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
)

func WebsocketHandler(c *websocket.Conn) {
	userID := c.Locals("userID").(uuid.UUID)
	log.Printf("New WebSocket connection from user: %s", userID)

	// Create new client
	client := &Client{
		Connection:     c,
		UserId:         userID,
		Channels:       make(map[uuid.UUID]bool),
		DirectMessages: make(map[uuid.UUID]bool),
	}

	// Register client
	WsManager.register <- client

	// Subscribe to user's channels and conversations
	subscribeToUserChannels(client)
	subscribeToUserConversations(client)

	// Set up a ping handler to detect disconnections
	c.SetPingHandler(func(string) error {
		return c.WriteControl(websocket.PongMessage, []byte{}, time.Now().Add(time.Second))
	})

	// Create a done channel to signal goroutine cleanup
	done := make(chan struct{})
	defer close(done)

	// Start a goroutine to handle ping/pong
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := c.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(time.Second)); err != nil {
					log.Printf("ping error: %v", err)
					WsManager.unregister <- client
					return
				}
			case <-done:
				return
			}
		}
	}()

	defer func() {
		log.Printf("WebSocket connection closed for user: %s", userID)
		WsManager.unregister <- client
	}()

	for {
		messageType, msg, err := c.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error for user %s: %v", userID, err)
			}
			break
		}

		if messageType == websocket.TextMessage {
			var payload interface{}
			if err := json.Unmarshal(msg, &payload); err != nil {
				continue
			}
			handleIncomingMessage(payload, userID)
		}
	}
}

func handleIncomingMessage(payload interface{}, userID uuid.UUID) {
	// First, determine the type of message
	var messageType struct {
		Type string `json:"type"`
	}

	payloadBytes, _ := json.Marshal(payload)
	if err := json.Unmarshal(payloadBytes, &messageType); err != nil {
		return
	}

	switch messageType.Type {
	case "typing", "stop_typing":
		var typingPayload types.TypingPayload
		if err := json.Unmarshal(payloadBytes, &typingPayload); err != nil {
			return
		}

		var user models.User
		if err := database.DB.Where("id = ?", userID).First(&user).Error; err != nil {
			return
		}

		typingPayload.UserAvatar = &user.ProfileImage
		typingPayload.UserDisplayName = &user.DisplayName

		var message types.Message

		message.ChannelID = typingPayload.ChannelId
		message.ConversationID = typingPayload.ConversationId
		message.Message.SenderID = userID.String()
		message.Type = typingPayload.Type
		message.Data = typingPayload

		WsManager.BroadcastMessage(message)
	}
}

func subscribeToUserChannels(client *Client) {
	var channels []models.TeamChannel
	if err := database.DB.
		Joins("JOIN team_members ON team_members.team_id = team_channels.team_id").
		Where("team_members.user_id = ?", client.UserId).
		Find(&channels).Error; err != nil {
		return
	}

	for _, channel := range channels {
		WsManager.SubscribeToChannel(client, channel.ID)
	}
}

func subscribeToUserConversations(client *Client) {
	var conversations []models.Conversation
	if err := database.DB.
		Where("user1_id = ? OR user2_id = ?", client.UserId, client.UserId).
		Find(&conversations).Error; err != nil {
		return
	}

	for _, conv := range conversations {
		WsManager.SubscribeToConversation(client, conv.ID)
	}
}
