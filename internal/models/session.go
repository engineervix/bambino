package models

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID        string    `gorm:"type:varchar(255);primary_key"`
	UserID    uuid.UUID `gorm:"type:uuid;not null"`
	ExpiresAt time.Time `gorm:"not null;index"`
	CreatedAt time.Time
	User      User `gorm:"foreignKey:UserID"`
}
