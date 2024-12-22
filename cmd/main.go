package main

import (
	"hotel-booking/internal/auth"
	"hotel-booking/internal/hotels"
	"hotel-booking/internal/storage"
	"hotel-booking/internal/users"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Загружаем .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}

	// Подключение базы данных
	storage.ConnectDatabase()

	// Выполнение миграций
	err = storage.DB.AutoMigrate(&users.User{}, &hotels.Hotel{}, &hotels.Room{})
	if err != nil {
		log.Fatal("Ошибка миграции:", err)
	}

	r := gin.Default()

	r.POST("/auth/register", auth.RegisterHandler)
	r.POST("/auth/login", auth.LoginHandler)

	authorized := r.Group("/")
	authorized.Use(auth.AuthMiddleware())
	// Добавляй защищённые маршруты здесь

	if err := r.Run(":8080"); err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}
