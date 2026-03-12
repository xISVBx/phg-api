package web

import (
	"github.com/gin-gonic/gin"

	httpctl "photogallery/api_go/internal/web/controllers/http"
)

func registerCustomerRoutes(secure *gin.RouterGroup, c *httpctl.Controller) {
	secure.GET("/customers", c.ListCustomers)
	secure.POST("/customers", c.CreateCustomer)
	secure.GET("/customers/:id", c.GetCustomer)
	secure.PUT("/customers/:id", c.UpdateCustomer)
	secure.PATCH("/customers/:id/activate", c.ActivateCustomer)
	secure.PATCH("/customers/:id/deactivate", c.DeactivateCustomer)
}
