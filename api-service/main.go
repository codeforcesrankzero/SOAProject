package main

import (
	"log"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	userServiceURL := os.Getenv("USER_SERVICE_URL")
	if userServiceURL == "" {
		userServiceURL = "http://user-service:8081"
	}

	target, err := url.Parse(userServiceURL)
	if err != nil {
		log.Fatalf("Ошибка при парсинге URL: %v", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(target)

	r.Any("/*path", func(c *gin.Context) {
		log.Printf("API Gateway получил запрос: %s %s", c.Request.Method, c.Request.URL.Path)
		
		proxy.ServeHTTP(c.Writer, c.Request)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("API Gateway запущен на порту %s", port)
	log.Fatal(r.Run(":" + port))
}
