package main

import (
	"log"
	"os"
	"promocodes-service/api"
	"promocodes-service/pkg/grpc"

	"github.com/gin-gonic/gin"
)

func main() {
    jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "super_secret_key"
	}

    grpcServerAddr := os.Getenv("GRPC_SERVER_ADDR")
    if grpcServerAddr == "" {
        grpcServerAddr = "promocode-grpc:50051"
    }

    port := os.Getenv("PORT")
    if port == "" {
        port = "8082"
    }

    grpcClient, err := grpc.NewClient(grpcServerAddr)
    if err != nil {
        log.Fatalf("Failed to connect to gRPC server: %v", err)
    }
    defer grpcClient.Close()

    router := gin.Default()

    handler := api.NewHandler(grpcClient, jwtSecret)
    api.SetupRoutes(router, handler)

    if err := router.Run(":" + port); err != nil {
        log.Fatalf("Failed to run server: %v", err)
    }
}
