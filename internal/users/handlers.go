package users

import (
	"hotel-booking/internal/storage"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UpdateRoleInput struct {
	Role string `json:"role" binding:"required"`
}

func UpdateRoleHandler(c *gin.Context) {
	// Проверка роли администратора
	role := c.GetString("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Только администратор может изменять роли"})
		return
	}

	// Получаем ID пользователя
	userID := c.Param("id")

	var input UpdateRoleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверяем наличие пользователя
	var user User
	if err := storage.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
		return
	}

	// Обновляем роль
	user.Role = input.Role
	if err := storage.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось обновить роль"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Роль успешно обновлена"})
}

// Управление пользователями для администратора
func GetUsersHandler(c *gin.Context) {
	role := c.GetString("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Только администратор может просматривать пользователей"})
		return
	}

	var users []User
	if err := storage.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении пользователей"})
		return
	}

	c.JSON(http.StatusOK, users)

}
