package web

import (
	"github.com/gin-gonic/gin"

	httpctl "photogallery/api_go/internal/web/controllers/http"
)

func registerCashRoutes(secure *gin.RouterGroup, c *httpctl.Controller) {
	secure.GET("/cash/categories", c.ListCashCategories)
	secure.POST("/cash/categories", c.CreateCashCategory)
	secure.GET("/cash/categories/:id", c.GetCashCategory)
	secure.PUT("/cash/categories/:id", c.UpdateCashCategory)
	secure.GET("/cash/movements", c.ListCashMovements)
	secure.POST("/cash/movements", c.CreateCashMovement)
	secure.GET("/cash/movements/:id", c.GetCashMovement)
	secure.PUT("/cash/movements/:id", c.UpdateCashMovement)
}
