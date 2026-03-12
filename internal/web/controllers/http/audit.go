package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
	webutils "photogallery/api_go/internal/web/utils"
)

// @Summary ListAudit
// @Description Lista paginada de auditorías con búsqueda y orden opcional.
// @Tags audit
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page" default(1)
// @Param pageSize query int false "Page size" default(25)
// @Param q query string false "Search text"
// @Param sort query string false "Sort field"
// @Param dir query string false "Sort direction" Enums(asc,desc)
// @Success 200 {object} SwaggerAuditListResponse
// @Failure 500 {object} SwaggerErrorResponse
// @Router /api/v1/audit-logs [get]
func (h *Controller) ListAudit(c *gin.Context) {
	items, total, err := h.uc.Audit.List(c.Request.Context(), parseQuery(c))
	if err != nil {
		webutils.Fail(c, http.StatusInternalServerError, "INTERNAL", err.Error())
		return
	}
	webutils.OKMeta(c, items, gin.H{"total": total})
}

// @Summary GetAudit
// @Description Obtiene un registro de auditoría por identificador.
// @Tags audit
// @Produce json
// @Security BearerAuth
// @Param id path string true "Resource ID (UUID)"
// @Success 200 {object} SwaggerAuditResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Failure 404 {object} SwaggerErrorResponse
// @Router /api/v1/audit-logs/:id [get]
func (h *Controller) GetAudit(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	out, err := h.uc.Audit.Get(c.Request.Context(), id)
	if err != nil {
		webutils.Fail(c, http.StatusNotFound, "NOT_FOUND", err.Error())
		return
	}
	webutils.OK(c, out)
}
