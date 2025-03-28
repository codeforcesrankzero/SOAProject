package api

import (
	"promocodes-service/internal/auth"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, handler *Handler) {
	promocodeGroup := router.Group("/api/promocodes")
	promocodeGroup.Use(auth.AuthMiddleware(handler.secret))
	{
		promocodeGroup.POST("", handler.CreatePromocode)
		promocodeGroup.GET("/:id", handler.GetPromocode)
		promocodeGroup.PUT("/:id", handler.UpdatePromocode)
		promocodeGroup.DELETE("/:id", handler.DeletePromocode)
		promocodeGroup.GET("", handler.ListPromocodes)
	}
}
