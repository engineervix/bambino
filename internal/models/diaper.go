package models

import "github.com/google/uuid"

type DiaperActivity struct {
	ActivityID  uuid.UUID `gorm:"type:varchar(36);primary_key"`
	Wet         bool      `gorm:"default:false"`
	Dirty       bool      `gorm:"default:false"`
	Color       string    `gorm:"type:varchar(20)"`
	Consistency string    `gorm:"type:varchar(20)"`
	Activity    Activity  `gorm:"foreignKey:ActivityID"`
}
