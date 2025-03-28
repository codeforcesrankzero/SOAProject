package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"promocodes-service/api"
	"promocodes-service/internal/auth"
	"promocodes-service/pkg/grpc"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/golang-jwt/jwt"
)

const testSecretKey = "some-secret"

type MockGrpcClient struct {
	mock.Mock
}

func (m *MockGrpcClient) CreatePromocode(ctx context.Context, req *grpc.CreatePromocodeRequest) (*grpc.PromocodeResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*grpc.PromocodeResponse), args.Error(1)
}

func (m *MockGrpcClient) GetPromocode(ctx context.Context, req *grpc.GetPromocodeRequest) (*grpc.PromocodeResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*grpc.PromocodeResponse), args.Error(1)
}

func (m *MockGrpcClient) UpdatePromocode(ctx context.Context, req *grpc.UpdatePromocodeRequest) (*grpc.PromocodeResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*grpc.PromocodeResponse), args.Error(1)
}

func (m *MockGrpcClient) DeletePromocode(ctx context.Context, req *grpc.DeletePromocodeRequest) (*grpc.DeletePromocodeResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*grpc.DeletePromocodeResponse), args.Error(1)
}

func (m *MockGrpcClient) ListPromocodes(ctx context.Context, req *grpc.ListPromocodesRequest) (*grpc.ListPromocodesResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*grpc.ListPromocodesResponse), args.Error(1)
}

func generateJWT(userID int64, secretKey string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

func TestCreatePromocode(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockClient := new(MockGrpcClient)
	handler := api.NewHandler(mockClient, testSecretKey)

	router := gin.New()
	router.Use(auth.AuthMiddleware(testSecretKey))
	router.POST("/promocodes", handler.CreatePromocode)

	expectedResponse := &grpc.PromocodeResponse{Id: 1}
	mockClient.On("CreatePromocode", mock.Anything, mock.AnythingOfType("*grpc.CreatePromocodeRequest")).Return(expectedResponse, nil)

	w := httptest.NewRecorder()
	body := `{"name": "Promo", "description": "Description", "discount": 10, "code": "CODE123"}`
	req, _ := http.NewRequest(http.MethodPost, "/promocodes", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	testToken, _ := generateJWT(1, testSecretKey)
	req.Header.Set("Authorization", "Bearer "+testToken)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var response struct{ Id int }
	json.NewDecoder(w.Body).Decode(&response)
	assert.Equal(t, 1, response.Id)
	mockClient.AssertExpectations(t)
}

func TestGetPromocode(t *testing.T) {
	mockClient := new(MockGrpcClient)
	handler := api.NewHandler(mockClient, testSecretKey)

	router := gin.New()
	router.Use(auth.AuthMiddleware(testSecretKey))
	router.GET("/promocodes/:id", handler.GetPromocode)

	expectedResponse := &grpc.PromocodeResponse{Id: 1}
	mockClient.On("GetPromocode", mock.Anything, mock.AnythingOfType("*grpc.GetPromocodeRequest")).Return(expectedResponse, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/promocodes/1", nil)

	testToken, _ := generateJWT(1, testSecretKey)
	req.Header.Set("Authorization", "Bearer "+testToken)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response struct{ Id int }
	json.NewDecoder(w.Body).Decode(&response)
	assert.Equal(t, 1, response.Id)
	mockClient.AssertExpectations(t)
}

func TestUpdatePromocode(t *testing.T) {
	mockClient := new(MockGrpcClient)
	handler := api.NewHandler(mockClient, testSecretKey)

	router := gin.New()
	router.Use(auth.AuthMiddleware(testSecretKey))
	router.PUT("/promocodes/:id", handler.UpdatePromocode)

	expectedResponse := &grpc.PromocodeResponse{Id: 1}
	mockClient.On("UpdatePromocode", mock.Anything, mock.AnythingOfType("*grpc.UpdatePromocodeRequest")).Return(expectedResponse, nil)

	w := httptest.NewRecorder()
	body := `{"name": "Updated Promo", "description": "Updated Description", "discount": 15, "code": "NEWCODE123"}`
	req, _ := http.NewRequest(http.MethodPut, "/promocodes/1", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	testToken, _ := generateJWT(1, testSecretKey)
	req.Header.Set("Authorization", "Bearer "+testToken)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response struct{ Id int }
	json.NewDecoder(w.Body).Decode(&response)
	assert.Equal(t, 1, response.Id)
	mockClient.AssertExpectations(t)
}

func TestDeletePromocode(t *testing.T) {
	mockClient := new(MockGrpcClient)
	handler := api.NewHandler(mockClient, testSecretKey)

	router := gin.New()
	router.Use(auth.AuthMiddleware(testSecretKey))
	router.DELETE("/promocodes/:id", handler.DeletePromocode)

	expectedResponse := &grpc.DeletePromocodeResponse{Success: true}
	mockClient.On("DeletePromocode", mock.Anything, mock.AnythingOfType("*grpc.DeletePromocodeRequest")).Return(expectedResponse, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/promocodes/1", nil)

	testToken, _ := generateJWT(1, testSecretKey)
	req.Header.Set("Authorization", "Bearer "+testToken)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response struct{ Success bool }
	json.NewDecoder(w.Body).Decode(&response)
	assert.Equal(t, true, response.Success)
	mockClient.AssertExpectations(t)
}

func TestListPromocodes(t *testing.T) {
	mockClient := new(MockGrpcClient)
	handler := api.NewHandler(mockClient, testSecretKey)

	router := gin.New()
	router.Use(auth.AuthMiddleware(testSecretKey))
	router.GET("/promocodes", handler.ListPromocodes)

	expectedResponse := &grpc.ListPromocodesResponse{
		Promocodes: []*grpc.PromocodeResponse{{Id: 1}, {Id: 2}},
		Total:      2,
		Page:       1,
		PerPage:    10,
	}
	mockClient.On("ListPromocodes", mock.Anything, mock.AnythingOfType("*grpc.ListPromocodesRequest")).Return(expectedResponse, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/promocodes", nil)

	testToken, _ := generateJWT(1, testSecretKey)
	req.Header.Set("Authorization", "Bearer "+testToken)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response struct {
		Items    []*grpc.PromocodeResponse `json:"items"`
		Total    int32                     `json:"total"`
		Page     int32                     `json:"page"`
		PerPage  int32                     `json:"per_page"`
	}
	json.NewDecoder(w.Body).Decode(&response)

	assert.Len(t, response.Items, 2)
	assert.Equal(t, int32(2), response.Total)
	assert.Equal(t, int32(1), response.Page)
	assert.Equal(t, int32(10), response.PerPage)
	mockClient.AssertExpectations(t)
}
