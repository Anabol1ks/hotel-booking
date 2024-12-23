package bookings

import (
	"time"

	"gorm.io/gorm"
)

type Booking struct {
	gorm.Model
	RoomID        uint      `gorm:"not null"`
	UserID        uint      `gorm:"not null"`
	StartDate     time.Time `gorm:"not null"`
	EndDate       time.Time `gorm:"not null"`
	TotalCost     float64   `gorm:"not null"` //Итоговая стоимость
	PaymentStatus string    `gorm:"type:varchar(20);default:'pending'"`
}
