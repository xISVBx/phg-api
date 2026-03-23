package web

import (
	"github.com/gin-gonic/gin"

	appif "photogallery/api_go/internal/application/interfaces/services"
	"photogallery/api_go/internal/application/use_cases"
	"photogallery/api_go/internal/infrastructure/bootstrap"
	httpctl "photogallery/api_go/internal/web/controllers/http"
	"photogallery/api_go/internal/web/middlewares"
)

func NewServer(uc *use_cases.UseCases, jwt appif.IJWTService, corsAllowedOrigins []string, bootstrapSvc *bootstrap.Service) *gin.Engine {
	r := gin.Default()
	r.Use(middlewares.CORS(corsAllowedOrigins))
	ctl := httpctl.NewController(uc)
	bootstrapCtl := httpctl.NewBootstrapController(bootstrapSvc)
	RegisterRoutes(r, ctl, bootstrapCtl, middlewares.Auth(jwt))
	return r
}
