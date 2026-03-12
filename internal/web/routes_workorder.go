package web

import (
	"github.com/gin-gonic/gin"

	httpctl "photogallery/api_go/internal/web/controllers/http"
)

func registerWorkOrderRoutes(secure *gin.RouterGroup, c *httpctl.Controller) {
	secure.GET("/work-orders", c.ListWorkOrders)
	secure.POST("/work-orders", c.CreateWorkOrder)
	secure.GET("/work-orders/:id", c.GetWorkOrder)
	secure.PUT("/work-orders/:id", c.UpdateWorkOrder)
}
