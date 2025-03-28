package models

import (
	"time"
)

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Login     string    `json:"login" gorm:"unique;not null"`
	Password  string    `json:"-" gorm:"not null"` // Пароль не возвращается в JSON
	Email     string    `json:"email" gorm:"unique;not null"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	BirthDate time.Time `json:"birth_date"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type RegisterRequest struct {
	Login    string `json:"login" binding:"required,min=4,max=20"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email" binding:"required,email"`
}

type LoginRequest struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UpdateProfileRequest struct {
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	BirthDate *time.Time `json:"birth_date"`
	Email     string     `json:"email" binding:"email"`
	Phone     string     `json:"phone"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
