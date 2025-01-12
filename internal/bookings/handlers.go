package bookings

import (
	"encoding/json"
	"fmt"
	"hotel-booking/internal/auth"
	"hotel-booking/internal/email"
	"hotel-booking/internal/hotels"
	"hotel-booking/internal/storage"
	"hotel-booking/internal/users"
	"log"
	"net/http"
	"os"
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

	if input.StartDate.After(input.EndDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Дата заезда не может быть позже даты выезда"})
		return
	}

	if input.StartDate.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Дата заезда не может быть в прошлом"})
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
		CreatedAt: time.Now(),
	}

	if err := storage.DB.Create(&booking).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании бронирования"})
		return
	}
	NotificationCreateBooking(userID, booking)

	c.JSON(http.StatusCreated, booking)
}

func NotificationCreateBooking(userID uint, booking Booking) {
	var user users.User
	if err := storage.DB.First(&user, userID).Error; err != nil {
		log.Printf("Ошибка при получении пользователя: %v", err)
		return
	}

	name := user.Name
	emailUs := user.Email

	// Использование системного токена
	systemToken := auth.GetSystemToken()

	createPayUrl := fmt.Sprintf(os.Getenv("URL_BACKEND")+"/bookings/%d/pay", booking.ID)
	req, err := http.NewRequest("POST", createPayUrl, nil)
	if err != nil {
		log.Printf("Ошибка при создании запроса: %v", err)
		return
	}
	req.Header.Set("Authorization", "Bearer "+systemToken)
	req.Header.Set("Content-Type", "application/json")

	// Отправка запроса
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Ошибка при отправке запроса: %v", err)
	}
	defer resp.Body.Close()

	// Проверка статуса ответа
	if resp.StatusCode != http.StatusOK {
		log.Printf("Ошибка: статус ответа %d", resp.StatusCode)
	}

	// Парсинг ответа
	var responseBody map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		log.Printf("Ошибка при чтении ответа: %v", err)
	}

	// Получение ссылки
	paymentURL, ok := responseBody["payment_url"]
	if !ok {
		log.Printf("Поле 'payment_url' отсутствует в ответе")
	}

	emailTemplate := `<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Вами было создано бронирование</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            background-color: #f9f9f9;
            margin: 0;
            padding: 0;
        }
        .email-container {
            max-width: 600px;
            margin: 20px auto;
            background: #ffffff;
            border: 1px solid #ddd;
            border-radius: 8px;
            overflow: hidden;
        }
        .email-header {
            background-color: #007bff;
            color: #ffffff;
            padding: 20px;
            text-align: center;
        }
        .email-header h1 {
            margin: 0;
            font-size: 24px;
        }
        .email-body {
            padding: 20px;
            color: #333333;
        }
        .email-body p {
            margin: 0 0 15px;
            line-height: 1.5;
        }
        .email-footer {
            background-color: #f4f4f9;
            text-align: center;
            padding: 10px;
            font-size: 12px;
            color: #777;
        }
        .reset-button {
            display: inline-block;
            margin-top: 20px;
            padding: 10px 20px;
            background-color:rgb(0, 255, 94);
            color: #ffffff;
            text-decoration: none;
            border-radius: 5px;
            font-size: 16px;
        }
        .reset-button:hover {
            background-color: #0056b3;
        }
    </style>
</head>
<body>
    <div class="email-container">
        <div class="email-header">
            <h1>Вами было создано бронирование</h1>
        </div>
        <div class="email-body">
            <p>Здравствуйте, %v</p>
            <p>Вы только что забронировали номер в отеле.</p>
						<p>Подробности бронирования:</p>
						<p>Номер: %d</p>
						<p>Дата заезда: %s</p>
						<p>Дата выезда: %s</p>
						<p>Стоимость: %f</p>
						<p> Пожалуйста, оплатите его в течении 30 минут с момента отправки письма. Оплатить можно по кнопке снизу, или в вашем списке бронирований: </p>
            <p>https://hotel-booking-sandy.vercel.app/my-bookings</p>
            <a href="%s" class="reset-button">Перейти к оплате</a>
            <p>Если если это были не вы, или бронь было оформленна случайно, то отмените её в списках своих бронирований, или бронь удалиться сама через 30 минут.</p>
            <p>С уважением,<br>Команда поддержки</p>
        </div>
        <div class="email-footer">
            Это письмо было отправлено автоматически. Пожалуйста, не отвечайте на него.
        </div>
    </div>
</body>
</html>`

	subject := "Вами было создано бронирование"
	body := fmt.Sprintf(emailTemplate, name, booking.RoomID, booking.StartDate.Format("02.01.2006"), booking.EndDate.Format("02.01.2006"), booking.TotalCost, paymentURL)
	if err := email.SendEmail(emailUs, subject, body); err != nil {
		log.Printf("Ошибка при отправке письма: %v", err)
	}

}

type CreateOfflineBookingInput struct {
	RoomID      uint      `json:"room_id" binding:"required"`
	StartDate   time.Time `json:"start_date" binding:"required"`
	EndDate     time.Time `json:"end_date" binding:"required"`
	PhoneNumber string    `json:"phone_number" binding:"required"`
	Name        string    `json:"name" binding:"required"`
}

