package model

import (
	"time"
)

type Promocode struct {
	ID          int64     `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"size:100;not null"`
	Description string    `json:"description" gorm:"size:255"`
	CreatorID   int64     `json:"creator_id" gorm:"not null"`
	Discount    float32   `json:"discount" gorm:"not null"`
	Code        string    `json:"code" gorm:"size:50;uniqueIndex;not null"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type CreatePromocodeRequest struct {
	Name        string  `json:"name" binding:"required,min=3,max=100"`
	Description string  `json:"description" binding:"max=255"`
	Discount    float32 `json:"discount" binding:"required,min=0,max=100"`
	Code        string  `json:"code" binding:"required,min=3,max=50"`
}

type UpdatePromocodeRequest struct {
	Name        string  `json:"name" binding:"required,min=3,max=100"`
	Description string  `json:"description" binding:"max=255"`
	Discount    float32 `json:"discount" binding:"required,min=0,max=100"`
	Code        string  `json:"code" binding:"required,min=3,max=50"`
}

type PaginationRequest struct {
	Page    int `json:"page" form:"page" binding:"min=1"`
	PerPage int `json:"per_page" form:"per_page" binding:"min=1,max=100"`
}

type PaginatedResponse struct {
	Items   interface{} `json:"items"`
	Total   int         `json:"total"`
	Page    int         `json:"page"`
	PerPage int         `json:"per_page"`
}
