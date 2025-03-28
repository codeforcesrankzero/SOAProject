package repository

import (
	"context"
	"promocodes-service/internal/model"
	"testing"
)

type MockPromocodeRepository struct {
    promocodesByCode map[string]*model.Promocode
    promocodesByID   map[int64]*model.Promocode
    idCounter        int64
}

func NewMockPromocodeRepository() *MockPromocodeRepository {
    return &MockPromocodeRepository{
        promocodesByCode: make(map[string]*model.Promocode),
        promocodesByID:   make(map[int64]*model.Promocode),
        idCounter:        1,
    }
}

func (r *MockPromocodeRepository) CreatePromocode(ctx context.Context, promocode *model.Promocode) error {
    promocode.ID = r.idCounter
    r.idCounter++
    r.promocodesByCode[promocode.Code] = promocode
    r.promocodesByID[promocode.ID] = promocode
    return nil
}

func (r *MockPromocodeRepository) GetPromocodeByID(ctx context.Context, id int64) (*model.Promocode, error) {
    promocode, exists := r.promocodesByID[id]
    if !exists {
        return nil, nil
    }
    return promocode, nil
}

func (r *MockPromocodeRepository) GetPromocodeByCode(ctx context.Context, code string) (*model.Promocode, error) {
    promocode, exists := r.promocodesByCode[code]
    if !exists {
        return nil, nil
    }
    return promocode, nil
}

func TestCreatePromocode(t *testing.T) {
    repo := NewMockPromocodeRepository()

    promocode := &model.Promocode{
        Code: "TEST123",
    }

    err := repo.CreatePromocode(context.Background(), promocode)
    if err != nil {
        t.Errorf("Expected no error, got %v", err)
    }
    if promocode.ID != 1 {
        t.Errorf("Expected ID=1, got %d", promocode.ID)
    }

    savedPromocode, err := repo.GetPromocodeByCode(context.Background(), "TEST123")
    if err != nil {
        t.Errorf("Error fetching promocode: %v", err)
    }
    if savedPromocode == nil {
        t.Error("Expected to find promocode, got nil")
    }
}

func TestGetPromocodeByID(t *testing.T) {
    repo := NewMockPromocodeRepository()

    testPromocode := &model.Promocode{
        Code: "TEST123",
    }
    repo.CreatePromocode(context.Background(), testPromocode)

    promocode, err := repo.GetPromocodeByID(context.Background(), 1)
    if err != nil {
        t.Errorf("Expected successful retrieval, got: %v", err)
    }
    if promocode == nil {
        t.Error("Promocode should be found")
    }

    promocode, err = repo.GetPromocodeByID(context.Background(), 999)
    if err != nil {
        t.Errorf("Expected nil without error, got: %v", err)
    }
    if promocode != nil {
        t.Error("Promocode should not be found")
    }
}
