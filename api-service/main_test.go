// api-service/main_test.go
package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestProxyRequest(t *testing.T) {
    mockUserService := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"message":"Hello from user service"}`))
    }))
    defer mockUserService.Close()

    gin.SetMode(gin.TestMode)
    r := gin.Default()
    
    r.GET("/test", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"message": "Hello from user service"})
    })
    
    req, _ := http.NewRequest("GET", "/test", nil)
    w := httptest.NewRecorder()
    
    r.ServeHTTP(w, req)
    
    if w.Code != http.StatusOK {
        t.Errorf("Ожидается код 200, получен: %d", w.Code)
    }
    
    expectedBody := `{"message":"Hello from user service"}`
    if w.Body.String() != expectedBody {
        t.Errorf("Ожидается: %s, получено: %s", expectedBody, w.Body.String())
    }
}
