package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
	systemreq "photogallery/api_go/internal/application/dtos/request/system"
	webutils "photogallery/api_go/internal/web/utils"
)

// @Summary GetSystemSetting
// @Description Obtiene una configuración del sistema por clave.
// @Tags settings
// @Produce json
// @Security BearerAuth
// @Param key path string true "Setting key"
// @Success 200 {object} SwaggerSettingResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Failure 404 {object} SwaggerErrorResponse
// @Router /api/v1/settings/:key [get]
func (h *Controller) GetSystemSetting(c *gin.Context) {
	key := c.Param("key")
	out, err := h.uc.System.GetSetting(c.Request.Context(), key)
	if err != nil {
		webutils.Fail(c, http.StatusNotFound, "NOT_FOUND", err.Error())
		return
	}
	webutils.OK(c, out)
}

// @Summary SetSystemSetting
// @Description Actualiza el valor de una configuración del sistema.
// @Tags settings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body SwaggerSetSettingRequest true "Set setting payload"
// @Param key path string true "Setting key"
// @Success 200 {object} SwaggerOKFlagResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Router /api/v1/settings/:key [put]
func (h *Controller) SetSystemSetting(c *gin.Context) {
	key := c.Param("key")
	var in systemreq.SetAppSettingRequestDTO
	if err := c.ShouldBindJSON(&in); err != nil {
		webutils.Fail(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	if err := h.uc.System.SetSetting(c.Request.Context(), actorID(c), key, in.Value); err != nil {
		webutils.Fail(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	webutils.OK(c, gin.H{"ok": true})
}
