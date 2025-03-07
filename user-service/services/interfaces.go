// services/interfaces.go
package services

import (
	"user-service/models"
)

type UserServiceInterface interface {
    Register(req models.RegisterRequest) error
    Login(req models.LoginRequest) (string, error)
    GetUserProfile(userID uint) (*models.User, error)
    UpdateUserProfile(userID uint, req models.UpdateProfileRequest) error
    ValidateToken(tokenString string) (uint, error)
}