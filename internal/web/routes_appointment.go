package web

import (
	"github.com/gin-gonic/gin"

	httpctl "photogallery/api_go/internal/web/controllers/http"
)

func registerAppointmentRoutes(secure *gin.RouterGroup, c *httpctl.Controller) {
	secure.GET("/appointments", c.ListAppointments)
	secure.POST("/appointments", c.CreateAppointment)
	secure.GET("/appointments/:id", c.GetAppointment)
	secure.PUT("/appointments/:id", c.UpdateAppointment)
}
