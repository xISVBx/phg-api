package web

import (
	"github.com/gin-gonic/gin"

	appif "photogallery/api_go/internal/application/interfaces/services"
	"photogallery/api_go/internal/application/use_cases"
	httpctl "photogallery/api_go/internal/web/controllers/http"
	"photogallery/api_go/internal/web/middlewares"
)

func NewServer(uc *use_cases.UseCases, jwt appif.IJWTService, corsAllowedOrigins []string) *gin.Engine {
	r := gin.Default()
	r.Use(middlewares.CORS(corsAllowedOrigins))
	ctl := httpctl.NewController(uc)
	RegisterRoutes(r, ctl, middlewares.Auth(jwt))
	return r
}
