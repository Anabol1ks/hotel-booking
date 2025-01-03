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
// @Description Возвращает отфильтрованный список номеров с возможностью фильтрации по цене, вместимости, датам бронирования и отелю
// @Tags rooms
// @Produce json
// @Param min_price query string false "Минимальная цена"
// @Param max_price query string false "Максимальная цена"
// @Param capacity query string false "Минимальная вместимость"
// @Param start_date query string false "Дата начала (YYYY-MM-DD)"
// @Param end_date query string false "Дата окончания (YYYY-MM-DD)"
// @Param hotel_id query string false "ID отеля"
// @Success 200 {array} response.RoomResponse "Список номеров"
// @Failure 500 {object} response.ErrorResponse "Ошибка при получении номеров"
// @Router /rooms [get]
func GetRoomsHandler(c *gin.Context) {
	var rooms []Room

	query := storage.DB

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if startDate != "" && endDate != "" {
		query = query.Where("id NOT IN (?)",
			storage.DB.Table("bookings").
				Select("room_id").
				Where("(start_date <= ? AND end_date >= ?) OR (start_date <= ? AND end_date >= ?) OR (start_date >= ? AND end_date <= ?)",
					endDate, startDate, endDate, startDate, startDate, endDate))
	}

	hotelID := c.Query("hotel_id")
	if hotelID != "" {
		query = query.Where("hotel_id = ?", hotelID)
	}
	// фильтры цен
	minPrice := c.Query("min_price")
	maxPrice := c.Query("max_price")

	if minPrice != "" {
		query = query.Where("price >= ?", minPrice)
	}

	if maxPrice != "" {
		query = query.Where("price <= ?", maxPrice)
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

// @Security BearerAuth
// ChangeRoomHandler godoc
// @Summary Изменение номера
// @Description Изменяет существующий номер. Доступно только для владельцев.
// @Tags rooms
// @Accept json
// @Produce json
// @Param id path int true "ID номера"
// @Param input body CreateRoomInput true "Данные номера"
// @Success 200 {object} response.MessageResponse "Номер успешно обновлен"
// @Failure 400 {object} response.ErrorResponse "Ошибка валидации"
// @Failure 403 {object} response.ErrorResponse "Доступ запрещен или номер не принадлежит владельцу"
// @Failure 404 {object} response.ErrorResponse "Номер не найден"
// @Failure 500 {object} response.ErrorResponse "Ошибка при обновлении номера"
// @Router /owners/{id}/room [put]
func ChangeRoomHandler(c *gin.Context) {
	roomID := c.Param("id")
	ownerID := c.GetUint("user_id")
	role := c.GetString("role")

	if role != "owner" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Доступ запрещен"})
		return
	}

	var room CreateRoomInput
	if err := c.ShouldBindJSON(&room); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingRoom Room
	if err := storage.DB.First(&existingRoom, roomID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Номер не найден"})
		return
	}

	var hotel Hotel
	if err := storage.DB.First(&hotel, existingRoom.HotelID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Отель не найден"})
		return
	}

	if hotel.OwnerID != ownerID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Номер не принадлежит вам"})
		return
	}

	if storage.DB.Model(&existingRoom).Updates(Room{
		RoomType:  room.RoomType,
		Price:     room.Price,
		Amenities: room.Amenities,
		Capacity:  room.Capacity,
	}).Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обновлении номера"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Номер успешно обновлен"})
}

// @Security BearerAuth
// DeleteRoomHandler godoc
// @Summary Удаление номера
// @Description Удаляет существующий номер. Доступно только для владельцев.
// @Tags rooms
// @Accept json
// @Produce json
// @Param id path int true "ID номера"
// @Success 200 {object} response.MessageResponse "Номер успешно удален"
// @Failure 403 {object} response.ErrorResponse "Доступ запрещен или номер не принадлежит владельцу"
// @Failure 404 {object} response.ErrorResponse "Номер не найден"
// @Failure 500 {object} response.ErrorResponse "Ошибка при удалении номера"
// @Router /owners/{id}/room [delete]
func DeleteRoomHandler(c *gin.Context) {
	roomID := c.Param("id")
	ownerID := c.GetUint("user_id")
	role := c.GetString("role")

	if role != "owner" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Доступ запрещен"})
		return
	}

	var existingRoom Room
	if err := storage.DB.First(&existingRoom, roomID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Номер не найден"})
		return
	}

	var hotel Hotel
	if err := storage.DB.First(&hotel, existingRoom.HotelID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Отель не найден"})
		return
	}

	if hotel.OwnerID != ownerID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Номер не принадлежит вам"})
		return
	}

	if err := storage.DB.Delete(&existingRoom).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при удалении номера"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Номер успешно удален"})
}
