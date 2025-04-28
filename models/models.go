package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Base struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt,omitempty" gorm:"index"`
}

type Conversation struct {
	Base
	Type            string    `json:"type" gorm:"type:conversation_type;default:'direct'"` // direct, team_channel
	LastMessageAt   time.Time `json:"lastMessageAt" gorm:"index"`
	LastMessageText string    `json:"lastMessageText"` // Cached last message for preview
	LastMessageBy   uuid.UUID `json:"lastMessageBy" gorm:"index"`
	UnreadCount     int       `json:"unreadCount" gorm:"default:0"` // Cached unread count
	Status          string    `json:"status" gorm:"type:conversation_status;default:'active'"`

	HasMessages bool `gorm:"default:false"` // New field

	// For direct messages
	User1ID uuid.UUID `json:"user1Id" gorm:"index;default:null"` // First participant
	User2ID uuid.UUID `json:"user2Id" gorm:"index;default:null"` // Second participant

	// For team channels
	TeamID      uuid.UUID `json:"teamId" gorm:"index;default:null"`
	TeamName    string    `json:"teamName"`                                  // Cached team name for quick access
	ChannelID   uuid.UUID `json:"channelId" gorm:"uniqueIndex;default:null"` // Make it uniqueIndex
	ChannelName string    `json:"channelName"`                               // Cached channel name for quick access
}

type TeamChannel struct {
	Base
	TeamID      uuid.UUID `json:"teamId" gorm:"index"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description"`
	Type        string    `json:"type" gorm:"type:team_channel_type;default:'general'"`
	CreatedBy   uuid.UUID `json:"createdBy"`
	Position    int       `json:"position" gorm:"default:0"`
	IsDefault   bool      `json:"isDefault" gorm:"default:false"`
	IsArchived  bool      `json:"isArchived" gorm:"default:false"`

	ReadOnly     bool     `json:"readOnly" gorm:"default:false"`
	AllowedRoles []string `json:"allowedRoles" gorm:"type:text[]"`

	ConversationID uuid.UUID     `json:"conversationId" gorm:"uniqueIndex;not null"`
	Conversation   *Conversation `json:"-" gorm:"constraint:OnDelete:CASCADE;"`
}

type User struct {
	Base
	Email           string    `json:"email" gorm:"uniqueIndex;not null"`
	PasswordHash    string    `json:"-"` // "-" means this field won't be included in JSON
	FullName        string    `json:"fullName"`
	DisplayName     string    `json:"displayName"`
	ProfileImage    string    `json:"profileImage"`
	Bio             string    `json:"bio" gorm:"type:text"`
	JobTitle        string    `json:"jobTitle"`
	Company         string    `json:"company"`
	Location        string    `json:"location"`
	Website         string    `json:"website"`
	Status          string    `json:"status" gorm:"type:user_status;default:'offline'"`
	LastActive      time.Time `json:"lastActive"`
	ProfileComplete bool      `json:"profileComplete" gorm:"default:false"`
	Provider        string    `json:"provider" gorm:"type:user_auth_provider"`

	// Settings interface{} `json:"settings" gorm:"foreignKey:UserID"`

	// Teams             []interface{}
	// Meetings          []interface{}
	// OwnedTeams        []interface{} `gorm:"foreignKey:OwnerID"`
	// OrganizedMeetings []interface{} `gorm:"foreignKey:OrganizerID"`
}
