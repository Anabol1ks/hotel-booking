package payments

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hotel-booking/internal/bookings"
	"hotel-booking/internal/storage"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PaymentRequest struct {
	Amount struct {
		Value    string `json:"value"`
		Currency string `json:"currency"`
	} `json:"amount"`
	Confirmation struct {
		Type      string `json:"type"`
		ReturnURL string `json:"return_url"`
	} `json:"confirmation"`
	Capture     bool                   `json:"capture"`
	Description string                 `json:"description"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// @Security BearerAuth
// @Summary Создание платежа для бронирования
// @Description Создает платеж через YooKassa для указанного бронирования и возвращает ссылку для оплаты.
// @Tags payments
// @Accept json
// @Produce json
// @Param id path int true "Идентификатор бронирования"
// @Success 200 {object} response.CreatePaymentResponse "Ссылка для оплаты успешно создана"
// @Failure 400 {object} response.ErrorResponse "Некорректный запрос или бронирование уже оплачено"
// @Failure 404 {object} response.ErrorResponse "Бронирование не найдено"
// @Failure 500 {object} response.ErrorResponse "Внутренняя ошибка сервера или ошибка платежной системы"
// @Router /bookings/{id}/pay [post]
func CreatePaymentHandler(c *gin.Context) {
	bookingID := c.Param("id")

	var booking bookings.Booking
	if err := storage.DB.First(&booking, bookingID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Бронирование не найдено"})
		return
	}

	if booking.PaymentStatus == "succeeded" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Бронирование уже оплачено"})
		return
	}

	paymentRequest := PaymentRequest{
		Capture:     true,
		Description: fmt.Sprintf("Оплата бронирования %s", bookingID),
		Metadata:    map[string]interface{}{"booking_id": bookingID}, // Указываем booking_id
	}
	paymentRequest.Amount.Value = fmt.Sprintf("%.2f", booking.TotalCost)
	paymentRequest.Amount.Currency = "RUB"
	paymentRequest.Confirmation.Type = "redirect"
	paymentRequest.Confirmation.ReturnURL = "http://localhost:8080/payment/success"

	requestBody, err := json.Marshal(paymentRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании платежа"})
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://api.yookassa.ru/v3/payments", bytes.NewBuffer(requestBody))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании запроса"})
		return
	}
	req.SetBasicAuth(os.Getenv("YOKASSA_SHOP_ID"), os.Getenv("YOKASSA_SECRET_KEY"))
	req.Header.Set("Content-Type", "application/json")

	// Генерация и установка Idempotence-Key
	idempotenceKey := uuid.New().String()
	req.Header.Set("Idempotence-Key", idempotenceKey)

	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при подключении к ЮKassa"})
		return
	}
	defer resp.Body.Close()

	var responseData map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обработки ответа от ЮKassa"})
		return
	}

	// Проверка статуса ответа
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		c.JSON(resp.StatusCode, responseData)
		return
	}

	// Проверка наличия confirmation_url
	confirmation, ok := responseData["confirmation"].(map[string]interface{})
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось получить платёжную ссылку"})
		return
	}

	confirmationURL, ok := confirmation["confirmation_url"].(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось получить платёжную ссылку"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"payment_url": confirmationURL})
}

// PaymentCallbackHandler обрабатывает уведомления о статусе оплаты от платежной системы.
// @Summary Webhook для обработки статуса оплаты
// @Description Обрабатывает уведомления от платежной системы и обновляет статус оплаты для указанного бронирования.
// @Tags payments
// @Accept json
// @Produce json
// @Param request body PaymentCallbackRequest true "Данные вебхука от платежной системы"
// @Success 200 {object} response.SuccessResponse "Статус оплаты обновлен"
// @Failure 400 {object} response.ErrorResponse "Некорректные данные запроса"
// @Failure 404 {object} response.ErrorResponse "Бронирование не найдено"
// @Failure 500 {object} response.ErrorResponse "Внутренняя ошибка сервера"
// @Router /payments/callback [post]
func PaymentCallbackHandler(c *gin.Context) {
	var callbackData map[string]interface{}
	if err := json.NewDecoder(c.Request.Body).Decode(&callbackData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные"})
		return
	}

	log.Printf("Получен Webhook: %v\n", callbackData)

	// Извлекаем объект "object" из Webhook
	object, ok := callbackData["object"].(map[string]interface{})
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный формат 'object'"})
		return
	}

	// Проверяем статус оплаты
	paymentStatus, ok := object["status"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный формат поля 'status'"})
		return
	}

	// Проверяем наличие и формат metadata
	metadata, ok := object["metadata"].(map[string]interface{})
	if !ok || len(metadata) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Поле 'metadata' отсутствует или пусто"})
		return
	}

	// Проверяем наличие booking_id
	bookingIDRaw, ok := metadata["booking_id"]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Поле 'booking_id' отсутствует в 'metadata'"})
		return
	}

	// Преобразуем booking_id в строку
	bookingID := fmt.Sprintf("%v", bookingIDRaw)

	// Проверяем существование бронирования
	var booking bookings.Booking
	if err := storage.DB.First(&booking, bookingID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Бронирование не найдено"})
		return
	}

	// Обновляем статус оплаты
	booking.PaymentStatus = paymentStatus
	if err := storage.DB.Save(&booking).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обновлении статуса оплаты"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Статус оплаты обновлен"})
}

type PaymentCallbackRequest struct {
	Object PaymentObject `json:"object"` // Основной объект данных
}

// PaymentObject описывает объект `object`, содержащий статус и метаданные.
type PaymentObject struct {
	Status   string          `json:"status" example:"succeeded"` // Статус оплаты
	Metadata PaymentMetadata `json:"metadata"`                   // Метаданные оплаты
}

// PaymentMetadata описывает объект `metadata` с деталями бронирования.
type PaymentMetadata struct {
	BookingID string `json:"booking_id" example:"1"` // Уникальный идентификатор бронирования
}
