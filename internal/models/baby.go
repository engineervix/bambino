package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Baby struct {
	ID          uuid.UUID `gorm:"type:varchar(36);primary_key"`
	UserID      uuid.UUID `gorm:"type:varchar(36);not null;index"`
	Name        string    `gorm:"type:varchar(100);not null"`
	BirthDate   time.Time `gorm:"type:date;not null"`
	BirthWeight *float64  `gorm:"type:decimal(5,2)"`
	BirthHeight *float64  `gorm:"type:decimal(5,2)"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	User        User       `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Activities  []Activity `gorm:"foreignKey:BabyID;constraint:OnDelete:CASCADE"`
}

func (b *Baby) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}

// BeforeSave hook to validate required fields
func (b *Baby) BeforeSave(tx *gorm.DB) error {
	if b.UserID == uuid.Nil {
		return gorm.ErrInvalidField
	}
	if b.Name == "" {
		return gorm.ErrInvalidField
	}
	return nil
}
