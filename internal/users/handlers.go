package users

import (
	"hotel-booking/internal/storage"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UpdateRoleInput struct {
	Role string `json:"role" binding:"required" enums:"owner,admin,client"`
}

// @Security BearerAuth
// UpdateRoleHandler godoc
// @Summary Обновление роли пользователя
// @Description Обновление роли пользователя через панель администратора
// @Tags admin
// @Accept json
// @Produce json
// @Param input body UpdateRoleInput true "Данные для обновления роли пользователя. Возможные значения: owner, admin, client"
// @Param id path int true "ID пользователя"
// @Success 200 {object} response.SuccessResponse "Роль успешно обновлена"
// @Failure 400 {object} response.ErrorResponse "Описание ошибки валидации"
// @Failure 403 {object} response.ErrorResponse "Только администратор может изменять роли"
// @Failure 404 {object} response.ErrorResponse "Пользователь не найден"
// @Failure 500 {object} response.ErrorResponse "Не удалось обновить роль"
// @Router /admin/users/:id/role [put]
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
	validRoles := map[string]bool{"owner": true, "client": true, "admin": true}
	if !validRoles[input.Role] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверная роль"})
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

// @Security BearerAuth
// GetUsersHandler godoc
// @Summary Получение списка пользователей
// @Description Получение списка пользователей через панель администратора
// @Tags admin
// @Accept json
// @Produce json
// @Success 200 {array} response.UserResponse "Список пользователей"
// @Failure 403 {object} response.ErrorResponse "Только администратор может просматривать пользователей"
// @Failure 500 {object} response.ErrorResponse "Ошибка при получении пользователей"
// @Router /admin/users [get]
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
