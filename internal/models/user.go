package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID           uuid.UUID `gorm:"type:varchar(36);primary_key"`
	Username     string    `gorm:"type:varchar(50);unique;not null"`
	PasswordHash string    `gorm:"type:varchar(255);not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Babies       []Baby `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

// BeforeSave hook to validate required fields
func (u *User) BeforeSave(tx *gorm.DB) error {
	if u.Username == "" {
		return gorm.ErrInvalidField
	}
	if u.PasswordHash == "" {
		return gorm.ErrInvalidField
	}
	return nil
}
