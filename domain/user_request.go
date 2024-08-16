package domain

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRequest struct {
	gorm.Model
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	RequestID uuid.UUID `gorm:"not null;index"`
	Request   Request   `gorm:"foreignKey:RequestID"`
	UserID    uuid.UUID `gorm:"not null;index"`
	User      User      `gorm:"foreignKey:UserID"`
	Share     float64   `gorm:"not null"`
	Paid      float64   `gorm:"not null;default:0"`
	Status    bool      `gorm:"not null;default:true"`
}
