package middlewares

import "github.com/gin-gonic/gin"

func RequirePermission(_ string, _ string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: validar permisos efectivos por submenu/permiso usando AuthUseCase.MyPermissions.
		c.Next()
	}
}
