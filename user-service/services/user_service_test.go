package services

import (
	"testing"
	"time"
	"user-service/models"
	"user-service/repository" // Импортируем пакет repository

	"golang.org/x/crypto/bcrypt"
)

type MockUserRepository struct {
    users     map[string]*models.User
    usersById map[uint]*models.User
    idCounter uint
}

var _ repository.UserRepositoryInterface = (*MockUserRepository)(nil)

func NewMockUserRepository() *MockUserRepository {
    return &MockUserRepository{
        users:     make(map[string]*models.User),
        usersById: make(map[uint]*models.User),
        idCounter: 1,
    }
}

func (r *MockUserRepository) CreateUser(user *models.User) error {
    user.ID = r.idCounter
    r.idCounter++
    r.users[user.Login] = user
    r.usersById[user.ID] = user
    return nil
}

func (r *MockUserRepository) GetUserByLogin(login string) (*models.User, error) {
    user, exists := r.users[login]
    if !exists {
        return nil, nil
    }
    return user, nil
}

func (r *MockUserRepository) GetUserByID(id uint) (*models.User, error) {
    user, exists := r.usersById[id]
    if !exists {
        return nil, nil
    }
    return user, nil
}

func (r *MockUserRepository) UpdateUser(user *models.User) error {
    r.users[user.Login] = user
    r.usersById[user.ID] = user
    return nil
}

func TestRegister(t *testing.T) {
    mockRepo := NewMockUserRepository()
    service := NewUserService(mockRepo, "test_secret", 24*time.Hour)

    err := service.Register(models.RegisterRequest{
        Login:    "testuser",
        Password: "password123",
        Email:    "test@example.com",
    })

    if err != nil {
        t.Errorf("Ожидается успешная регистрация, получена ошибка: %v", err)
    }

    user, _ := mockRepo.GetUserByLogin("testuser")
    if user == nil {
        t.Error("Пользователь должен быть создан")
    }
    if user.Email != "test@example.com" {
        t.Errorf("Ожидается email test@example.com, получен: %s", user.Email)
    }

    err = service.Register(models.RegisterRequest{
        Login:    "testuser",
        Password: "anotherpassword",
        Email:    "another@example.com",
    })

    if err == nil {
        t.Error("Ожидается ошибка дубликата логина")
    }
}

func TestLogin(t *testing.T) {
    mockRepo := NewMockUserRepository()
    service := NewUserService(mockRepo, "test_secret", 24*time.Hour)

    hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
    mockRepo.users["testuser"] = &models.User{
        ID:       1,
        Login:    "testuser",
        Password: string(hashedPassword),
        Email:    "test@example.com",
    }
    mockRepo.usersById[1] = mockRepo.users["testuser"]

    token, err := service.Login(models.LoginRequest{
        Login:    "testuser",
        Password: "password123",
    })

    if err != nil {
        t.Errorf("Ожидается успешный вход, получена ошибка: %v", err)
    }
    if token == "" {
        t.Error("Токен не должен быть пустым")
    }

    _, err = service.Login(models.LoginRequest{
        Login:    "testuser",
        Password: "wrongpassword",
    })

    if err == nil {
        t.Error("Ожидается ошибка неверного пароля")
    }
}
