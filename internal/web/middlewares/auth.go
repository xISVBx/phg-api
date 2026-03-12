package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	appsvc "photogallery/api_go/internal/application/interfaces/services"
	webutils "photogallery/api_go/internal/web/utils"
)

func Auth(jwt appsvc.IJWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if h == "" || !strings.HasPrefix(h, "Bearer ") {
			webutils.Fail(c, http.StatusUnauthorized, "UNAUTHORIZED", "missing bearer token")
			c.Abort()
			return
		}
		tok := strings.TrimSpace(strings.TrimPrefix(h, "Bearer "))
		claims, err := jwt.Parse(tok)
		if err != nil {
			webutils.Fail(c, http.StatusUnauthorized, "UNAUTHORIZED", "invalid token")
			c.Abort()
			return
		}
		webutils.SetUserID(c, claims.UserID)
		c.Next()
	}
}
