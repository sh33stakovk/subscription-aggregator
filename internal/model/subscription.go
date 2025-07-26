package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Subscription struct {
	gorm.Model
	ServiceName string    `gorm:"not null"                  json:"service_name"`
	Price       uint      `gorm:"not null;check:price >= 0" json:"price"`
	UserID      uuid.UUID `gorm:"type:uuid;not null"        json:"user_id"`
	StartDate   time.Time `gorm:"type:date;not null"        json:"start_date"`
}
