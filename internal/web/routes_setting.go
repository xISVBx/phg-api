package web

import (
	"github.com/gin-gonic/gin"

	httpctl "photogallery/api_go/internal/web/controllers/http"
)

func registerSettingRoutes(secure *gin.RouterGroup, c *httpctl.Controller) {
	secure.GET("/settings/:key", c.GetSystemSetting)
	secure.PUT("/settings/:key", c.SetSystemSetting)
}
