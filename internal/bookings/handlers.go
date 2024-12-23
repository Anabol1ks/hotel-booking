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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Номер уже забронирован в этот период"})
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

func GetRoomBookingsHandler(c *gin.Context) {
	roomID := c.Param("id")

	var bookings []Booking
	if err := storage.DB.Where("room_id = ?", roomID).Find(&bookings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении списка бронирований"})
		return
	}

	c.JSON(http.StatusOK, bookings)
}
