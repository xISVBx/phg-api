package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
	authreq "photogallery/api_go/internal/application/dtos/request/auth"
	webutils "photogallery/api_go/internal/web/utils"
	"time"
)

// @Summary Health
// @Description Verifica el estado de la API y el tiempo del servidor.
// @Tags system
// @Produce json
// @Success 200 {object} SwaggerHealthResponse
// @Router /api/v1/health [get]
func (h *Controller) Health(c *gin.Context) {
	webutils.OK(c, gin.H{"status": "ok", "timeUtc": time.Now().UTC(), "version": "v1"})
}

// @Summary SystemInfo
// @Description Retorna información básica del servicio en ejecución.
// @Tags system
// @Produce json
// @Success 200 {object} SwaggerSystemInfoResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Router /api/v1/system/info [get]
func (h *Controller) SystemInfo(c *gin.Context) {
	webutils.OK(c, gin.H{"name": "photo-gallery-api", "timeUtc": time.Now().UTC()})
}

// @Summary BackupRun
// @Description Dispara una tarea de respaldo (placeholder MVP).
// @Tags system
// @Produce json
// @Success 200 {object} SwaggerBackupResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Router /api/v1/system/backup/run [post]
func (h *Controller) BackupRun(c *gin.Context) {
	webutils.OK(c, gin.H{"status": "not_implemented_mvp"})
}

// @Summary BackupStatus
// @Description Retorna el estado del respaldo (placeholder MVP).
// @Tags system
// @Produce json
// @Success 200 {object} SwaggerBackupResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Router /api/v1/system/backup/status [get]
func (h *Controller) BackupStatus(c *gin.Context) {
	webutils.OK(c, gin.H{"status": "not_implemented_mvp"})
}

// @Summary Login
// @Description Autentica credenciales y devuelve tokens de acceso y refresh.
// @Tags auth
// @Accept json
// @Produce json
// @Param body body SwaggerLoginRequest true "Login payload"
// @Success 200 {object} SwaggerAuthLoginResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Failure 401 {object} SwaggerErrorResponse
// @Router /api/v1/auth/login [post]
func (h *Controller) Login(c *gin.Context) {
	var in authreq.LoginRequestDTO
	if err := c.ShouldBindJSON(&in); err != nil {
		webutils.Fail(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	out, err := h.uc.Auth.Login(c.Request.Context(), in)
	if err != nil {
		webutils.Fail(c, http.StatusUnauthorized, "UNAUTHORIZED", err.Error())
		return
	}
	webutils.OK(c, out)
}

// @Summary Refresh
// @Description Genera nuevos tokens a partir de un refresh token válido.
// @Tags auth
// @Accept json
// @Produce json
// @Param body body SwaggerRefreshRequest true "Refresh payload"
// @Success 200 {object} SwaggerAuthLoginResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Failure 401 {object} SwaggerErrorResponse
// @Router /api/v1/auth/refresh [post]
func (h *Controller) Refresh(c *gin.Context) {
	var in authreq.RefreshRequestDTO
	if err := c.ShouldBindJSON(&in); err != nil {
		webutils.Fail(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	out, err := h.uc.Auth.Refresh(c.Request.Context(), in)
	if err != nil {
		webutils.Fail(c, http.StatusUnauthorized, "UNAUTHORIZED", err.Error())
		return
	}
	webutils.OK(c, out)
}

// @Summary Me
// @Description Retorna el perfil del usuario autenticado.
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} SwaggerAuthMeResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Failure 404 {object} SwaggerErrorResponse
// @Router /api/v1/auth/me [get]
func (h *Controller) Me(c *gin.Context) {
	out, err := h.uc.Auth.Me(c.Request.Context(), actorID(c))
	if err != nil {
		webutils.Fail(c, http.StatusNotFound, "NOT_FOUND", err.Error())
		return
	}
	webutils.OK(c, out)
}

// @Summary Logout
// @Description Cierra la sesión actual en cliente (invalidación de token opcional).
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} SwaggerOKFlagResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Router /api/v1/auth/logout [post]
func (h *Controller) Logout(c *gin.Context) { webutils.OK(c, gin.H{"ok": true}) }

// @Summary ChangePassword
// @Description Cambia la contraseña del usuario autenticado.
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body SwaggerChangePasswordRequest true "Change password payload"
// @Success 200 {object} SwaggerOKFlagResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Router /api/v1/auth/change-password [post]
func (h *Controller) ChangePassword(c *gin.Context) {
	var in authreq.ChangePasswordRequestDTO
	if err := c.ShouldBindJSON(&in); err != nil {
		webutils.Fail(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	if err := h.uc.Auth.ChangePassword(c.Request.Context(), actorID(c), in); err != nil {
		webutils.Fail(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	webutils.OK(c, gin.H{"ok": true})
}

// @Summary MyPermissions
// @Description Retorna los permisos efectivos calculados por rol y overrides.
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} SwaggerPermissionTreeResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Router /api/v1/auth/my-permissions [get]
func (h *Controller) MyPermissions(c *gin.Context) {
	out, err := h.uc.Auth.MyPermissions(c.Request.Context(), actorID(c))
	if err != nil {
		webutils.Fail(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	webutils.OK(c, out)
}
