package api

import (
	"context"
	"net/http"
	"promocodes-service/internal/auth"
	"promocodes-service/internal/model"
	"promocodes-service/pkg/grpc"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GrpcClient interface {
	CreatePromocode(ctx context.Context, req *grpc.CreatePromocodeRequest) (*grpc.PromocodeResponse, error)
	GetPromocode(ctx context.Context, req *grpc.GetPromocodeRequest) (*grpc.PromocodeResponse, error)
	UpdatePromocode(ctx context.Context, req *grpc.UpdatePromocodeRequest) (*grpc.PromocodeResponse, error)
	DeletePromocode(ctx context.Context, req *grpc.DeletePromocodeRequest) (*grpc.DeletePromocodeResponse, error)
	ListPromocodes(ctx context.Context, req *grpc.ListPromocodesRequest) (*grpc.ListPromocodesResponse, error)
}

type Handler struct {
	client GrpcClient
	secret string
}

func NewHandler(client GrpcClient, secret string) *Handler {
	return &Handler{client: client, secret: secret}
}

func (h *Handler) CreatePromocode(c *gin.Context) {
	var request model.CreatePromocodeRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := auth.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.client.CreatePromocode(context.Background(), &grpc.CreatePromocodeRequest{
		Name:        request.Name,
		Description: request.Description,
		CreatorId:   userID,
		Discount:    request.Discount,
		Code:        request.Code,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (h *Handler) GetPromocode(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid promocode ID"})
		return
	}

	resp, err := h.client.GetPromocode(context.Background(), &grpc.GetPromocodeRequest{Id: id})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "promocode not found"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) UpdatePromocode(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid promocode ID"})
		return
	}

	var request model.UpdatePromocodeRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := auth.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.client.UpdatePromocode(context.Background(), &grpc.UpdatePromocodeRequest{
		Id:          id,
		Name:        request.Name,
		Description: request.Description,
		CreatorId:   userID,
		Discount:    request.Discount,
		Code:        request.Code,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) DeletePromocode(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid promocode ID"})
		return
	}

	userID, err := auth.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.client.DeletePromocode(context.Background(), &grpc.DeletePromocodeRequest{
		Id:        id,
		CreatorId: userID,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) ListPromocodes(c *gin.Context) {
	var pagination model.PaginationRequest
	pagination.Page = 1
	pagination.PerPage = 10

	if page, err := strconv.Atoi(c.DefaultQuery("page", "1")); err == nil && page > 0 {
		pagination.Page = page
	}

	if perPage, err := strconv.Atoi(c.DefaultQuery("per_page", "10")); err == nil && perPage > 0 && perPage <= 100 {
		pagination.PerPage = perPage
	}

	resp, err := h.client.ListPromocodes(context.Background(), &grpc.ListPromocodesRequest{
		Page:    int32(pagination.Page),
		PerPage: int32(pagination.PerPage),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items":    resp.Promocodes,
		"total":    resp.Total,
		"page":     resp.Page,
		"per_page": resp.PerPage,
	})
}
