package main

import (
	"hotel-booking/internal/auth"
	"hotel-booking/internal/bookings"
	"hotel-booking/internal/hotels"
	"hotel-booking/internal/payments"
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
	err = storage.DB.AutoMigrate(&users.User{}, &hotels.Hotel{}, &hotels.Room{}, &bookings.Booking{})
	if err != nil {
		log.Fatal("Ошибка миграции:", err)
	}

	r := gin.Default()

	r.POST("/auth/register", auth.RegisterHandler)
	r.POST("/auth/login", auth.LoginHandler)

	authorized := r.Group("/")
	authorized.Use(auth.AuthMiddleware())
	// смена роли
	authorized.PUT("/users/:id/role", users.UpdateRoleHandler)

	authorized.POST("/hotels", hotels.CreateHotelHandler)
	authorized.POST("/rooms", hotels.CreateRoomHandler)

	authorized.POST("/bookings", bookings.CreateBookingHandler)

	r.GET("/hotels", hotels.GetHotelsHandler)
	r.GET("/rooms", hotels.GetRoomsHandler)
	r.GET("/rooms/:id/bookings", bookings.GetRoomBookingsHandler)

	authorized.POST("/bookings/:id/pay", payments.CreatePaymentHandler)
	r.POST("/payments/callback", payments.PaymentCallbackHandler)

	owners := authorized.Group("/owners")
	owners.GET("/hotels", hotels.GetOwnerHotelsHandler)
	owners.GET("/bookings", bookings.GetOwnerBookingsHandler)

	admins := authorized.Group("/admin")
	admins.GET("/users", users.GetUsersHandler)
	admins.PUT("/users/:id/role", users.UpdateRoleHandler)

	if err := r.Run(":8080"); err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}