// @Security BearerAuth
// CreateOfflineBookingHandler godoc
// @Summary Создание брони для офлайн клиента
// @Description Создание брони менеджером или владельцем для клиента без аккаунта
// @Tags bookings
// @Accept json
// @Produce json
// @Param input body CreateOfflineBookingInput true "Данные для офлайн бронирования"
// @Success 201 {object} response.BookingResponse
// @Failure 403 {object} response.ErrorResponse "Только менеджеры и владельцы могут создавать офлайн бронирования"
// @Failure 400 {object} response.ErrorResponse "Ошибка валидации"
// @Failure 404 {object} response.ErrorResponse "Номер не найден"
// @Failure 409 {object} response.ErrorResponse "Номер уже забронирован в этот период"
// @Failure 500 {object} response.ErrorResponse "Ошибка при проверке доступности номера или при создании бронирования"
// @Router /bookings/offline [post]
func CreateOfflineBookingHandler(c *gin.Context) {
	// Проверка роли
	role := c.GetString("role")
	if role != "owner" && role != "manager" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Только менеджеры и владельцы могут создавать офлайн бронирования"})
		return
	}

	var input CreateOfflineBookingInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.StartDate.After(input.EndDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Дата заезда не может быть позже даты выезда"})
		return
	}

	if input.StartDate.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Дата заезда не может быть в прошлом"})
		return
	}

	// Поиск существующего пользователя по телефону
	var user users.User
	result := storage.DB.Where("phone = ?", input.PhoneNumber).First(&user)

	if result.Error != nil {
		// Создаем нового пользователя
		user = users.User{
			Phone: input.PhoneNumber,
			Name:  input.Name,
			Role:  "client",
			// Генерируем временный пароль или оставляем пустым
		}
		if err := storage.DB.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании пользователя"})
			return
		}
	}

	// Проверка номера
	var room hotels.Room
	if err := storage.DB.First(&room, input.RoomID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Номер не найден"})
		return
	}

	// Проверка доступности
	var overlappingBookings []Booking
	if err := storage.DB.Where(
		"room_id = ? AND NOT (end_date <= ? OR start_date >= ?)",
		input.RoomID,
		input.StartDate,
		input.EndDate,
	).Find(&overlappingBookings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при проверке доступности номера"})
		return
	}

	if len(overlappingBookings) > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Номер уже забронирован в этот период"})
		return
	}

	// Расчет стоимости
	days := input.EndDate.Sub(input.StartDate).Hours() / 24
	totalPrice := days * float64(room.Price)

	// Создание бронирования
	booking := Booking{
		RoomID:           input.RoomID,
		UserID:           user.ID,
		StartDate:        input.StartDate,
		EndDate:          input.EndDate,
		TotalCost:        totalPrice,
		CreatedAt:        time.Now(),
		PaymentStatus:    "pending", // Офлайн бронирования считаются оплаченными
		IsOfflineBooking: true,
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

// @Security BearerAuth
// GetManagerBookingsHandler godoc
// @Summary Получунеи своих бронирований
// @Description Получение бронирований для пользователя
// @Tags bookings
// @Produce json
// @Success 200 {array} response.BookingResponse "Данные о бранировании"
// Failure 500 {object} response.ErrorResponse "Ошибка при получении бронирований"
// @Router /bookings/my [get]
func GetYourBookingsHandler(c *gin.Context) {
	userID := c.GetUint("user_id")

	var bookings []Booking
	if err := storage.DB.Where("user_id = ?", userID).Find(&bookings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении бронирований"})
		return
	}
	c.JSON(http.StatusOK, bookings)
}

// @Security BearerAuth
// CancelBookingHandler godoc
// @Summary Отмена бронирования
// @Description Отмена бронирования пользователем
// @Tags bookings
// @Param id path int true "ID бронирования"
// @Produce json
// @Success 200 {object} response.MessageResponse "Бронирование успешно отменено"
// @Failure 400 {object} response.ErrorResponse "Бронирование уже оплачено и не может быть отменено"
// @Failure 403 {object} response.ErrorResponse "Вы не можете отменить бронирование, которое не принадлежит вам"
// @Failure 404 {object} response.ErrorResponse "Бронирование не найдено"
// @Failure 500 {object} response.ErrorResponse "Ошибка при отмене бронирования"
// @Router /bookings/{id} [delete]
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

func init() {
	// Start the cleanup goroutine
	go cleanupUnpaidBookings()
}

func cleanupUnpaidBookings() {
	log.Println("Запуск очистки просроченных бронирований...")
	ticker := time.NewTicker(3 * time.Minute)
	for range ticker.C {
		var bookings []Booking
		thirtyMinutesAgo := time.Now().Add(-30 * time.Minute)

		// Находим только онлайн-бронирования с истекшим сроком
		if err := storage.DB.Where(
			"created_at <= ? AND payment_status = ? AND is_offline_booking = ?",
			thirtyMinutesAgo,
			"pending",
			false,
		).Find(&bookings).Error; err != nil {
			log.Printf("Ошибка при обнаружении просроченных бронирований: %v", err)
			continue
		}

		for _, booking := range bookings {
			if err := storage.DB.Delete(&booking).Error; err != nil {
				log.Printf("Ошибка при отмене бронирования с истекшим сроком действия %d: %v", booking.ID, err)
				continue
			}
			log.Printf("Отмененное бронирование с истекшим сроком действия %d", booking.ID)
		}
	}
}
