package response

import (
	"time"
)

// ErrorResponse представляет стандартный формат ответа при ошибке
// @Description Стандартный ответ при ошибке
type ErrorResponse struct {
	Error string `json:"error"`
}

// SuccessResponse представляет стандартный формат успешного ответа
// @Description Стандартный ответ при успешном выполнении
type SuccessResponse struct {
	Message string `json:"message"`
}

type TokenResponse struct {
	Token string `json:"token" example:"Ваш токен"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type UserResponse struct {
	ID    uint   `json:"ID"`
	Name  string `json:"Name"`
	Email string `json:"Email"`
	Phone string `json:"Phone"`
	Role  string `json:"Role"`
}

type HotelResponse struct {
	Name        string         `json:"name"`
	Address     string         `json:"address"`
	Description string         `json:"description"`
	OwnerID     uint           `json:"owner_id"`
	Rooms       []RoomResponse `json:"rooms"`
}

type RoomResponse struct {
	HotelID   uint    `json:"hotel_id"`  // ID отеля
	RoomType  string  `json:"room_type"` // Тип номера (стандартный, люкс и т.д.)
	Price     float64 `json:"price"`     // Цена за ночь
	Amenities string  `json:"amenities"` // Удобства
	Capacity  int     `json:"capacity"`  // Количество гостей
	Available bool    `json:"available"` // Наличие
}

type BookingResponse struct {
	RoomID        uint      `json:"room_id"`
	UserID        uint      `json:"user_id"`
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`
	TotalCost     float64   `json:"total_cost"`     //Итоговая стоимость
	PaymentStatus string    `json:"payment_status"` //Статус оплаты
}

type CreatePaymentResponse struct {
	PaymentURL string `json:"payment_url" example:"ссылка на оплату"` // Ссылка для оплаты
}

type HotelRatingResponse struct {
	HotelID uint    `json:"hotel_id"`
	UserID  uint    `json:"user_id"`
	Rating  float64 `json:"rating"`
	Comment string  `json:"comment"`
}

type RoomRatingResponse struct {
	RoomID  uint    `json:"room_id"`
	UserID  uint    `json:"user_id"`
	Rating  float64 `json:"rating"`
	Comment string  `json:"comment"`
}
