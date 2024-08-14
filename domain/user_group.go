package domain

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserGroup struct {
	gorm.Model
	ID             uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	UserID         uuid.UUID `gorm:"not null;index"`
	User           User      `gorm:"foreignKey:UserID"`
	GroupID        uuid.UUID `gorm:"not null;index"`
	Group          Group     `gorm:"foreignKey:GroupID"`
	UnreadMessages int       `gorm:"not null;default:0"`
	Notification   bool      `gorm:"not null;default:true"`
	Status         bool      `gorm:"not null;default:true"`
}
