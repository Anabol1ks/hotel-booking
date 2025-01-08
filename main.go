package main

import (
	_ "hotel-booking/docs"
	"hotel-booking/internal/auth"
	"hotel-booking/internal/bookings"
	"hotel-booking/internal/email"
	"hotel-booking/internal/hotels"
	"hotel-booking/internal/payments"
	"hotel-booking/internal/storage"
	"hotel-booking/internal/users"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// @Title Система бронирования номеров
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	key := os.Getenv("TEST_ENV")
	if key == "" {
		log.Println("\nПеременной среды нет, используется .env")
		// Загружаем .env
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Ошибка загрузки .env файла")
		}
	}

	// Подключение базы данных
	storage.ConnectDatabase()

	// Выполнение миграций
	err := storage.DB.AutoMigrate(&users.User{}, &hotels.Favorite{}, &hotels.Hotel{}, &hotels.Room{}, &bookings.Booking{})
	if err != nil {
		log.Fatal("Ошибка миграции:", err)
	}

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Укажи адрес фронтенда React
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.POST("/auth/register", auth.RegisterHandler)
	r.POST("/auth/login", auth.LoginHandler)

	r.GET("/hotels", hotels.GetHotelsHandler)
	r.GET("/rooms", hotels.GetRoomsHandler)
	r.GET("/rooms/:id/bookings", bookings.GetRoomBookingsHandler)

	r.GET("/email/test", email.SendTestEmailHandler)
	r.POST("/auth/reset-password-request", auth.ResetPasswordRequestHandler)
	r.POST("/auth/reset-password", auth.ResetPasswordHandler)

	authorized := r.Group("/")
	{
		authorized.Use(auth.AuthMiddleware())
		authorized.POST("/bookings", bookings.CreateBookingHandler)
		authorized.POST("/bookings/:id/pay", payments.CreatePaymentHandler)
		authorized.DELETE("/bookings/:id", bookings.CancelBookingHandler)
		authorized.POST("/bookings/:id/refund", payments.RefundPaymentHandler)
		authorized.POST("/favorites/:room_id", hotels.AddToFavoritesHandler)
		authorized.GET("/favorites", hotels.GetFavoritesHandler)
		authorized.DELETE("/favorites/:room_id", hotels.RemoveFromFavoritesHandler)
		authorized.POST("/booking/offline", bookings.CreateOfflineBookingHandler)
	}
	r.POST("/payments/callback", payments.PaymentCallbackHandler)

	owners := authorized.Group("/owners")
	{
		owners.POST("/hotels", hotels.CreateHotelHandler)
		owners.POST("/rooms", hotels.CreateRoomHandler)
		owners.GET("/hotels", hotels.GetOwnerHotelsHandler)
		owners.GET("/bookings", bookings.GetOwnerBookingsHandler)
		owners.PUT("/:id/room", hotels.ChangeRoomHandler)
		owners.DELETE("/:id/room", hotels.DeleteRoomHandler)
	}

	admins := authorized.Group("/admin")
	{
		admins.GET("/users", users.GetUsersHandler)
		admins.PUT("/users/:id/role", users.UpdateRoleHandler)
	}

	if err := r.Run(":8080"); err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}
