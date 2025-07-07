package models

import "github.com/google/uuid"

type FeedType string

const (
	FeedTypeBottle      FeedType = "bottle"
	FeedTypeBreastLeft  FeedType = "breast_left"
	FeedTypeBreastRight FeedType = "breast_right"
	FeedTypeSolid       FeedType = "solid"
)

type FeedActivity struct {
	ActivityID      uuid.UUID `gorm:"type:uuid;primary_key"`
	FeedType        FeedType  `gorm:"type:varchar(20);not null"`
	AmountML        *float64  `gorm:"type:decimal(5,1)"`
	DurationMinutes *int
	Activity        Activity `gorm:"foreignKey:ActivityID"`
}
