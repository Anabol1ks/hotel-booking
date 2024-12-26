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

// @Security BearerAuth
// CreateHotelHandler godoc
// @Summary Создание отеля владельцем
// @Description Обрабатывает запрос на создание нового отеля. Только владельцы могут создавать новые отели.
// @Tags hotels
// @Accept json
// @Produce json
// @Param input body CreateHotelInput true "Данные для создания нового отеля"
// @Success 201 {object} response.HotelResponse "Новый отель успешно создан"
// @Failure 400 {object} response.ErrorResponse "Ошибка валидации входных данных"
// @Failure 403 {object} response.ErrorResponse "Только владельцы могут создавать отели"
// @Failure 500 {object} response.ErrorResponse "Ошибка при создании отеля"
// @Router /owners/hotels [post]
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

	c.JSON(http.StatusCreated, hotel)
}

// GetHotelsHandler godoc
// @Summary Получение списка отелей
// @Description Возвращает список всех отелей, включая связанные номера.
// @Tags hotels
// @Produce json
// @Success 200 {object} []response.HotelResponse "Список отелей"
// @Failure 500 {object} response.ErrorResponse "Ошибка при получении отелей"
// @Router /hotels [get]
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

// @Security BearerAuth
// CreateRoomHandler godoc
// @Summary Создание нового номера
// @Description Создает новый номер в отеле. Доступно только для владельцев.
// @Tags rooms
// @Accept json
// @Produce json
// @Param input body CreateRoomInput true "Данные номера"
// @Success 201 {object} response.RoomResponse "Созданный номер"
// @Failure 400 {object} response.ErrorResponse "Ошибка валидации"
// @Failure 403 {object} response.ErrorResponse "Доступ запрещен"
// @Failure 404 {object} response.ErrorResponse "Отель не найден"
// @Failure 500 {object} response.ErrorResponse "Ошибка сервера"
// @Router /owners/rooms [post]
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

	c.JSON(http.StatusCreated, room)
}

// GetRoomsHandler godoc
// @Summary Получение списка номеров
// @Description Возвращает отфильтрованный список номеров с возможностью фильтрации по цене и вместимости
// @Tags rooms
// @Produce json
// @Param min_price query string false "Минимальная цена"
// @Param max_price query string false "Максимальная цена"
// @Param capacity query string false "Минимальная вместимость"
// @Success 200 {array} response.RoomResponse "Список номеров"
// @Failure 500 {object} response.ErrorResponse "Ошибка при получении номеров"
// @Router /rooms [get]
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

// @Security BearerAuth
// GetOwnerHotelsHandler godoc
// @Summary Получение списка отелей владельца
// @Description Возвращает список отелей, принадлежащих текущему владельцу
// @Tags hotels
// @Produce json
// @Success 200 {array} response.HotelResponse "Список отелей"
// @Failure 403 {object} response.ErrorResponse "Доступ запрещен"
// @Failure 500 {object} response.ErrorResponse "Ошибка при получении отелей"
// @Router /owners/hotels [get]
func GetOwnerHotelsHandler(c *gin.Context) {
	ownerID := c.GetUint("user_id")
	role := c.GetString("role")

	if role != "owner" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Доступ запрещен"})
		return
	}

	var hotels []Hotel
	if err := storage.DB.Where("owner_id = ?", ownerID).Find(&hotels).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении отелей"})
		return
	}

	c.JSON(http.StatusOK, hotels)
}
