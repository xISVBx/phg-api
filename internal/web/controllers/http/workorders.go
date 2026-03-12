package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
	workorderreq "photogallery/api_go/internal/application/dtos/request/workorder"
	webutils "photogallery/api_go/internal/web/utils"
)

// @Summary ListWorkOrders
// @Description Lista paginada de órdenes de trabajo con búsqueda y orden opcional.
// @Tags work-orders
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page" default(1)
// @Param pageSize query int false "Page size" default(25)
// @Param q query string false "Search text"
// @Param sort query string false "Sort field"
// @Param dir query string false "Sort direction" Enums(asc,desc)
// @Success 200 {object} SwaggerWorkOrderListResponse
// @Failure 500 {object} SwaggerErrorResponse
// @Router /api/v1/work-orders [get]
func (h *Controller) ListWorkOrders(c *gin.Context) {
	items, total, err := h.uc.WorkOrders.List(c.Request.Context(), parseQuery(c))
	if err != nil {
		webutils.Fail(c, http.StatusInternalServerError, "INTERNAL", err.Error())
		return
	}
	webutils.OKMeta(c, items, gin.H{"total": total})
}

// @Summary CreateWorkOrder
// @Description Crea una nueva orden de trabajo.
// @Tags work-orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body SwaggerCreateWorkOrderRequest true "Create work order payload"
// @Success 201 {object} SwaggerWorkOrderResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Router /api/v1/work-orders [post]
func (h *Controller) CreateWorkOrder(c *gin.Context) {
	var in workorderreq.CreateWorkOrderRequestDTO
	if err := c.ShouldBindJSON(&in); err != nil {
		webutils.Fail(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	out, err := h.uc.WorkOrders.Create(c.Request.Context(), actorID(c), in)
	if err != nil {
		webutils.Fail(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	webutils.Created(c, out)
}

// @Summary GetWorkOrder
// @Description Obtiene una orden de trabajo por identificador.
// @Tags work-orders
// @Produce json
// @Security BearerAuth
// @Param id path string true "Resource ID (UUID)"
// @Success 200 {object} SwaggerWorkOrderResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Failure 404 {object} SwaggerErrorResponse
// @Router /api/v1/work-orders/:id [get]
func (h *Controller) GetWorkOrder(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	out, err := h.uc.WorkOrders.Get(c.Request.Context(), id)
	if err != nil {
		webutils.Fail(c, http.StatusNotFound, "NOT_FOUND", err.Error())
		return
	}
	webutils.OK(c, out)
}

// @Summary UpdateWorkOrder
// @Description Actualiza una orden de trabajo existente.
// @Tags work-orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body SwaggerUpdateWorkOrderRequest true "Update work order payload"
// @Param id path string true "Resource ID (UUID)"
// @Success 200 {object} SwaggerWorkOrderResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Router /api/v1/work-orders/:id [put]
func (h *Controller) UpdateWorkOrder(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	var in workorderreq.UpdateWorkOrderRequestDTO
	if err := c.ShouldBindJSON(&in); err != nil {
		webutils.Fail(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	out, err := h.uc.WorkOrders.Update(c.Request.Context(), actorID(c), id, in)
	if err != nil {
		webutils.Fail(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	webutils.OK(c, out)
}
