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
	ID                uuid.UUID    `gorm:"type:varchar(36);primary_key"`
	BabyID            uuid.UUID    `gorm:"type:varchar(36);not null;index"`
	Type              ActivityType `gorm:"type:varchar(20);not null"`
	StartTime         time.Time    `gorm:"not null"`
	EndTime           *time.Time
	Notes             string `gorm:"type:text"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	Baby              Baby               `gorm:"foreignKey:BabyID;constraint:OnDelete:CASCADE"`
	FeedActivity      *FeedActivity      `gorm:"foreignKey:ActivityID;constraint:OnDelete:CASCADE"`
	DiaperActivity    *DiaperActivity    `gorm:"foreignKey:ActivityID;constraint:OnDelete:CASCADE"`
	SleepActivity     *SleepActivity     `gorm:"foreignKey:ActivityID;constraint:OnDelete:CASCADE"`
	PumpActivity      *PumpActivity      `gorm:"foreignKey:ActivityID;constraint:OnDelete:CASCADE"`
	GrowthMeasurement *GrowthMeasurement `gorm:"foreignKey:ActivityID;constraint:OnDelete:CASCADE"`
	HealthRecord      *HealthRecord      `gorm:"foreignKey:ActivityID;constraint:OnDelete:CASCADE"`
	Milestone         *Milestone         `gorm:"foreignKey:ActivityID;constraint:OnDelete:CASCADE"`
}

func (a *Activity) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

// BeforeSave hook to validate required fields
func (a *Activity) BeforeSave(tx *gorm.DB) error {
	if a.BabyID == uuid.Nil {
		return gorm.ErrInvalidField
	}
	if a.Type == "" {
		return gorm.ErrInvalidField
	}
	return nil
}
