package hotels

import (
	"gorm.io/gorm"
)

type Hotel struct {
	gorm.Model
	Name        string `gorm:"type:varchar(100);not null"`
	Address     string `gorm:"type:varchar(255);not null"`
	Description string `gorm:"type:text"`
	OwnerID     uint   `gorm:"not null"` // ID владельца отеля
	Rooms       []Room
}

type Room struct {
	gorm.Model
	HotelID   uint    `gorm:"not null"`                  // ID отеля
	RoomType  string  `gorm:"type:varchar(50);not null"` // Тип номера (стандартный, люкс и т.д.)
	Price     float64 `gorm:"not null"`                  // Цена за ночь
	Amenities string  `gorm:"type:text"`                 // Удобства
	Capacity  int     `gorm:"not null"`                  // Количество гостей
	Available bool    `gorm:"default:true"`              // Наличие
}

type Favorite struct {
	gorm.Model
	UserID uint `gorm:"user_id"`
	RoomID uint `gorm:"room_id"`
}
