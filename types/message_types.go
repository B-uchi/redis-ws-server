package types

import (
	"github.com/google/uuid"
)

type MessageResponse struct {
	ID               string               `json:"id"`
	Content          string               `json:"content"`
	ContentType      string               `json:"contentType"`
	ReplyToMessageID string               `json:"replyToMessageId,omitempty"`
	IsEdited         bool                 `json:"isEdited,omitempty"`
	FilePath         string               `json:"filePath,omitempty"`
	SenderID         string               `json:"senderId"`
	SenderName       string               `json:"senderName"`
	SenderAvatar     string               `json:"senderAvatar,omitempty"`
	CreatedAt        string               `json:"createdAt"`
	ChannelID        string               `json:"channelId,omitempty"`
	ConversationID   string               `json:"conversationId,omitempty"`
	IsMe             bool                 `json:"isMe"`
	Reactions        []ReactionResponse   `json:"reactions,omitempty"`
	Attachments      []AttachmentResponse `json:"attachments,omitempty"`
}

type AttachmentResponse struct {
	ID          string `json:"id"`
	FileName    string `json:"fileName"`
	FileSize    int64  `json:"fileSize"`
	FileType    string `json:"fileType"`
	ContentType string `json:"contentType"`
	URL         string `json:"url"`
}

type ReactionResponse struct {
	Emoji      string     `json:"emoji"`
	Count      int        `json:"count"`
	Users      []UserInfo `json:"users"`
	HasReacted bool       `json:"hasReacted"`
}

type UserInfo struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar,omitempty"`
}

type MessageUpdateData struct {
	MessageID string `json:"messageId"`
	Content   string `json:"content"`
	IsEdited  bool   `json:"isEdited"`
	UserId    string `json:"userId"`
	Status    string `json:"status"`
}

type Message struct {
	Type           string                `json:"type"`
	ConversationID *uuid.UUID            `json:"conversationId,omitempty"`
	ChannelID      *uuid.UUID            `json:"channelId,omitempty"`
	Message        MessageResponse       `json:"message,omitempty"`
	Event          string                `json:"event,omitempty"`
	Reaction       *MessageReactionEvent `json:"reaction,omitempty"`
	Data           interface{}           `json:"data,omitempty"`
}

type MessageReactionEvent struct {
	MessageID string                  `json:"messageId"`
	Reaction  MessageReactionResponse `json:"reaction"`
	Action    string                  `json:"action"` // "add" or "remove"
}

type MessageReactionResponse struct {
	ID              string `json:"id"`
	Emoji           string `json:"emoji"`
	MessageID       string `json:"messageId"`
	UserId          string `json:"userId"`
	UserAvatar      string `json:"userAvatar"`
	UserDisplayName string `json:"userDisplayName"`
}

type TypingPayload struct {
	Type            string     `json:"type"`
	ConversationId  *uuid.UUID `json:"conversationId"`
	ChannelId       *uuid.UUID `json:"channelId"`
	UserID          string     `json:"userId"`
	UserAvatar      *string    `json:"userAvatar"`
	UserDisplayName *string    `json:"userDisplayName"`
}
