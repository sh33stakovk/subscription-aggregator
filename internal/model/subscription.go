package model

import (
	monthyear "subscription-aggregator/pkg/month-year"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Subscription struct {
	gorm.Model
	ServiceName string              `gorm:"not null"                  json:"service_name"`
	Price       uint                `gorm:"not null;check:price >= 0" json:"price"`
	UserID      uuid.UUID           `gorm:"type:uuid;not null"        json:"user_id"`
	StartDate   monthyear.MonthYear `gorm:"type:date;not null"        json:"start_date"`
}
