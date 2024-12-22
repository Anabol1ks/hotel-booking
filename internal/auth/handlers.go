package auth

import (
	"hotel-booking/internal/storage"
	"hotel-booking/internal/users"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type RegisterInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
}

func RegisterHandler(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверка уникальности почты и телефона
	var existingUser users.User
	if err := storage.DB.Where("email = ? or phone = ?", input.Email, input.Phone).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Почта или телефон уже зарегистрированы"})
		return
	}

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось хешировать пароль"})
		return
	}

	// Создаём пользователя
	user := users.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
		Phone:    input.Phone,
		Role:     "client",
	}

	if err := storage.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось создать пользователя"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Регистрация успешна"})
}

var jwtSecret = []byte(os.Getenv("JWT_KEY"))

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func LoginHandler(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверяем пользователя
	var user users.User
	if err := storage.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный email или пароль"})
		return
	}

	// Проверяем пароль
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный email или пароль"})
		return
	}

	// Генерируем JWT
	token := GenerateJWT(user.ID, user.Role)
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func GenerateJWT(userID uint, role string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(72 * time.Hour).Unix(),
	})
	tokenString, _ := token.SignedString(jwtSecret)
	return tokenString
}
