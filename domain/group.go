package domain

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Group struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Name        string    `gorm:"not null"`
	Image       string    `gorm:"not null;default:''"`
	AdminID     uuid.UUID `gorm:"not null;index"`
	Admin       User      `gorm:"foreignKey:AdminID"`
	LastMessage string    `gorm:"not null;default:''"`
	Status      bool      `gorm:"not null;default:true"`
}

type GroupListResponse struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	Image          string    `json:"image"`
	UnreadMessages int       `json:"unread_messages"`
	LastMessage    string    `json:"last_message"`
}

func (group *Group) ToGroupListResponse() GroupListResponse {

	return GroupListResponse{
		ID:             group.ID,
		Name:           group.Name,
		Image:          group.Image,
		UnreadMessages: 0,
		LastMessage:    "",
	}
}

type GroupDetailsResponse struct {
	Name  string             `json:"name"`
	Image string             `json:"image"`
	Users []UserListResponse `json:"users"`
}
