package models

import "github.com/google/uuid"

type SleepActivity struct {
	ActivityID uuid.UUID `gorm:"type:varchar(36);primary_key"`
	Location   string    `gorm:"type:varchar(50)"` // crib, bassinet, car_seat, etc.
	Quality    *int      `gorm:"check:quality >= 1 AND quality <= 5"`
	Activity   Activity  `gorm:"foreignKey:ActivityID"`
}
