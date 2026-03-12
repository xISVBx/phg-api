package web

import (
	"github.com/gin-gonic/gin"

	httpctl "photogallery/api_go/internal/web/controllers/http"
)

func registerCatalogRoutes(secure *gin.RouterGroup, c *httpctl.Controller) {
	secure.GET("/categories", c.ListCategories)
	secure.POST("/categories", c.CreateCategory)
	secure.GET("/categories/:id", c.GetCategory)
	secure.PUT("/categories/:id", c.UpdateCategory)
	secure.PATCH("/categories/:id/activate", c.ActivateCategory)
	secure.PATCH("/categories/:id/deactivate", c.DeactivateCategory)

	secure.GET("/products", c.ListProducts)
	secure.POST("/products", c.CreateProduct)
	secure.GET("/products/:id", c.GetProduct)
	secure.PUT("/products/:id", c.UpdateProduct)
	secure.PATCH("/products/:id/activate", c.ActivateProduct)
	secure.PATCH("/products/:id/deactivate", c.DeactivateProduct)
}
