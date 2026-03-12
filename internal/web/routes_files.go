package web

import (
	"github.com/gin-gonic/gin"

	httpctl "photogallery/api_go/internal/web/controllers/http"
)

func registerFileRoutes(secure *gin.RouterGroup, c *httpctl.Controller) {
	secure.GET("/files", c.ListFiles)
	secure.POST("/files/upload", c.UploadFile)
	secure.GET("/files/:id/download", c.DownloadFile)
	secure.DELETE("/files/:id", c.DeleteFile)
}
