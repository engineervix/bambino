package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ActivityType string

const (
	ActivityTypeFeed      ActivityType = "feed"
	ActivityTypePump      ActivityType = "pump"
	ActivityTypeDiaper    ActivityType = "diaper"
	ActivityTypeSleep     ActivityType = "sleep"
	ActivityTypeGrowth    ActivityType = "growth"
	ActivityTypeHealth    ActivityType = "health"
	ActivityTypeMilestone ActivityType = "milestone"
)

type Activity struct {
	ID        uuid.UUID    `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	BabyID    uuid.UUID    `gorm:"type:uuid;not null"`
	Type      ActivityType `gorm:"type:varchar(20);not null"`
	StartTime time.Time    `gorm:"not null"`
	EndTime   *time.Time
	Notes     string `gorm:"type:text"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Baby      Baby `gorm:"foreignKey:BabyID"`
}

func (a *Activity) BeforeCreate(tx *gorm.DB) error {
	a.ID = uuid.New()
	return nil
}
