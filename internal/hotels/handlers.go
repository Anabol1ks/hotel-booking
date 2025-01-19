package hotels

import (
	"hotel-booking/internal/storage"
	"io"
	"net/http"
	"os"

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

	query := storage.DB.Preload("Images")

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
// GetOwnerRoomsHandler godoc
// @Summary Получение списка номеров владельца
// @Description Возвращает список всех номеров в отелях, принадлежащих текущему владельцу
// @Tags rooms
// @Produce json
// @Param hotel_id query string false "ID отеля для фильтрации"
// @Success 200 {array} response.RoomResponse "Список номеров"
// @Failure 403 {object} response.ErrorResponse "Доступ запрещен"
// @Failure 500 {object} response.ErrorResponse "Ошибка при получении номеров"
// @Router /owners/rooms [get]
func GetOwnerRoomsHandler(c *gin.Context) {
	ownerID := c.GetUint("user_id")
	role := c.GetString("role")
	hotelID := c.Query("hotel_id")

	if role != "owner" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Доступ запрещен"})
		return
	}

	query := storage.DB.Table("rooms").Joins("JOIN hotels ON rooms.hotel_id = hotels.id").Joins("JOIN users ON hotels.owner_id = users.id").Where("users.id = ?", ownerID)

	if hotelID != "" {
		query = query.Where("hotels.id = ?", hotelID)
	}

	var rooms []Room
	if err := query.Find(&rooms).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении номеров"})
		return
	}

	c.JSON(http.StatusOK, rooms)
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

// @Security BearerAuth
// AddToFavoritesHandler godoc
// @Summary Добавление номера в избранное
// @Description Добавляет номер в список избранных пользователя
// @Tags favorites
// @Accept json
// @Produce json
// @Param room_id path int true "ID номера"
// @Success 201 {object} response.MessageResponse "Номер успешно добавлен в избранное"
// @Failure 400 {object} response.ErrorResponse "Номер уже в избранном"
// @Failure 404 {object} response.ErrorResponse "Номер не найден"
// @Router /favorites/{room_id} [post]
func AddToFavoritesHandler(c *gin.Context) {
	userID := c.GetUint("user_id")
	roomID := c.Param("room_id")

	// Проверяем существование номера
	var room Room
	if err := storage.DB.First(&room, roomID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Номер не найден"})
		return
	}

	// Проверяем, не добавлен ли номер уже в избранное
	var existing Favorite
	result := storage.DB.Where("user_id = ? AND room_id = ?", userID, roomID).First(&existing)
	if result.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Номер уже в избранном"})
		return
	}

	// Создаем новую запись в избранном
	favorite := Favorite{
		UserID: userID,
		RoomID: room.ID,
	}

	if err := storage.DB.Create(&favorite).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при добавлении в избранное"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Номер успешно добавлен в избранное"})
}

// @Security BearerAuth
// GetFavoritesHandler godoc
// @Summary Получение списка избранных номеров
// @Description Возвращает список избранных номеров пользователя
// @Tags favorites
// @Produce json
// @Success 200 {array} response.RoomResponse "Список избранных номеров"
// @Failure 500 {object} response.ErrorResponse "Ошибка при получении списка избранного"
// @Router /favorites [get]
func GetFavoritesHandler(c *gin.Context) {
	userID := c.GetUint("user_id")

	var rooms []Room
	if err := storage.DB.Joins("JOIN favorites ON rooms.id = favorites.room_id").
		Where("favorites.user_id = ? AND favorites.deleted_at IS NULL", userID).
		Find(&rooms).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении избранных номеров"})
		return
	}

	c.JSON(http.StatusOK, rooms)
}

// @Security BearerAuth
// RemoveFromFavoritesHandler godoc
// @Summary Удаление номера из избранного
// @Description Удаляет номер из списка избранных пользователя
// @Tags favorites
// @Produce json
// @Param room_id path int true "ID номера"
// @Success 200 {object} response.MessageResponse "Номер успешно удален из избранного"
// @Failure 404 {object} response.ErrorResponse "Номер не найден в избранном"
// @Router /favorites/{room_id} [delete]
func RemoveFromFavoritesHandler(c *gin.Context) {
	userID := c.GetUint("user_id")
	roomID := c.Param("room_id")

	result := storage.DB.Where("user_id = ? AND room_id = ?", userID, roomID).Delete(&Favorite{})
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Номер не найден в избранном"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Номер успешно удален из избранного"})
}

// Рейтинговая система
type RatingInput struct {
	Rating  int    `json:"rating" binding:"required"`
	Comment string `json:"comment"`
}

