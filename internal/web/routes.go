package web

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	httpctl "photogallery/api_go/internal/web/controllers/http"
)

func RegisterRoutes(r *gin.Engine, c *httpctl.Controller, auth gin.HandlerFunc) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/scalar", scalarReferenceUI)
	r.GET("/scalar/", scalarReferenceUI)

	v1 := r.Group("/api/v1")
	registerSystemPublicRoutes(v1, c)
	registerAuthPublicRoutes(v1, c)

	secure := v1.Group("")
	secure.Use(auth)
	registerAuthSecureRoutes(secure, c)
	registerSecurityRoutes(secure, c)
	registerAuditRoutes(secure, c)
	registerCatalogRoutes(secure, c)
	registerCustomerRoutes(secure, c)
	registerSalesRoutes(secure, c)
	registerWorkOrderRoutes(secure, c)
	registerAppointmentRoutes(secure, c)
	registerFileRoutes(secure, c)
	registerCashRoutes(secure, c)
	registerWorkerRoutes(secure, c)
	registerSettingRoutes(secure, c)
}

func scalarReferenceUI(c *gin.Context) {
	const html = `<!doctype html>
<html>
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title>PhotoGallery API - Scalar</title>
  </head>
  <body>
    <script id="api-reference" data-url="/swagger/doc.json"></script>
    <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
  </body>
</html>`
	c.Data(200, "text/html; charset=utf-8", []byte(html))
}
