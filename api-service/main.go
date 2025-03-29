package main

import (
	"log"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	userServiceURL := os.Getenv("USER_SERVICE_URL")
	if userServiceURL == "" {
		userServiceURL = "http://user-service:8081"
	}

	promoServiceURL := os.Getenv("PROMO_SERVICE_URL")
	if promoServiceURL == "" {
		promoServiceURL = "http://promocode-api:8082" // Предполагаемый адрес сервиса промокодов
	}

	userTarget, err := url.Parse(userServiceURL)
	if err != nil {
		log.Fatalf("Ошибка при парсинге URL сервиса пользователей: %v", err)
	}

	promoTarget, err := url.Parse(promoServiceURL)
	if err != nil {
		log.Fatalf("Ошибка при парсинге URL сервиса промокодов: %v", err)
	}

	userProxy := httputil.NewSingleHostReverseProxy(userTarget)
	promoProxy := httputil.NewSingleHostReverseProxy(promoTarget)

	r.Any("/*path", func(c *gin.Context) {
		path := c.Param("path")
		log.Printf("API Gateway получил запрос: %s %s", c.Request.Method, path)
		if strings.HasPrefix(path, "/promo") || strings.HasPrefix(path, "/api/promocodes") {
			log.Printf("Проксирование запроса на сервис промокодов: %s", path)
			promoProxy.ServeHTTP(c.Writer, c.Request)
		} else if strings.HasPrefix(path, "/profile") || strings.HasPrefix(path, "/login") || strings.HasPrefix(path, "/register") {
			log.Printf("Проксирование запроса на сервис пользователей: %s", path)
			userProxy.ServeHTTP(c.Writer, c.Request)
		}
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("API Gateway запущен на порту %s", port)
	log.Fatal(r.Run(":" + port))
}
