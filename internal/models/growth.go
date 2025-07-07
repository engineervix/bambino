package models

import "github.com/google/uuid"

type GrowthMeasurement struct {
	ActivityID          uuid.UUID `gorm:"type:uuid;primary_key"`
	WeightKG            *float64  `gorm:"type:decimal(5,2)"`
	HeightCM            *float64  `gorm:"type:decimal(5,1)"`
	HeadCircumferenceCM *float64  `gorm:"type:decimal(4,1)"`
	Activity            Activity  `gorm:"foreignKey:ActivityID"`
}
