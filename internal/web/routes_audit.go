package web

import (
	"github.com/gin-gonic/gin"

	httpctl "photogallery/api_go/internal/web/controllers/http"
)

func registerAuditRoutes(secure *gin.RouterGroup, c *httpctl.Controller) {
	secure.GET("/audit-logs", c.ListAudit)
	secure.GET("/audit-logs/:id", c.GetAudit)
}
