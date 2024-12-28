package email

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

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
