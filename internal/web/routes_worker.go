package web

import (
	"github.com/gin-gonic/gin"

	httpctl "photogallery/api_go/internal/web/controllers/http"
)

func registerWorkerRoutes(secure *gin.RouterGroup, c *httpctl.Controller) {
	secure.GET("/workers", c.ListWorkers)
	secure.POST("/workers", c.CreateWorker)
	secure.GET("/workers/:id", c.GetWorker)
	secure.PUT("/workers/:id", c.UpdateWorker)
	secure.PATCH("/workers/:id/activate", c.ActivateWorker)
	secure.PATCH("/workers/:id/deactivate", c.DeactivateWorker)
	secure.POST("/workers/pay-commissions", c.PayCommissions)
	secure.POST("/workers/pay-salary", c.PaySalary)
}
