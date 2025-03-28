// repository/interfaces.go
package repository

import (
	"user-service/models"
)

type UserRepositoryInterface interface {
    CreateUser(user *models.User) error
    GetUserByLogin(login string) (*models.User, error)
    GetUserByID(id uint) (*models.User, error)
    UpdateUser(user *models.User) error
}
