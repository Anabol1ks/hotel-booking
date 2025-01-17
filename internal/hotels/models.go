package hotels

import (
	"gorm.io/gorm"
)

type Hotel struct {
	gorm.Model
	Name          string  `gorm:"type:varchar(100);not null"`
	Address       string  `gorm:"type:varchar(255);not null"`
	Description   string  `gorm:"type:text"`
	OwnerID       uint    `gorm:"not null"` // ID владельца отеля
	AverageRating float64 `gorm:"default:0"`
	RatingsCount  int     `gorm:"default:0"`
	Rooms         []Room
	Ratings       []HotelRating
}

type Room struct {
	gorm.Model
	HotelID       uint    `gorm:"not null"`                  // ID отеля
	RoomType      string  `gorm:"type:varchar(50);not null"` // Тип номера (стандартный, люкс и т.д.)
	Price         float64 `gorm:"not null"`                  // Цена за ночь
	Amenities     string  `gorm:"type:text"`                 // Удобства
	Capacity      int     `gorm:"not null"`                  // Количество гостей
	Available     bool    `gorm:"default:true"`              // Наличие
	AverageRating float64 `gorm:"default:0"`
	RatingsCount  int     `gorm:"default:0"`
	Ratings       []RoomRating
	Images        []RoomImage
}

type RoomImage struct {
	gorm.Model
	RoomID    uint   `gorm:"not null"`                   // ID номера
	ImageURL  string `gorm:"type:varchar(255);not null"` // URL изображения
	ImageName string `gorm:"type:varchar(100);not null"` // Имя изображения
}

type Favorite struct {
	gorm.Model
	UserID uint `gorm:"user_id"`
	RoomID uint `gorm:"room_id"`
}

type HotelRating struct {
	gorm.Model
	HotelID uint    `gorm:"not null"`
	UserID  uint    `gorm:"not null"`
	Rating  float64 `gorm:"not null;check:rating >= 1 AND rating <= 5"`
	Comment string  `gorm:"type:text"`
}

type RoomRating struct {
	gorm.Model
	RoomID  uint    `gorm:"not null"`
	UserID  uint    `gorm:"not null"`
	Rating  float64 `gorm:"not null;check:rating >= 1 AND rating <= 5"`
	Comment string  `gorm:"type:text"`
}
