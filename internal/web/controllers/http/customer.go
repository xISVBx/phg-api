package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
	customerreq "photogallery/api_go/internal/application/dtos/request/customer"
	webutils "photogallery/api_go/internal/web/utils"
)

// @Summary ListCustomers
// @Description Lista paginada de clientes con búsqueda y orden opcional.
// @Tags customers
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page" default(1)
// @Param pageSize query int false "Page size" default(25)
// @Param q query string false "Search text"
// @Param sort query string false "Sort field"
// @Param dir query string false "Sort direction" Enums(asc,desc)
// @Success 200 {object} SwaggerCustomerListResponse
// @Failure 500 {object} SwaggerErrorResponse
// @Router /api/v1/customers [get]
func (h *Controller) ListCustomers(c *gin.Context) {
	items, total, err := h.uc.Customers.List(c.Request.Context(), parseQuery(c))
	if err != nil {
		webutils.Fail(c, http.StatusInternalServerError, "INTERNAL", err.Error())
		return
	}
	webutils.OKMeta(c, items, gin.H{"total": total})
}

// @Summary CreateCustomer
// @Description Crea un nuevo cliente.
// @Tags customers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body SwaggerCreateCustomerRequest true "Create customer payload"
// @Success 201 {object} SwaggerCustomerResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Router /api/v1/customers [post]
func (h *Controller) CreateCustomer(c *gin.Context) {
	var in customerreq.CreateCustomerRequestDTO
	if err := c.ShouldBindJSON(&in); err != nil {
		webutils.Fail(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	out, err := h.uc.Customers.Create(c.Request.Context(), actorID(c), in)
	if err != nil {
		webutils.Fail(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	webutils.Created(c, out)
}

// @Summary GetCustomer
// @Description Obtiene un cliente por identificador.
// @Tags customers
// @Produce json
// @Security BearerAuth
// @Param id path string true "Resource ID (UUID)"
// @Success 200 {object} SwaggerCustomerResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Failure 404 {object} SwaggerErrorResponse
// @Router /api/v1/customers/:id [get]
func (h *Controller) GetCustomer(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	out, err := h.uc.Customers.Get(c.Request.Context(), id)
	if err != nil {
		webutils.Fail(c, http.StatusNotFound, "NOT_FOUND", err.Error())
		return
	}
	webutils.OK(c, out)
}

// @Summary UpdateCustomer
// @Description Actualiza un cliente existente.
// @Tags customers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body SwaggerUpdateCustomerRequest true "Update customer payload"
// @Param id path string true "Resource ID (UUID)"
// @Success 200 {object} SwaggerCustomerResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Router /api/v1/customers/:id [put]
func (h *Controller) UpdateCustomer(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	var in customerreq.UpdateCustomerRequestDTO
	if err := c.ShouldBindJSON(&in); err != nil {
		webutils.Fail(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	out, err := h.uc.Customers.Update(c.Request.Context(), actorID(c), id, in)
	if err != nil {
		webutils.Fail(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	webutils.OK(c, out)
}

// @Summary ActivateCustomer
// @Description Activa un cliente.
// @Tags customers
// @Produce json
// @Security BearerAuth
// @Param id path string true "Resource ID (UUID)"
// @Success 200 {object} SwaggerOKFlagResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Router /api/v1/customers/:id/activate [patch]
func (h *Controller) ActivateCustomer(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	if err := h.uc.Customers.SetActive(c.Request.Context(), actorID(c), id, true); err != nil {
		webutils.Fail(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	webutils.OK(c, gin.H{"ok": true})
}

// @Summary DeactivateCustomer
// @Description Desactiva un cliente.
// @Tags customers
// @Produce json
// @Security BearerAuth
// @Param id path string true "Resource ID (UUID)"
// @Success 200 {object} SwaggerOKFlagResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Router /api/v1/customers/:id/deactivate [patch]
func (h *Controller) DeactivateCustomer(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	if err := h.uc.Customers.SetActive(c.Request.Context(), actorID(c), id, false); err != nil {
		webutils.Fail(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	webutils.OK(c, gin.H{"ok": true})
}
