package service

import (
	"context"
	"promocodes-service/internal/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPromocodeRepository struct {
	mock.Mock
}

func (m *MockPromocodeRepository) CreatePromocode(ctx context.Context, promocode *model.Promocode) error {
	args := m.Called(ctx, promocode)
	return args.Error(0)
}

func (m *MockPromocodeRepository) GetPromocodeByID(ctx context.Context, id int64) (*model.Promocode, error) {
	args := m.Called(ctx, id)
	if promocode, ok := args.Get(0).(*model.Promocode); ok {
		return promocode, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockPromocodeRepository) UpdatePromocode(ctx context.Context, promocode *model.Promocode) error {
	args := m.Called(ctx, promocode)
	return args.Error(0)
}

func (m *MockPromocodeRepository) DeletePromocode(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPromocodeRepository) ListPromocodes(ctx context.Context, page, perPage int) ([]*model.Promocode, int, error) {
	args := m.Called(ctx, page, perPage)
	if promocodes, ok := args.Get(0).([]*model.Promocode); ok {
		return promocodes, args.Int(1), args.Error(2)
	}
	return nil, 0, args.Error(2)
}

func TestCreatePromocode(t *testing.T) {
	mockRepo := new(MockPromocodeRepository)
	service := NewPromocodeService(mockRepo)

	promocode := &model.Promocode{ID: 1, CreatorID: 123, Code: "DISCOUNT10"}

	mockRepo.On("CreatePromocode", mock.Anything, promocode).Return(nil)

	err := service.CreatePromocode(context.Background(), promocode)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestGetPromocode(t *testing.T) {
	mockRepo := new(MockPromocodeRepository)
	service := NewPromocodeService(mockRepo)

	expectedPromocode := &model.Promocode{ID: 1, CreatorID: 123, Code: "DISCOUNT10"}

	mockRepo.On("GetPromocodeByID", mock.Anything, int64(1)).Return(expectedPromocode, nil)

	promocode, err := service.GetPromocode(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, expectedPromocode, promocode)
	mockRepo.AssertExpectations(t)
}

func TestUpdatePromocode(t *testing.T) {
	mockRepo := new(MockPromocodeRepository)
	service := NewPromocodeService(mockRepo)

	existingPromocode := &model.Promocode{ID: 1, CreatorID: 123, Code: "OLD"}
	updatedPromocode := &model.Promocode{ID: 1, CreatorID: 123, Code: "NEW"}

	mockRepo.On("GetPromocodeByID", mock.Anything, int64(1)).Return(existingPromocode, nil)
	mockRepo.On("UpdatePromocode", mock.Anything, updatedPromocode).Return(nil)

	err := service.UpdatePromocode(context.Background(), updatedPromocode)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeletePromocode_Unauthorized(t *testing.T) {
	mockRepo := new(MockPromocodeRepository)
	service := NewPromocodeService(mockRepo)

	existingPromocode := &model.Promocode{ID: 1, CreatorID: 123}

	mockRepo.On("GetPromocodeByID", mock.Anything, int64(1)).Return(existingPromocode, nil)

	_, err := service.DeletePromocode(context.Background(), 1, 456)
	assert.Error(t, err)
	assert.Equal(t, "unauthorized: you can only delete your own promocodes", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestListPromocodes(t *testing.T) {
	mockRepo := new(MockPromocodeRepository)
	service := NewPromocodeService(mockRepo)

	promocodes := []*model.Promocode{
		{ID: 1, CreatorID: 123, Code: "PROMO1"},
		{ID: 2, CreatorID: 123, Code: "PROMO2"},
	}

	mockRepo.On("ListPromocodes", mock.Anything, 1, 2).Return(promocodes, 2, nil)

	result, count, err := service.ListPromocodes(context.Background(), 1, 2)
	assert.NoError(t, err)
	assert.Equal(t, 2, count)
	assert.Equal(t, promocodes, result)
	mockRepo.AssertExpectations(t)
}
