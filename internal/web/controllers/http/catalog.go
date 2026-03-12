package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
	catalogreq "photogallery/api_go/internal/application/dtos/request/catalog"
	webutils "photogallery/api_go/internal/web/utils"
)

// @Summary ListCategories
// @Description Lista paginada de categorías con búsqueda y orden opcional.
// @Tags catalog
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page" default(1)
// @Param pageSize query int false "Page size" default(25)
// @Param q query string false "Search text"
// @Param sort query string false "Sort field"
// @Param dir query string false "Sort direction" Enums(asc,desc)
// @Success 200 {object} SwaggerCategoryListResponse
// @Failure 500 {object} SwaggerErrorResponse
// @Router /api/v1/categories [get]
func (h *Controller) ListCategories(c *gin.Context) {
	items, total, err := h.uc.Catalog.ListCategories(c.Request.Context(), parseQuery(c))
	if err != nil {
		webutils.Fail(c, http.StatusInternalServerError, "INTERNAL", err.Error())
		return
	}
	webutils.OKMeta(c, items, gin.H{"total": total})
}

// @Summary CreateCategory
// @Description Crea una nueva categoría.
// @Tags catalog
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body SwaggerCreateCategoryRequest true "Create category payload"
// @Success 201 {object} SwaggerCategoryResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Router /api/v1/categories [post]
func (h *Controller) CreateCategory(c *gin.Context) {
	var in catalogreq.CreateCategoryRequestDTO
	if err := c.ShouldBindJSON(&in); err != nil {
		webutils.Fail(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	out, err := h.uc.Catalog.CreateCategory(c.Request.Context(), actorID(c), in)
	if err != nil {
		webutils.Fail(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	webutils.Created(c, out)
}

// @Summary GetCategory
// @Description Obtiene una categoría por identificador.
// @Tags catalog
// @Produce json
// @Security BearerAuth
// @Param id path string true "Resource ID (UUID)"
// @Success 200 {object} SwaggerCategoryResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Failure 404 {object} SwaggerErrorResponse
// @Router /api/v1/categories/:id [get]
func (h *Controller) GetCategory(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	out, err := h.uc.Catalog.GetCategory(c.Request.Context(), id)
	if err != nil {
		webutils.Fail(c, http.StatusNotFound, "NOT_FOUND", err.Error())
		return
	}
	webutils.OK(c, out)
}

// @Summary UpdateCategory
// @Description Actualiza una categoría existente.
// @Tags catalog
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body SwaggerUpdateCategoryRequest true "Update category payload"
// @Param id path string true "Resource ID (UUID)"
// @Success 200 {object} SwaggerCategoryResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Router /api/v1/categories/:id [put]
func (h *Controller) UpdateCategory(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	var in catalogreq.UpdateCategoryRequestDTO
	if err := c.ShouldBindJSON(&in); err != nil {
		webutils.Fail(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	out, err := h.uc.Catalog.UpdateCategory(c.Request.Context(), actorID(c), id, in)
	if err != nil {
		webutils.Fail(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	webutils.OK(c, out)
}

// @Summary ActivateCategory
// @Description Activa una categoría.
// @Tags catalog
// @Produce json
// @Security BearerAuth
// @Param id path string true "Resource ID (UUID)"
// @Success 200 {object} SwaggerOKFlagResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Router /api/v1/categories/:id/activate [patch]
func (h *Controller) ActivateCategory(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	if err := h.uc.Catalog.SetCategoryActive(c.Request.Context(), actorID(c), id, true); err != nil {
		webutils.Fail(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	webutils.OK(c, gin.H{"ok": true})
}

// @Summary DeactivateCategory
// @Description Desactiva una categoría.
// @Tags catalog
// @Produce json
// @Security BearerAuth
// @Param id path string true "Resource ID (UUID)"
// @Success 200 {object} SwaggerOKFlagResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Router /api/v1/categories/:id/deactivate [patch]
func (h *Controller) DeactivateCategory(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	if err := h.uc.Catalog.SetCategoryActive(c.Request.Context(), actorID(c), id, false); err != nil {
		webutils.Fail(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	webutils.OK(c, gin.H{"ok": true})
}

// @Summary ListProducts
// @Description Lista paginada de productos con búsqueda y orden opcional.
// @Tags catalog
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page" default(1)
// @Param pageSize query int false "Page size" default(25)
// @Param q query string false "Search text"
// @Param sort query string false "Sort field"
// @Param dir query string false "Sort direction" Enums(asc,desc)
// @Success 200 {object} SwaggerProductListResponse
// @Failure 500 {object} SwaggerErrorResponse
// @Router /api/v1/products [get]
func (h *Controller) ListProducts(c *gin.Context) {
	items, total, err := h.uc.Catalog.ListProducts(c.Request.Context(), parseQuery(c))
	if err != nil {
		webutils.Fail(c, http.StatusInternalServerError, "INTERNAL", err.Error())
		return
	}
	webutils.OKMeta(c, items, gin.H{"total": total})
}

// @Summary CreateProduct
// @Description Crea un nuevo producto.
// @Tags catalog
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body SwaggerCreateProductRequest true "Create product payload"
// @Success 201 {object} SwaggerProductResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Router /api/v1/products [post]
func (h *Controller) CreateProduct(c *gin.Context) {
	var in catalogreq.CreateProductRequestDTO
	if err := c.ShouldBindJSON(&in); err != nil {
		webutils.Fail(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	out, err := h.uc.Catalog.CreateProduct(c.Request.Context(), actorID(c), in)
	if err != nil {
		webutils.Fail(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	webutils.Created(c, out)
}

// @Summary GetProduct
// @Description Obtiene un producto por identificador.
// @Tags catalog
// @Produce json
// @Security BearerAuth
// @Param id path string true "Resource ID (UUID)"
// @Success 200 {object} SwaggerProductResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Failure 404 {object} SwaggerErrorResponse
// @Router /api/v1/products/:id [get]
func (h *Controller) GetProduct(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	out, err := h.uc.Catalog.GetProduct(c.Request.Context(), id)
	if err != nil {
		webutils.Fail(c, http.StatusNotFound, "NOT_FOUND", err.Error())
		return
	}
	webutils.OK(c, out)
}

// @Summary UpdateProduct
// @Description Actualiza un producto existente.
// @Tags catalog
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body SwaggerUpdateProductRequest true "Update product payload"
// @Param id path string true "Resource ID (UUID)"
// @Success 200 {object} SwaggerProductResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Router /api/v1/products/:id [put]
func (h *Controller) UpdateProduct(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	var in catalogreq.UpdateProductRequestDTO
	if err := c.ShouldBindJSON(&in); err != nil {
		webutils.Fail(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	out, err := h.uc.Catalog.UpdateProduct(c.Request.Context(), actorID(c), id, in)
	if err != nil {
		webutils.Fail(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	webutils.OK(c, out)
}

// @Summary ActivateProduct
// @Description Activa un producto.
// @Tags catalog
// @Produce json
// @Security BearerAuth
// @Param id path string true "Resource ID (UUID)"
// @Success 200 {object} SwaggerOKFlagResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Router /api/v1/products/:id/activate [patch]
func (h *Controller) ActivateProduct(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	if err := h.uc.Catalog.SetProductActive(c.Request.Context(), actorID(c), id, true); err != nil {
		webutils.Fail(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	webutils.OK(c, gin.H{"ok": true})
}

// @Summary DeactivateProduct
// @Description Desactiva un producto.
// @Tags catalog
// @Produce json
// @Security BearerAuth
// @Param id path string true "Resource ID (UUID)"
// @Success 200 {object} SwaggerOKFlagResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Router /api/v1/products/:id/deactivate [patch]
func (h *Controller) DeactivateProduct(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	if err := h.uc.Catalog.SetProductActive(c.Request.Context(), actorID(c), id, false); err != nil {
		webutils.Fail(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	webutils.OK(c, gin.H{"ok": true})
}
