package domain

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID        uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Mobile    string         `gorm:"uniqueIndex;not null"`
	Password  string         `gorm:"not null"`
	Name      string         `gorm:"not null"`
	Image     string         `gorm:"not null"`
	FcmTokens pq.StringArray `gorm:"type:text[]"`
}

type UserResponse struct {
	ID     uuid.UUID `json:"id"`
	Mobile string    `json:"mobile"`
	Name   string    `json:"name"`
	Image  string    `json:"image"`
	Token  string    `json:"token"`
}

func (u *User) ToUserResponse() (user UserResponse) {
	return UserResponse{
		ID:     u.ID,
		Mobile: u.Mobile,
		Name:   u.Name,
		Image:  u.Image,
	}
}

type UserListResponse struct {
	ID      uuid.UUID `json:"id"`
	Mobile  string    `json:"mobile"`
	Name    string    `json:"name"`
	Image   string    `json:"image"`
	IsAdmin bool      `json:"is_admin"`
}

func (u *User) ToUserListResponse() (user UserListResponse) {
	return UserListResponse{
		ID:     u.ID,
		Mobile: u.Mobile,
		Name:   u.Name,
		Image:  u.Image,
	}
}

func (u *User) ToGroupDetailsUser(isAdmin bool) (user UserListResponse) {
	return UserListResponse{
		ID:      u.ID,
		Mobile:  u.Mobile,
		Name:    u.Name,
		Image:   u.Image,
		IsAdmin: isAdmin,
	}
}
