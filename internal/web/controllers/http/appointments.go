package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
	appointmentreq "photogallery/api_go/internal/application/dtos/request/appointment"
	webutils "photogallery/api_go/internal/web/utils"
)

// @Summary ListAppointments
// @Description Lista paginada de citas con búsqueda y orden opcional.
// @Tags appointments
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page" default(1)
// @Param pageSize query int false "Page size" default(25)
// @Param q query string false "Search text"
// @Param sort query string false "Sort field"
// @Param dir query string false "Sort direction" Enums(asc,desc)
// @Success 200 {object} SwaggerAppointmentListResponse
// @Failure 500 {object} SwaggerErrorResponse
// @Router /api/v1/appointments [get]
func (h *Controller) ListAppointments(c *gin.Context) {
	items, total, err := h.uc.Appointments.List(c.Request.Context(), parseQuery(c))
	if err != nil {
		webutils.Fail(c, http.StatusInternalServerError, "INTERNAL", err.Error())
		return
	}
	webutils.OKMeta(c, items, gin.H{"total": total})
}

// @Summary CreateAppointment
// @Description Crea una nueva cita.
// @Tags appointments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body SwaggerCreateAppointmentRequest true "Create appointment payload"
// @Success 201 {object} SwaggerAppointmentResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Router /api/v1/appointments [post]
func (h *Controller) CreateAppointment(c *gin.Context) {
	var in appointmentreq.CreateAppointmentRequestDTO
	if err := c.ShouldBindJSON(&in); err != nil {
		webutils.Fail(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	out, err := h.uc.Appointments.Create(c.Request.Context(), actorID(c), in)
	if err != nil {
		webutils.Fail(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	webutils.Created(c, out)
}

// @Summary GetAppointment
// @Description Obtiene una cita por identificador.
// @Tags appointments
// @Produce json
// @Security BearerAuth
// @Param id path string true "Resource ID (UUID)"
// @Success 200 {object} SwaggerAppointmentResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Failure 404 {object} SwaggerErrorResponse
// @Router /api/v1/appointments/:id [get]
func (h *Controller) GetAppointment(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	out, err := h.uc.Appointments.Get(c.Request.Context(), id)
	if err != nil {
		webutils.Fail(c, http.StatusNotFound, "NOT_FOUND", err.Error())
		return
	}
	webutils.OK(c, out)
}

// @Summary UpdateAppointment
// @Description Actualiza una cita existente.
// @Tags appointments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body SwaggerUpdateAppointmentRequest true "Update appointment payload"
// @Param id path string true "Resource ID (UUID)"
// @Success 200 {object} SwaggerAppointmentResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Router /api/v1/appointments/:id [put]
func (h *Controller) UpdateAppointment(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	var in appointmentreq.UpdateAppointmentRequestDTO
	if err := c.ShouldBindJSON(&in); err != nil {
		webutils.Fail(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	out, err := h.uc.Appointments.Update(c.Request.Context(), actorID(c), id, in)
	if err != nil {
		webutils.Fail(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	webutils.OK(c, out)
}
