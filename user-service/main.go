package main

import (
	"log"
	"os"
	"time"

	"user-service/handlers"
	"user-service/models"
	"user-service/repository"
	"user-service/services"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=" + os.Getenv("DB_HOST") +
		" user=" + os.Getenv("DB_USER") +
		" password=" + os.Getenv("DB_PASSWORD") +
		" dbname=" + os.Getenv("DB_NAME") +
		" port=" + os.Getenv("DB_PORT") +
		" sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}

	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("Ошибка миграции базы данных: %v", err)
	}

	userRepo := repository.NewUserRepository(db)

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "my_secret_key"
	}
	userService := services.NewUserService(userRepo, jwtSecret, 24*time.Hour)

	userHandler := handlers.NewUserHandler(userService)

	r := gin.Default()

	r.POST("/register", userHandler.Register)
	r.POST("/login", userHandler.Login)

	protected := r.Group("/")
	protected.Use(userHandler.AuthMiddleware())
	{
		protected.GET("/profile", userHandler.GetProfile)
		protected.PUT("/profile", userHandler.UpdateProfile)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	log.Fatal(r.Run(":" + port))
}
