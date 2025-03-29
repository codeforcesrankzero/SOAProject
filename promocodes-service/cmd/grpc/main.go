package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"promocodes-service/internal/model"
	"promocodes-service/internal/repository"
	"promocodes-service/internal/service"
	grpcHandler "promocodes-service/pkg/grpc"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"google.golang.org/grpc"
)

func main() {
    dbHost := getEnv("DB_HOST", "postgres")
    dbPort := getEnv("DB_PORT", "5432")
    dbUser := getEnv("DB_USER", "postgres")
    dbPassword := getEnv("DB_PASSWORD", "postgres")
    dbName := getEnv("DB_NAME", "promocodedb")
    
    grpcPort := getEnv("GRPC_PORT", "50051")

    dbURI := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
        dbHost, dbPort, dbUser, dbName, dbPassword)

    var db *gorm.DB
    var err error
    
    maxRetries := 5
    for i := 0; i < maxRetries; i++ {
        db, err = gorm.Open("postgres", dbURI)
        if err == nil {
            break
        }
        log.Printf("Failed to connect to database (attempt %d/%d): %v", i+1, maxRetries, err)
        time.Sleep(time.Second * 3)
    }
    
    if err != nil {
        log.Fatalf("Could not connect to database after %d attempts: %v", maxRetries, err)
    }
    
    defer db.Close()

    db.AutoMigrate(&model.Promocode{})

    promocodeRepo := repository.NewPromocodeRepository(db)
    promocodeService := service.NewPromocodeService(promocodeRepo)

    lis, err := net.Listen("tcp", ":"+grpcPort)
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }

    grpcServer := grpc.NewServer()
    grpcHandler.RegisterPromocodeServiceServer(grpcServer, grpcHandler.NewPromocodeServer(promocodeService))

    log.Printf("Starting gRPC server on :%s", grpcPort)
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("Failed to serve: %v", err)
    }
}

func getEnv(key, defaultValue string) string {
    value := os.Getenv(key)
    if value == "" {
        return defaultValue
    }
    return value
}