// @Security BearerAuth
// RateHotelHandler godoc
// @Summary Оценка отеля
// @Description Оценивает отель пользователем
// @Tags ratings
// @Accept json
// @Produce json
// @Param hotel_id path int true "ID отеля"
// @Param input body RatingInput true "Рейтинг и комментарий"
// @Success 200 {object} response.MessageResponse "Оценка успешно добавлена"
// @Failure 400 {object} response.ErrorResponse "Недопустимый рейтинг/Вы уже оценили этот отель"
// @Failure 404 {object} response.ErrorResponse "Отель не найден"
// @Router /hotels/{hotel_id}/rate [post]
func RateHotelHandler(c *gin.Context) {
	userID := c.GetUint("user_id")
	hotelID := c.Param("hotel_id")

	var input RatingInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Rating < 1 || input.Rating > 5 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Недопустимый рейтинг"})
		return
	}

	var hotel Hotel
	if err := storage.DB.First(&hotel, hotelID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Отель не найден"})
		return
	}

	var existingRating HotelRating
	result := storage.DB.Where("user_id = ? AND hotel_id = ?", userID, hotelID).First(&existingRating)
	if result.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Вы уже оценили этот отель"})
		return
	}

	rating := HotelRating{
		HotelID: hotel.ID,
		UserID:  userID,
		Rating:  float64(input.Rating),
		Comment: input.Comment,
	}

	tx := storage.DB.Begin()

	if err := tx.Create(&rating).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при сохранении оценки"})
		return
	}

	newAvgRating := (hotel.AverageRating*float64(hotel.RatingsCount) + float64(input.Rating)) / float64(hotel.RatingsCount+1)
	if err := tx.Model(&hotel).Updates(map[string]interface{}{
		"average_rating": newAvgRating,
		"ratings_count":  hotel.RatingsCount + 1,
	}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обновлении рейтинга"})
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"message": "Оценка успешно добавлена"})

}

// @Security BearerAuth
// RateRoomHandler godoc
// @Summary Оценка номера
// @Description Оценивает номер пользователем
// @Tags ratings
// @Accept json
// @Produce json
// @Param room_id path int true "ID номера"
// @Param input body RatingInput true "Рейтинг и комментарий"
// @Success 200 {object} response.MessageResponse "Оценка успешно добавлена"
// @Failure 400 {object} response.ErrorResponse "Недопусти рейтинг/Вы уже оценили этот номер"
// @Failure 404 {object} response.ErrorResponse "Номер не найден"
// @Router /rooms/{room_id}/rate [post]
func RateRoomHandler(c *gin.Context) {
	userID := c.GetUint("user_id")
	roomID := c.Param("id")

	var input RatingInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Rating < 1 || input.Rating > 5 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Недопустимый рейтинг"})
		return
	}

	var room Room
	if err := storage.DB.First(&room, roomID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Номер не найден"})
		return
	}

	var existingRating RoomRating
	result := storage.DB.Where("user_id = ? AND room_id = ?", userID, roomID).First(&existingRating)
	if result.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Вы уже оценили этот номер"})
		return
	}

	rating := RoomRating{
		RoomID:  room.ID,
		UserID:  userID,
		Rating:  float64(input.Rating),
		Comment: input.Comment,
	}

	tx := storage.DB.Begin()

	if err := tx.Create(&rating).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при сохранении оценки"})
		return
	}

	newAvgRating := (room.AverageRating*float64(room.RatingsCount) + float64(input.Rating)) / float64(room.RatingsCount+1)
	if err := tx.Model(&room).Updates(map[string]interface{}{
		"average_rating": newAvgRating,
		"ratings_count":  room.RatingsCount + 1,
	}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обновлении рейтинга"})
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"message": "Оценка успешно добавлена"})
}

// GetHotelsRatingsHandler godoc
// @Summary Получить оценки отеля
// @Description Получает оценки отеля
// @Tags ratings
// @Produce json
// @Param hotel_id path int true "ID отеля"
// @Success 200 {array} []response.HotelRatingResponse "Список оценок отеля"
// @Failure 500 {object} response.ErrorResponse "Ошибка при получении оценок"
// @Router /hotels/{hotel_id}/rate [get]
func GetHotelsRatingsHandler(c *gin.Context) {
	hotelID := c.Param("hotel_id")

	var retings []HotelRating
	if err := storage.DB.Where("hotel_id = ?", hotelID).Find(&retings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении оценок"})
		return
	}

	c.JSON(http.StatusOK, retings)
}

