package email

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Отправка тестового письма
// @Description Отправляет тестовое письмо на указанную почту.
// @Tags email
// @Accept json
// @Produce json
// @Param to query string true "Email получателя"
// @Success 200 {object} response.SuccessResponse "Тестовое письмо успешно отправлено"
// @Failure 400 {object} response.ErrorResponse "Не указан email получателя"
// @Failure 500 {object} response.ErrorResponse "Ошибка отправки письма"
// @Router /email/test [get]
func SendTestEmailHandler(c *gin.Context) {
	to := c.Query("to")
	if to == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Не указан email получателя"})
		return
	}

	subject := "Тестовое письмо"
	body := "Проверка работы отправки сообщений"

	if err := SendEmail(to, subject, body); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка отправки письма", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Тестовое письмо успешно отправлено"})

}
