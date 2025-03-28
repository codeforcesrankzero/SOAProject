// handlers/user_handler_test.go
package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"user-service/models"
	"user-service/services"

	"github.com/gin-gonic/gin"
)

type MockUserService struct {
    RegisterFunc func(models.RegisterRequest) error
    LoginFunc func(models.LoginRequest) (string, error)
    GetUserProfileFunc func(uint) (*models.User, error)
    UpdateUserProfileFunc func(uint, models.UpdateProfileRequest) error
    ValidateTokenFunc func(string) (uint, error)
}

var _ services.UserServiceInterface = (*MockUserService)(nil)

func (m *MockUserService) Register(req models.RegisterRequest) error {
    return m.RegisterFunc(req)
}

func (m *MockUserService) Login(req models.LoginRequest) (string, error) {
    return m.LoginFunc(req)
}

func (m *MockUserService) GetUserProfile(userID uint) (*models.User, error) {
    return m.GetUserProfileFunc(userID)
}

func (m *MockUserService) UpdateUserProfile(userID uint, req models.UpdateProfileRequest) error {
    return m.UpdateUserProfileFunc(userID, req)
}

func (m *MockUserService) ValidateToken(token string) (uint, error) {
    return m.ValidateTokenFunc(token)
}


func TestRegisterHandler(t *testing.T) {
    gin.SetMode(gin.TestMode)
    r := gin.Default()
    
    mockService := &MockUserService{
        RegisterFunc: func(req models.RegisterRequest) error {
            return nil
        },
    }
    
    handler := NewUserHandler(mockService)
    
    r.POST("/register", handler.Register)
    
    reqBody, _ := json.Marshal(models.RegisterRequest{
        Login:    "testuser",
        Password: "password123",
        Email:    "test@example.com",
    })
    req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(reqBody))
    req.Header.Set("Content-Type", "application/json")
    
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)
    
    if w.Code != http.StatusCreated {
        t.Errorf("Ожидается код 201, получен: %d", w.Code)
    }
    
    var response map[string]string
    json.Unmarshal(w.Body.Bytes(), &response)
    if message, exists := response["message"]; !exists || message != "Пользователь успешно зарегистрирован" {
        t.Errorf("Ожидается сообщение об успешной регистрации, получено: %v", response)
    }
}

func TestLoginHandler(t *testing.T) {
    gin.SetMode(gin.TestMode)
    r := gin.Default()
    
    mockService := &MockUserService{
        LoginFunc: func(req models.LoginRequest) (string, error) {
            return "test_token", nil
        },
    }
    
    handler := NewUserHandler(mockService)
    
    r.POST("/login", handler.Login)
    
    reqBody, _ := json.Marshal(models.LoginRequest{
        Login:    "testuser",
        Password: "password123",
    })
    req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(reqBody))
    req.Header.Set("Content-Type", "application/json")
    
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)
    
    if w.Code != http.StatusOK {
        t.Errorf("Ожидается код 200, получен: %d", w.Code)
    }
    
    var response models.LoginResponse
    json.Unmarshal(w.Body.Bytes(), &response)
    if response.Token != "test_token" {
        t.Errorf("Ожидается токен test_token, получен: %s", response.Token)
    }
}
