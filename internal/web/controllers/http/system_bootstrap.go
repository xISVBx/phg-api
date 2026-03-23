package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"photogallery/api_go/internal/infrastructure/bootstrap"
	webutils "photogallery/api_go/internal/web/utils"
)

const bootstrapTokenHeader = "X-Bootstrap-Token"

type BootstrapController struct {
	bootstrapSvc *bootstrap.Service
}

func NewBootstrapController(bootstrapSvc *bootstrap.Service) *BootstrapController {
	return &BootstrapController{bootstrapSvc: bootstrapSvc}
}

// @Summary RunBootstrap
// @Description Ejecuta el bootstrap inicial de seguridad de forma manual y excepcional.
// @Tags system
// @Accept json
// @Produce json
// @Param X-Bootstrap-Token header string true "Bootstrap token"
// @Success 200 {object} map[string]any
// @Failure 401 {object} SwaggerErrorResponse
// @Failure 403 {object} SwaggerErrorResponse
// @Failure 500 {object} SwaggerErrorResponse
// @Router /api/v1/system/bootstrap/run [post]
func (h *BootstrapController) RunBootstrap(c *gin.Context) {
	if h == nil || h.bootstrapSvc == nil {
		webutils.Fail(c, http.StatusInternalServerError, "INTERNAL", "bootstrap service not configured")
		return
	}

	token := c.GetHeader(bootstrapTokenHeader)
	if token == "" {
		webutils.Fail(c, http.StatusUnauthorized, "UNAUTHORIZED", "missing bootstrap token")
		return
	}
	if !h.bootstrapSvc.ValidateManualToken(token) {
		webutils.Fail(c, http.StatusForbidden, "FORBIDDEN", "invalid bootstrap token")
		return
	}

	result, err := h.bootstrapSvc.RunManual(c.Request.Context())
	if err != nil {
		webutils.Fail(c, http.StatusInternalServerError, "INTERNAL", err.Error())
		return
	}

	webutils.OK(c, result)
}
