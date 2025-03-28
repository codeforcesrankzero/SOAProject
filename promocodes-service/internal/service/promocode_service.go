package service

import (
	"context"
	"errors"
	"promocodes-service/internal/model"
	"promocodes-service/internal/repository"
)

type PromocodeService interface {
	CreatePromocode(ctx context.Context, promocode *model.Promocode) error
	GetPromocode(ctx context.Context, id int64) (*model.Promocode, error)
	UpdatePromocode(ctx context.Context, promocode *model.Promocode) error
	DeletePromocode(ctx context.Context, id int64, creatorID int64) (model.Promocode, error)
	ListPromocodes(ctx context.Context, page, perPage int) ([]*model.Promocode, int, error)
}

type promocodeService struct {
	repo repository.PromocodeRepository
}

func NewPromocodeService(repo repository.PromocodeRepository) PromocodeService {
	return &promocodeService{repo: repo}
}

func (s *promocodeService) CreatePromocode(ctx context.Context, promocode *model.Promocode) error {
	return s.repo.CreatePromocode(ctx, promocode)
}

func (s *promocodeService) GetPromocode(ctx context.Context, id int64) (*model.Promocode, error) {
	return s.repo.GetPromocodeByID(ctx, id)
}

func (s *promocodeService) UpdatePromocode(ctx context.Context, promocode *model.Promocode) error {
	existingPromocode, err := s.repo.GetPromocodeByID(ctx, promocode.ID)
	if err != nil {
		return err
	}
	
	if existingPromocode.CreatorID != promocode.CreatorID {
		return errors.New("unauthorized: you can only update your own promocodes")
	}
	
	return s.repo.UpdatePromocode(ctx, promocode)
}

func (s *promocodeService) DeletePromocode(ctx context.Context, id int64, creatorID int64) (model.Promocode, error) {
	existingPromocode, err := s.repo.GetPromocodeByID(ctx, id)
	if err != nil {
		return model.Promocode{}, err
	}
	
	if existingPromocode.CreatorID != creatorID {
		return model.Promocode{}, errors.New("unauthorized: you can only delete your own promocodes")
	}
	
	return *existingPromocode, s.repo.DeletePromocode(ctx, id)
}

func (s *promocodeService) ListPromocodes(ctx context.Context, page, perPage int) ([]*model.Promocode, int, error) {
	return s.repo.ListPromocodes(ctx, page, perPage)
}
