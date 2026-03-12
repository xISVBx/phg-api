package web

import (
	"github.com/gin-gonic/gin"

	httpctl "photogallery/api_go/internal/web/controllers/http"
)

func registerSalesRoutes(secure *gin.RouterGroup, c *httpctl.Controller) {
	secure.GET("/sales", c.ListSales)
	secure.POST("/sales", c.CreateSale)
	secure.GET("/sales/:id", c.GetSale)
	secure.POST("/sales/:id/payments", c.RegisterSalePayment)
}
