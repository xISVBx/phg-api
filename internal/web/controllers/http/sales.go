package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
	salesreq "photogallery/api_go/internal/application/dtos/request/sales"
	webutils "photogallery/api_go/internal/web/utils"
)

// @Summary ListSales
// @Description Lista paginada de ventas con búsqueda y orden opcional.
// @Tags sales
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page" default(1)
// @Param pageSize query int false "Page size" default(25)
// @Param q query string false "Search text"
// @Param sort query string false "Sort field"
// @Param dir query string false "Sort direction" Enums(asc,desc)
// @Success 200 {object} SwaggerSaleListResponse
// @Failure 500 {object} SwaggerErrorResponse
// @Router /api/v1/sales [get]
func (h *Controller) ListSales(c *gin.Context) {
	items, total, err := h.uc.Sales.List(c.Request.Context(), parseQuery(c))
	if err != nil {
		webutils.Fail(c, http.StatusInternalServerError, "INTERNAL", err.Error())
		return
	}
	webutils.OKMeta(c, items, gin.H{"total": total})
}

// @Summary CreateSale
// @Description Crea una nueva venta.
// @Tags sales
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body SwaggerCreateSaleRequest true "Create sale payload"
// @Success 201 {object} SwaggerSaleResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Router /api/v1/sales [post]
func (h *Controller) CreateSale(c *gin.Context) {
	var in salesreq.CreateSaleRequestDTO
	if err := c.ShouldBindJSON(&in); err != nil {
		webutils.Fail(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	out, err := h.uc.Sales.Create(c.Request.Context(), actorID(c), in)
	if err != nil {
		webutils.Fail(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	webutils.Created(c, out)
}

// @Summary GetSale
// @Description Obtiene una venta por identificador.
// @Tags sales
// @Produce json
// @Security BearerAuth
// @Param id path string true "Resource ID (UUID)"
// @Success 200 {object} SwaggerSaleDetailResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Failure 404 {object} SwaggerErrorResponse
// @Router /api/v1/sales/:id [get]
func (h *Controller) GetSale(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	sale, items, payments, err := h.uc.Sales.Get(c.Request.Context(), id)
	if err != nil {
		webutils.Fail(c, http.StatusNotFound, "NOT_FOUND", err.Error())
		return
	}
	webutils.OK(c, gin.H{"sale": sale, "items": items, "payments": payments})
}

// @Summary RegisterSalePayment
// @Description Registra un pago para una venta.
// @Tags sales
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body SwaggerRegisterSalePaymentRequest true "Register payment payload"
// @Param id path string true "Resource ID (UUID)"
// @Success 200 {object} SwaggerSalePaymentResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Router /api/v1/sales/:id/payments [post]
func (h *Controller) RegisterSalePayment(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	var in salesreq.RegisterSalePaymentRequestDTO
	if err := c.ShouldBindJSON(&in); err != nil {
		webutils.Fail(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	out, err := h.uc.Sales.RegisterPayment(c.Request.Context(), actorID(c), id, in)
	if err != nil {
		webutils.Fail(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	webutils.OK(c, out)
}
