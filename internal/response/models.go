package response

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
