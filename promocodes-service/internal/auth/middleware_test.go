package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware_ValidToken(t *testing.T) {
	router := gin.Default()
	router.Use(AuthMiddleware("secret"))

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = 1

	tokenString, _ := token.SignedString([]byte("secret"))

	router.GET("/test", func(c *gin.Context) {
		userID, err := GetUserID(c)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), userID)
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthMiddleware_MissingAuthorizationHeader(t *testing.T) {
	router := gin.Default()
	router.Use(AuthMiddleware("secret"))

	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Authorization header is required")
}

func TestAuthMiddleware_InvalidAuthorizationHeader(t *testing.T) {
	router := gin.Default()
	router.Use(AuthMiddleware("secret"))

	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "InvalidTokenFormat")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Authorization header format must be Bearer {token}")
}

func TestGetUserID_NoUserIDInContext(t *testing.T) {
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	id, err := GetUserID(ctx)
	assert.Error(t, err)
	assert.Equal(t, int64(0), id)
	assert.Equal(t, "user ID not found in context", err.Error())
}

func TestGetUserID_InvalidUserIDType(t *testing.T) {
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Set("userID", "notAnInt64")
	id, err := GetUserID(ctx)
	assert.Error(t, err)
	assert.Equal(t, int64(0), id)
	assert.Contains(t, err.Error(), "user ID has invalid type")
}
