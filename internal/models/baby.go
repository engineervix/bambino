package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Baby struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID      uuid.UUID `gorm:"type:uuid;not null"`
	Name        string    `gorm:"type:varchar(100);not null"`
	BirthDate   time.Time `gorm:"type:date;not null"`
	BirthWeight *float64  `gorm:"type:decimal(5,2)"`
	BirthHeight *float64  `gorm:"type:decimal(5,2)"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	User        User       `gorm:"foreignKey:UserID"`
	Activities  []Activity `gorm:"foreignKey:BabyID"`
}

func (b *Baby) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}
