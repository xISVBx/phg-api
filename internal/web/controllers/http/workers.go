package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
	workerreq "photogallery/api_go/internal/application/dtos/request/worker"
	webutils "photogallery/api_go/internal/web/utils"
)

// @Summary ListWorkers
// @Description Lista paginada de trabajadores con búsqueda y orden opcional.
// @Tags workers
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page" default(1)
// @Param pageSize query int false "Page size" default(25)
// @Param q query string false "Search text"
// @Param sort query string false "Sort field"
// @Param dir query string false "Sort direction" Enums(asc,desc)
// @Success 200 {object} SwaggerWorkerListResponse
// @Failure 500 {object} SwaggerErrorResponse
// @Router /api/v1/workers [get]
func (h *Controller) ListWorkers(c *gin.Context) {
	items, total, err := h.uc.Workers.List(c.Request.Context(), parseQuery(c))
	if err != nil {
		webutils.Fail(c, http.StatusInternalServerError, "INTERNAL", err.Error())
		return
	}
	webutils.OKMeta(c, items, gin.H{"total": total})
}

// @Summary CreateWorker
// @Description Crea un nuevo trabajador.
// @Tags workers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body SwaggerCreateWorkerRequest true "Create worker payload"
// @Success 201 {object} SwaggerWorkerResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Router /api/v1/workers [post]
func (h *Controller) CreateWorker(c *gin.Context) {
	var in workerreq.CreateWorkerRequestDTO
	if err := c.ShouldBindJSON(&in); err != nil {
		webutils.Fail(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	out, err := h.uc.Workers.Create(c.Request.Context(), actorID(c), in)
	if err != nil {
		webutils.Fail(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	webutils.Created(c, out)
}

// @Summary GetWorker
// @Description Obtiene un trabajador por identificador.
// @Tags workers
// @Produce json
// @Security BearerAuth
// @Param id path string true "Resource ID (UUID)"
// @Success 200 {object} SwaggerWorkerResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Failure 404 {object} SwaggerErrorResponse
// @Router /api/v1/workers/:id [get]
func (h *Controller) GetWorker(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	out, err := h.uc.Workers.Get(c.Request.Context(), id)
	if err != nil {
		webutils.Fail(c, http.StatusNotFound, "NOT_FOUND", err.Error())
		return
	}
	webutils.OK(c, out)
}

// @Summary UpdateWorker
// @Description Actualiza un trabajador existente.
// @Tags workers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body SwaggerUpdateWorkerRequest true "Update worker payload"
// @Param id path string true "Resource ID (UUID)"
// @Success 200 {object} SwaggerWorkerResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Router /api/v1/workers/:id [put]
func (h *Controller) UpdateWorker(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	var in workerreq.UpdateWorkerRequestDTO
	if err := c.ShouldBindJSON(&in); err != nil {
		webutils.Fail(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	out, err := h.uc.Workers.Update(c.Request.Context(), actorID(c), id, in)
	if err != nil {
		webutils.Fail(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	webutils.OK(c, out)
}

// @Summary ActivateWorker
// @Description Activa un trabajador.
// @Tags workers
// @Produce json
// @Security BearerAuth
// @Param id path string true "Resource ID (UUID)"
// @Success 200 {object} SwaggerOKFlagResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Router /api/v1/workers/:id/activate [patch]
func (h *Controller) ActivateWorker(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	if err := h.uc.Workers.SetActive(c.Request.Context(), actorID(c), id, true); err != nil {
		webutils.Fail(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	webutils.OK(c, gin.H{"ok": true})
}

// @Summary DeactivateWorker
// @Description Desactiva un trabajador.
// @Tags workers
// @Produce json
// @Security BearerAuth
// @Param id path string true "Resource ID (UUID)"
// @Success 200 {object} SwaggerOKFlagResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Router /api/v1/workers/:id/deactivate [patch]
func (h *Controller) DeactivateWorker(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	if err := h.uc.Workers.SetActive(c.Request.Context(), actorID(c), id, false); err != nil {
		webutils.Fail(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	webutils.OK(c, gin.H{"ok": true})
}

// @Summary PayCommissions
// @Description Ejecuta el pago de comisiones pendientes.
// @Tags workers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body SwaggerPayCommissionRequest true "Pay commissions payload"
// @Success 200 {object} SwaggerOKFlagResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Router /api/v1/workers/pay-commissions [post]
func (h *Controller) PayCommissions(c *gin.Context) {
	var in workerreq.PayCommissionRequestDTO
	if err := c.ShouldBindJSON(&in); err != nil {
		webutils.Fail(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	if err := h.uc.Workers.PayCommissionFIFO(c.Request.Context(), actorID(c), in); err != nil {
		webutils.Fail(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	webutils.OK(c, gin.H{"ok": true})
}

// @Summary PaySalary
// @Description Ejecuta el pago de salario.
// @Tags workers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body SwaggerPaySalaryRequest true "Pay salary payload"
// @Success 200 {object} SwaggerOKFlagResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Router /api/v1/workers/pay-salary [post]
func (h *Controller) PaySalary(c *gin.Context) {
	var in workerreq.PaySalaryRequestDTO
	if err := c.ShouldBindJSON(&in); err != nil {
		webutils.Fail(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	if err := h.uc.Workers.PaySalary(c.Request.Context(), actorID(c), in); err != nil {
		webutils.Fail(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	webutils.OK(c, gin.H{"ok": true})
}
