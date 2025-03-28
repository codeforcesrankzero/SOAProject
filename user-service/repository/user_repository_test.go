// repository/user_repository_test.go
package repository

import (
	"testing"
	"user-service/models"
)

type MockUserRepository struct {
    users     map[string]*models.User
    usersById map[uint]*models.User
    idCounter uint
}

var _ UserRepositoryInterface = (*MockUserRepository)(nil)

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

func TestCreateUser(t *testing.T) {
    repo := NewMockUserRepository()
    
    user := &models.User{
        Login: "testuser",
        Password: "hashedpassword",
        Email: "test@example.com",
    }
    
    err := repo.CreateUser(user)
    if err != nil {
        t.Errorf("Ожидается успешное создание пользователя, получено: %v", err)
    }
    if user.ID != 1 {
        t.Errorf("Ожидается ID=1, получено: %d", user.ID)
    }
    
    savedUser, err := repo.GetUserByLogin("testuser")
    if err != nil {
        t.Errorf("Ошибка при получении пользователя: %v", err)
    }
    if savedUser == nil {
        t.Error("Ожидается найденный пользователь, получен nil")
    }
}

func TestGetUserByLogin(t *testing.T) {
    repo := NewMockUserRepository()
    
    testUser := &models.User{
        Login: "testuser",
        Password: "hashedpassword",
        Email: "test@example.com",
    }
    repo.CreateUser(testUser)
    
    user, err := repo.GetUserByLogin("testuser")
    if err != nil {
        t.Errorf("Ожидается успешный поиск, получено: %v", err)
    }
    if user == nil {
        t.Error("Пользователь должен быть найден")
    }
    
    user, err = repo.GetUserByLogin("nonexistent")
    if err != nil {
        t.Errorf("Ожидается nil без ошибки, получено: %v", err)
    }
    if user != nil {
        t.Error("Пользователь не должен быть найден")
    }
}
