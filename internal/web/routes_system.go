package web

import (
	"github.com/gin-gonic/gin"

	httpctl "photogallery/api_go/internal/web/controllers/http"
)

func registerSystemPublicRoutes(v1 *gin.RouterGroup, c *httpctl.Controller) {
	v1.GET("/health", c.Health)
	v1.GET("/system/info", c.SystemInfo)
	v1.POST("/system/backup/run", c.BackupRun)
	v1.GET("/system/backup/status", c.BackupStatus)
}
