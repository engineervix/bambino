package models

import (
	"database/sql/driver"
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

func (a *ActivityType) Scan(value interface{}) error {
	*a = ActivityType(value.(string))
	return nil
}

func (a ActivityType) Value() (driver.Value, error) {
	return string(a), nil
}

type Activity struct {
	ID                uuid.UUID    `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	BabyID            uuid.UUID    `gorm:"type:uuid;not null"`
	Type              ActivityType `gorm:"type:varchar(20);not null"`
	StartTime         time.Time    `gorm:"not null"`
	EndTime           *time.Time
	Notes             string `gorm:"type:text"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	Baby              Baby               `gorm:"foreignKey:BabyID"`
	FeedActivity      *FeedActivity      `gorm:"foreignKey:ActivityID"`
	DiaperActivity    *DiaperActivity    `gorm:"foreignKey:ActivityID"`
	SleepActivity     *SleepActivity     `gorm:"foreignKey:ActivityID"`
	PumpActivity      *PumpActivity      `gorm:"foreignKey:ActivityID"`
	GrowthMeasurement *GrowthMeasurement `gorm:"foreignKey:ActivityID"`
	HealthRecord      *HealthRecord      `gorm:"foreignKey:ActivityID"`
	Milestone         *Milestone         `gorm:"foreignKey:ActivityID"`
}

func (a *Activity) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}
