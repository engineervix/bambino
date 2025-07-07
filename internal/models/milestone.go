package models

import "github.com/google/uuid"

type Milestone struct {
	ActivityID    uuid.UUID `gorm:"type:varchar(36);primary_key"`
	MilestoneType string    `gorm:"type:varchar(50);not null"`
	Description   string    `gorm:"type:text"`
	Activity      Activity  `gorm:"foreignKey:ActivityID"`
}