// GetRoomsRatingsHandler godoc
// @Summary Получить оценки номера
// @Description Получает оценки номера
// @Tags ratings
// @Produce json
// @Param room_id path int true "ID номера"
// @Success 200 {array} []response.RoomRatingResponse "Список оценок номера"
// @Failure 500 {object} response.ErrorResponse "Ошибка при получении оценок"
// @Router /rooms/{room_id}/rate [get]
func GetRoomsRatingsHandler(c *gin.Context) {
	roomID := c.Param("id")

	var retings []RoomRating
	if err := storage.DB.Where("room_id = ?", roomID).Find(&retings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении оценок"})
		return
	}

	c.JSON(http.StatusOK, retings)
}

// ---------------------------------------------------------------

// загрузка изображения для номеров и отелей

// @Security BearerAuth
// UploadHotelImagesHandler godoc
// @Summary Загрузка изображений для отеля
// @Description Загружает изображения для отеля
// @Tags images
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "ID номера"
// @Param images formData file true "Изображения"
// @Success 200 {object} response.MessageResponse "Изображения успешно загружены"
// @Failure 400 {object} response.ErrorResponse "Ошибка при загрузке изображений"
// @Failure 404 {object} response.ErrorResponse "Отель не найден"
// @Failure 403 {object} response.ErrorResponse "Доступ запрещен"
// @Router /owners/rooms/{id}/images [post]
func UploadRoomImagesHandler(c *gin.Context) {
	ownerID := c.GetUint("user_id")
	roomID := c.Param("id")

	var room Room
	if err := storage.DB.First(&room, roomID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Номер не найден"})
		return
	}

	var hotel Hotel
	if err := storage.DB.First(&hotel, room.HotelID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Отель не найден"})
		return
	}

	if hotel.OwnerID != ownerID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Доступ запрещен"})
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка при обработке формы"})
		return
	}

	files := form.File["images"]
	webdavService := NewWebDAVService(
		"https://webdav.cloud.mail.ru",
		os.Getenv("WEBDAV_USERNAME"),
		os.Getenv("WEBDAV_PASSWORD"),
	)
	if webdavService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка подключения к облачному хранилищу"})
		return
	}

	for _, file := range files {
		// проверка типа файла
		if !isImageFile(file.Filename) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Недопустимый тип файла"})
			return
		}

		// Чтение файла
		fileData, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при чтении файла"})
			return
		}
		defer fileData.Close()

		data, err := io.ReadAll(fileData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при чтении файла"})
			return
		}

		// Генерация уникального имени
		filename := generateUniqueFilename(file.Filename)

		// Загрузка файла на WebDAV
		imageUrl, err := webdavService.UploadImage(data, filename)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при загрузке файла"})
			return
		}

		roomImage := RoomImage{
			RoomID:    room.ID,
			ImageURL:  imageUrl,
			ImageName: filename,
		}

		if err := storage.DB.Create(&roomImage).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при сохранении изображения"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Изображения успешно загружены"})
}

// @Security BearerAuth
// DeleteRoomImageHandler godoc
// @Summary Удаление изображения номера
// @Description Удаляет изображение номера
// @Tags images
// @Produce json
// @Param room_id path int true "ID номера"
// @Param image_id path int true "ID изображения"
// @Success 200 {object} response.MessageResponse "Изображение успешно удалено"
// @Failure 403 {object} response.ErrorResponse "Доступ запрещен"
// @Failure 404 {object} response.ErrorResponse "Изображение не найдено"
// @Router /owners/rooms/{room_id}/images/{image_id} [delete]
func DeleteRoomImageHandler(c *gin.Context) {
	ownerID := c.GetUint("user_id")
	roomID := c.Param("id")
	imageID := c.Param("image_id")

	var room Room
	if err := storage.DB.First(&room, roomID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Номер не найден"})
		return
	}

	var hotel Hotel
	if err := storage.DB.First(&hotel, room.HotelID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Отель не найден"})
		return
	}

	if hotel.OwnerID != ownerID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Доступ запрещен"})
		return
	}

	var image RoomImage
	if err := storage.DB.Where("id = ? AND room_id = ?", imageID, roomID).First(&image).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Изображение не найдено"})
		return
	}

	webdavService := NewWebDAVService(
		"https://webdav.cloud.mail.ru",
		os.Getenv("WEBDAV_USERNAME"),
		os.Getenv("WEBDAV_PASSWORD"),
	)
	if webdavService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка подключения к облачному хранилищу"})
		return
	}

	if webdavService.DeleteImage(image.ImageName) != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при удалении файла"})
		return
	}

	if err := storage.DB.Delete(&image).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при удалении изображения"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Изображение успешно удалено"})
}

// -------------------------------------------------------------
