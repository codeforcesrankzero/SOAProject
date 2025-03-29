package repository

import (
	"context"
	"promocodes-service/internal/model"

	"github.com/jinzhu/gorm"
)


type PromocodeRepository interface {
	CreatePromocode(ctx context.Context, promocode *model.Promocode) error
	GetPromocodeByID(ctx context.Context, id int64) (*model.Promocode, error)
	UpdatePromocode(ctx context.Context, promocode *model.Promocode) error
	DeletePromocode(ctx context.Context, id int64) error
	ListPromocodes(ctx context.Context, page, perPage int) ([]*model.Promocode, int, error)
}

type promocodeRepository struct {
	db *gorm.DB
}

func NewPromocodeRepository(db *gorm.DB) PromocodeRepository {
	return &promocodeRepository{db: db}
}
func (r *promocodeRepository) CreatePromocode(ctx context.Context, promocode *model.Promocode) error {
	return r.db.Create(promocode).Error
}

func (r *promocodeRepository) GetPromocodeByID(ctx context.Context, id int64) (*model.Promocode, error) {
	var promocode model.Promocode
	err := r.db.First(&promocode, id).Error
	if err != nil {
		return nil, err
	}
	return &promocode, nil
}

func (r *promocodeRepository) UpdatePromocode(ctx context.Context, promocode *model.Promocode) error {
	return r.db.Save(promocode).Error
}

func (r *promocodeRepository) DeletePromocode(ctx context.Context, id int64) error {
	return r.db.Delete(&model.Promocode{}, id).Error
}

func (r *promocodeRepository) ListPromocodes(ctx context.Context, page, perPage int) ([]*model.Promocode, int, error) {
	var promocodes []*model.Promocode
	var count int64
	
	offset := (page - 1) * perPage
	
	err := r.db.Model(&model.Promocode{}).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	
	err = r.db.Offset(offset).Limit(perPage).Find(&promocodes).Error
	if err != nil {
		return nil, 0, err
	}
	
	return promocodes, int(count), nil
}
