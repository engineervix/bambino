package models

import "github.com/google/uuid"

type HealthRecordType string

const (
	HealthRecordTypeCheckup HealthRecordType = "checkup"
	HealthRecordTypeVaccine HealthRecordType = "vaccine"
	HealthRecordTypeIllness HealthRecordType = "illness"
)

type HealthRecord struct {
	ActivityID  uuid.UUID        `gorm:"type:varchar(36);primary_key"`
	RecordType  HealthRecordType `gorm:"type:varchar(20);not null"`
	Provider    string           `gorm:"type:varchar(100)"`
	VaccineName string           `gorm:"type:varchar(100)"`
	Symptoms    string           `gorm:"type:text"`
	Treatment   string           `gorm:"type:text"`
	Activity    Activity         `gorm:"foreignKey:ActivityID"`
}
