package hotels

import (
	"hotel-booking/internal/storage"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateHotelInput struct {
	Name        string `json:"name" binding:"required"`
	Addres      string `json:"address" binding:"required"`
	Description string `json:"description"`
}

// создание отеля
func CreateHotelHandler(c *gin.Context) {
	ownerID := c.GetUint("user_id")
	role := c.GetString("role")

	if role != "owner" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Доступ запрещен"})
		return
	}

	var input CreateHotelInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hotel := Hotel{
		Name:        input.Name,
		Address:     input.Addres,
		Description: input.Description,
		OwnerID:     ownerID,
	}

	if err := storage.DB.Create(&hotel).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании отеля"})
		return
	}

	c.JSON(http.StatusOK, hotel)
}

func GetHotelsHandler(c *gin.Context) {
	var hotels []Hotel
	if err := storage.DB.Preload("Rooms").Find(&hotels).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении отелей"})
		return
	}
	c.JSON(http.StatusOK, hotels)
}

type CreateRoomInput struct {
	HotelID   uint    `json:"hotel_id" binding:"required"`
	RoomType  string  `json:"room_type" binding:"required"`
	Price     float64 `json:"price" binding:"required"`
	Amenities string  `json:"amenities"`
	Capacity  int     `json:"capacity" binding:"required"`
}

func CreateRoomHandler(c *gin.Context) {
	role := c.GetString("role")

	if role != "owner" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Только владелец может создавать номера"})
		return
	}

	var input CreateRoomInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var hotel Hotel
	if err := storage.DB.First(&hotel, input.HotelID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Отель не найден"})
		return
	}

	room := Room{
		HotelID:   input.HotelID,
		RoomType:  input.RoomType,
		Price:     input.Price,
		Amenities: input.Amenities,
		Capacity:  input.Capacity,
	}

	if err := storage.DB.Create(&room).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании номера"})
		return
	}

	c.JSON(http.StatusOK, room)
}

func GetRoomsHandler(c *gin.Context) {
	var rooms []Room

	query := storage.DB

	// фильтры цен
	minPrice := c.Query("min_price")
	maxPrice := c.Query("max_price")
	if minPrice != "" && maxPrice != "" {
		query = query.Where("price BETWEEN ? AND ?", minPrice, maxPrice)
	}

	// фильтры количества гостей
	capacity := c.Query("capacity")
	if capacity != "" {
		query = query.Where("capacity >= ?", capacity)
	}

	if err := query.Find(&rooms).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении номеров"})
		return
	}

	c.JSON(http.StatusOK, rooms)
}
