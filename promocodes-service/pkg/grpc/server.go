package grpc

import (
	"context"
	"promocodes-service/internal/model"
	"promocodes-service/internal/service"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PromocodeServer struct {
	UnimplementedPromocodeServiceServer
	service service.PromocodeService
}

func NewPromocodeServer(service service.PromocodeService) *PromocodeServer {
	return &PromocodeServer{service: service}
}

func (s *PromocodeServer) CreatePromocode(ctx context.Context, req *CreatePromocodeRequest) (*PromocodeResponse, error) {
	promocode := &model.Promocode{
		Name:        req.Name,
		Description: req.Description,
		CreatorID:   req.CreatorId,
		Discount:    req.Discount,
		Code:        req.Code,
	}

	err := s.service.CreatePromocode(ctx, promocode)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &PromocodeResponse{
		Id:          promocode.ID,
		Name:        promocode.Name,
		Description: promocode.Description,
		CreatorId:   promocode.CreatorID,
		Discount:    promocode.Discount,
		Code:        promocode.Code,
		CreatedAt:   promocode.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   promocode.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (s *PromocodeServer) GetPromocode(ctx context.Context, req *GetPromocodeRequest) (*PromocodeResponse, error) {
	promocode, err := s.service.GetPromocode(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "promocode not found")
	}

	return &PromocodeResponse{
		Id:          promocode.ID,
		Name:        promocode.Name,
		Description: promocode.Description,
		CreatorId:   promocode.CreatorID,
		Discount:    promocode.Discount,
		Code:        promocode.Code,
		CreatedAt:   promocode.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   promocode.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (s *PromocodeServer) UpdatePromocode(ctx context.Context, req *UpdatePromocodeRequest) (*PromocodeResponse, error) {
	promocode := &model.Promocode{
		ID:          req.Id,
		Name:        req.Name,
		Description: req.Description,
		CreatorID:   req.CreatorId,
		Discount:    req.Discount,
		Code:        req.Code,
	}

	err := s.service.UpdatePromocode(ctx, promocode)
	if err != nil {
		if err.Error() == "unauthorized: you can only update your own promocodes" {
			return nil, status.Error(codes.PermissionDenied, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &PromocodeResponse{
		Id:          promocode.ID,
		Name:        promocode.Name,
		Description: promocode.Description,
		CreatorId:   promocode.CreatorID,
		Discount:    promocode.Discount,
		Code:        promocode.Code,
		CreatedAt:   promocode.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   promocode.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (s *PromocodeServer) DeletePromocode(ctx context.Context, req *DeletePromocodeRequest) (*DeletePromocodeResponse, error) {
	promocode, err := s.service.DeletePromocode(ctx, req.Id, req.CreatorId)
	if err != nil {
		if err.Error() == "unauthorized: you can only delete your own promocodes" {
			return nil, status.Error(codes.PermissionDenied, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &DeletePromocodeResponse{
		Success: true,
		Message: "Promocode successfully deleted",
		Id: promocode.ID,
		Name: promocode.Name,
		CreatorId: promocode.CreatorID,
		Discount: promocode.Discount,
		Code: promocode.Code,
	}, nil
}

func (s *PromocodeServer) ListPromocodes(ctx context.Context, req *ListPromocodesRequest) (*ListPromocodesResponse, error) {
	promocodes, total, err := s.service.ListPromocodes(ctx, int(req.Page), int(req.PerPage))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	response := &ListPromocodesResponse{
		Promocodes: make([]*PromocodeResponse, 0, len(promocodes)),
		Total:      int32(total),
		Page:       req.Page,
		PerPage:    req.PerPage,
	}

	for _, promocode := range promocodes {
		response.Promocodes = append(response.Promocodes, &PromocodeResponse{
			Id:          promocode.ID,
			Name:        promocode.Name,
			Description: promocode.Description,
			CreatorId:   promocode.CreatorID,
			Discount:    promocode.Discount,
			Code:        promocode.Code,
			CreatedAt:   promocode.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   promocode.UpdatedAt.Format(time.RFC3339),
		})
	}

	return response, nil
}
