package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
	cashreq "photogallery/api_go/internal/application/dtos/request/cash"
	webutils "photogallery/api_go/internal/web/utils"
)

// @Summary ListCashCategories
// @Description Lista paginada de categorías de caja con búsqueda y orden opcional.
// @Tags cash
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page" default(1)
// @Param pageSize query int false "Page size" default(25)
// @Param q query string false "Search text"
// @Param sort query string false "Sort field"
// @Param dir query string false "Sort direction" Enums(asc,desc)
// @Success 200 {object} SwaggerCashCategoryListResponse
// @Failure 500 {object} SwaggerErrorResponse
// @Router /api/v1/cash/categories [get]
func (h *Controller) ListCashCategories(c *gin.Context) {
	items, total, err := h.uc.Cash.ListCategories(c.Request.Context(), parseQuery(c))
	if err != nil {
		webutils.Fail(c, http.StatusInternalServerError, "INTERNAL", err.Error())
		return
	}
	webutils.OKMeta(c, items, gin.H{"total": total})
}

// @Summary CreateCashCategory
// @Description Crea una nueva categoría de caja.
// @Tags cash
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body SwaggerCreateCashCategoryRequest true "Create cash category payload"
// @Success 201 {object} SwaggerCashCategoryResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Router /api/v1/cash/categories [post]
func (h *Controller) CreateCashCategory(c *gin.Context) {
	var in cashreq.CreateCashCategoryRequestDTO
	if err := c.ShouldBindJSON(&in); err != nil {
		webutils.Fail(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	out, err := h.uc.Cash.CreateCategory(c.Request.Context(), actorID(c), in)
	if err != nil {
		webutils.Fail(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	webutils.Created(c, out)
}

// @Summary GetCashCategory
// @Description Obtiene una categoría de caja por identificador.
// @Tags cash
// @Produce json
// @Security BearerAuth
// @Param id path string true "Resource ID (UUID)"
// @Success 200 {object} SwaggerCashCategoryResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Failure 404 {object} SwaggerErrorResponse
// @Router /api/v1/cash/categories/:id [get]
func (h *Controller) GetCashCategory(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	out, err := h.uc.Cash.GetCategory(c.Request.Context(), id)
	if err != nil {
		webutils.Fail(c, http.StatusNotFound, "NOT_FOUND", err.Error())
		return
	}
	webutils.OK(c, out)
}

// @Summary UpdateCashCategory
// @Description Actualiza una categoría de caja existente.
// @Tags cash
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body SwaggerUpdateCashCategoryRequest true "Update cash category payload"
// @Param id path string true "Resource ID (UUID)"
// @Success 200 {object} SwaggerCashCategoryResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Router /api/v1/cash/categories/:id [put]
func (h *Controller) UpdateCashCategory(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	var in cashreq.UpdateCashCategoryRequestDTO
	if err := c.ShouldBindJSON(&in); err != nil {
		webutils.Fail(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	out, err := h.uc.Cash.UpdateCategory(c.Request.Context(), actorID(c), id, in)
	if err != nil {
		webutils.Fail(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	webutils.OK(c, out)
}

// @Summary ListCashMovements
// @Description Lista paginada de movimientos de caja con búsqueda y orden opcional.
// @Tags cash
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page" default(1)
// @Param pageSize query int false "Page size" default(25)
// @Param q query string false "Search text"
// @Param sort query string false "Sort field"
// @Param dir query string false "Sort direction" Enums(asc,desc)
// @Success 200 {object} SwaggerCashMovementListResponse
// @Failure 500 {object} SwaggerErrorResponse
// @Router /api/v1/cash/movements [get]
func (h *Controller) ListCashMovements(c *gin.Context) {
	items, total, err := h.uc.Cash.ListMovements(c.Request.Context(), parseQuery(c))
	if err != nil {
		webutils.Fail(c, http.StatusInternalServerError, "INTERNAL", err.Error())
		return
	}
	webutils.OKMeta(c, items, gin.H{"total": total})
}

// @Summary CreateCashMovement
// @Description Crea un nuevo movimiento de caja.
// @Tags cash
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body SwaggerCreateCashMovementRequest true "Create cash movement payload"
// @Success 201 {object} SwaggerCashMovementResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Router /api/v1/cash/movements [post]
func (h *Controller) CreateCashMovement(c *gin.Context) {
	var in cashreq.CreateCashMovementRequestDTO
	if err := c.ShouldBindJSON(&in); err != nil {
		webutils.Fail(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	out, err := h.uc.Cash.CreateMovement(c.Request.Context(), actorID(c), in)
	if err != nil {
		webutils.Fail(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	webutils.Created(c, out)
}

// @Summary GetCashMovement
// @Description Obtiene un movimiento de caja por identificador.
// @Tags cash
// @Produce json
// @Security BearerAuth
// @Param id path string true "Resource ID (UUID)"
// @Success 200 {object} SwaggerCashMovementResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Failure 404 {object} SwaggerErrorResponse
// @Router /api/v1/cash/movements/:id [get]
func (h *Controller) GetCashMovement(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	out, err := h.uc.Cash.GetMovement(c.Request.Context(), id)
	if err != nil {
		webutils.Fail(c, http.StatusNotFound, "NOT_FOUND", err.Error())
		return
	}
	webutils.OK(c, out)
}

// @Summary UpdateCashMovement
// @Description Actualiza un movimiento de caja existente.
// @Tags cash
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body SwaggerUpdateCashMovementRequest true "Update cash movement payload"
// @Param id path string true "Resource ID (UUID)"
// @Success 200 {object} SwaggerCashMovementResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Router /api/v1/cash/movements/:id [put]
func (h *Controller) UpdateCashMovement(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	var in cashreq.UpdateCashMovementRequestDTO
	if err := c.ShouldBindJSON(&in); err != nil {
		webutils.Fail(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	out, err := h.uc.Cash.UpdateMovement(c.Request.Context(), actorID(c), id, in)
	if err != nil {
		webutils.Fail(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	webutils.OK(c, out)
}
