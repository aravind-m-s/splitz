package domain

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Request struct {
	gorm.Model
	ID      uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Note    string    `gorm:"not null"`
	Amount  float64   `gorm:"not null"`
	OwnerID uuid.UUID `gorm:"not null;index"`
	Owner   User      `gorm:"foreignKey:OwnerID"`
	GroupID uuid.UUID `gorm:"not null;index"`
	Group   Group     `gorm:"foreignKey:GroupID"`
	Type    string    `gorm:"not null"`
	Status  bool      `gorm:"not null"`
}
