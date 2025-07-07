package models

import "github.com/google/uuid"

type PumpBreast string

const (
	PumpBreastLeft  PumpBreast = "left"
	PumpBreastRight PumpBreast = "right"
	PumpBreastBoth  PumpBreast = "both"
)

type PumpActivity struct {
	ActivityID      uuid.UUID  `gorm:"type:varchar(36);primary_key"`
	Breast          PumpBreast `gorm:"type:varchar(10)"`
	AmountML        *float64   `gorm:"type:decimal(5,1)"`
	DurationMinutes *int
	Activity        Activity `gorm:"foreignKey:ActivityID"`
}
