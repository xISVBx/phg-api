package web

import (
	"github.com/gin-gonic/gin"

	httpctl "photogallery/api_go/internal/web/controllers/http"
)

func registerAuthPublicRoutes(v1 *gin.RouterGroup, c *httpctl.Controller) {
	v1.POST("/auth/login", c.Login)
	v1.POST("/auth/refresh", c.Refresh)
}

func registerAuthSecureRoutes(secure *gin.RouterGroup, c *httpctl.Controller) {
	secure.GET("/auth/me", c.Me)
	secure.POST("/auth/logout", c.Logout)
	secure.POST("/auth/change-password", c.ChangePassword)
	secure.GET("/auth/my-permissions", c.MyPermissions)
}
