package bookings

import (
	"hotel-booking/internal/hotels"
	"hotel-booking/internal/storage"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type CreateBookingInput struct {
	RoomID    uint      `json:"room_id" binding:"required"`
	StartDate time.Time `json:"start_date" binding:"required"`
	EndDate   time.Time `json:"end_date" binding:"required"`
}

// @Security BearerAuth
// CreateBookingHandler godoc
// @Summary Бронирование номера
// @Description Бронирование номера только для авторизованных пользователей
// @Tags bookings
// @Produce json
// @Param input body CreateBookingInput true "Данные для бронирования"
// @Success 201 {object} response.BookingResponse "Данные о бранировании"
// @Failure 400 {object} response.ErrorResponse "Ошибка валидации"
// @Failure 404 {object} response.ErrorResponse "Номер не найден"
// @Failure 409 {object} response.ErrorResponse "Номер уже забронирован в этот период"
// @Failure 500 {object} response.ErrorResponse "Ошибка при проверке доступности номера или при создании бронирования"
// @Router /bookings [post]
func CreateBookingHandler(c *gin.Context) {
	userID := c.GetUint("user_id")

	var input CreateBookingInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверка номера
	var room hotels.Room
	if err := storage.DB.First(&room, input.RoomID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Номер не найден"})
		return
	}

	// доступность номера
	var overlappingBookings []Booking
	if err := storage.DB.Where("room_id = ? AND NOT (end_date <= ? OR start_date >= ?)", input.RoomID, input.StartDate, input.EndDate).Find(&overlappingBookings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при проверке доступности номера"})
		return
	}

	if len(overlappingBookings) > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Номер уже забронирован в этот период"})
		return
	}

	days := input.EndDate.Sub(input.StartDate).Hours() / 24
	totalPrice := days * float64(room.Price)

	// создание бронирования
	booking := Booking{
		RoomID:    input.RoomID,
		UserID:    userID,
		StartDate: input.StartDate,
		EndDate:   input.EndDate,
		TotalCost: totalPrice,
	}

	if err := storage.DB.Create(&booking).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании бронирования"})
		return
	}

	c.JSON(http.StatusCreated, booking)
}

// GetRoomBookingsHandler godoc
// @Summary Получение бронирований для номера
// @Description Получение бронирований для номера
// @Tags rooms
// @Produce json
// @Param id path int true "ID номера"
// @Success 201 {array} response.BookingResponse "Данные о бранировании"
// @Failure 500 {object} response.ErrorResponse "Ошибка при получении списка бронирований"
// @Router /rooms/{id}/bookings [get]
func GetRoomBookingsHandler(c *gin.Context) {
	roomID := c.Param("id")

	var bookings []Booking
	if err := storage.DB.Where("room_id = ?", roomID).Find(&bookings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении списка бронирований"})
		return
	}

	c.JSON(http.StatusOK, bookings)
}

// @Security BearerAuth
// GetOwnerBookingsHandler godoc
// @Summary Получение бронирований владельца
// @Description Получение бронирований для владельца
// @Tags bookings
// @Produce json
// @Success 201 {array} response.BookingResponse "Данные о бранировании"
// @Failure 403 {object} response.ErrorResponse "Доступ запрещен"
// @Failure 500 {object} response.ErrorResponse "Ошибка при получении бронирований"
// @Router /owners/bookings [get]
func GetOwnerBookingsHandler(c *gin.Context) {
	ownerID := c.GetUint("user_id")

	role := c.GetString("role")

	if role != "owner" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Доступ запрещен"})
		return
	}

	var bookings []Booking
	query := `
	SELECT b.*
	FROM bookings b
	JOIN rooms r ON b.room_id = r.id
	JOIN hotels h ON r.hotel_id = h.id
	WHERE h.owner_id = ?
	`

	if err := storage.DB.Raw(query, ownerID).Scan(&bookings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении бронирований"})
		return
	}

	c.JSON(http.StatusOK, bookings)
}

func CancelBookingHandler(c *gin.Context) {
	bookingID := c.Param("id")
	userID := c.GetUint("user_id")

	var booking Booking
	if err := storage.DB.First(&booking, bookingID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Бронирование не найдено"})
		return
	}
	if booking.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Вы не можете отменить бронирование, которое не принадлежит вам"})
		return
	}

	if booking.PaymentStatus == "succeeded" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Бронирование уже оплачено и не может быть отменено"})
		return
	}

	if err := storage.DB.Delete(&booking).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при отмене бронирования"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Бронирование успешно отменено"})
}
