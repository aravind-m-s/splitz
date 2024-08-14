package domain

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRequest struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	RequestID   uuid.UUID `gorm:"not null;inde"`
	Request     Request   `gorm:"foreignKey:RequestID"`
	Share       float64   `gorm:"not null"`
	Paid        float64   `gorm:"not null;default:0"`
	LastMessage string    `gorm:"not null;default:''"`
	Status      bool      `gorm:"not null;default:true"`
}
