package http

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	securityreq "photogallery/api_go/internal/application/dtos/request/security"
	"photogallery/api_go/internal/domain/entities"
	webutils "photogallery/api_go/internal/web/utils"
)

// @Summary ListRoles
// @Description Lista paginada de roles con búsqueda y orden opcional.
// @Tags security
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page" default(1)
// @Param pageSize query int false "Page size" default(25)
// @Param q query string false "Search text"
// @Param sort query string false "Sort field"
// @Param dir query string false "Sort direction" Enums(asc,desc)
// @Success 200 {object} SwaggerRoleListResponse
// @Failure 500 {object} SwaggerErrorResponse
// @Router /api/v1/roles [get]
func (h *Controller) ListRoles(c *gin.Context) {
	items, total, err := h.uc.Security.ListRoles(c.Request.Context(), parseQuery(c))
	if err != nil {
		webutils.Fail(c, http.StatusInternalServerError, "INTERNAL", err.Error())
		return
	}
	webutils.OKMeta(c, items, gin.H{"total": total})
}

// @Summary CreateRole
// @Description Crea un nuevo rol.
// @Tags security
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body SwaggerCreateRoleRequest true "Role payload"
// @Success 201 {object} SwaggerRoleDetailResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Router /api/v1/roles [post]
func (h *Controller) CreateRole(c *gin.Context) {
	var in securityreq.CreateRoleRequestDTO
	if err := c.ShouldBindJSON(&in); err != nil {
		webutils.Fail(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	out, err := h.uc.Security.CreateRole(c.Request.Context(), actorID(c), in)
	if err != nil {
		webutils.Fail(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	webutils.Created(c, out)
}

// @Summary GetRole
// @Description Obtiene un rol por identificador.
// @Tags security
// @Produce json
// @Security BearerAuth
// @Param id path string true "Resource ID (UUID)"
// @Success 200 {object} SwaggerRoleDetailResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Failure 404 {object} SwaggerErrorResponse
// @Router /api/v1/roles/:id [get]
func (h *Controller) GetRole(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	out, err := h.uc.Security.GetRole(c.Request.Context(), id)
	if err != nil {
		webutils.Fail(c, http.StatusNotFound, "NOT_FOUND", err.Error())
		return
	}
	webutils.OK(c, out)
}

// @Summary UpdateRole
// @Description Actualiza un rol existente.
// @Tags security
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body SwaggerUpdateRoleRequest true "Role payload"
// @Param id path string true "Resource ID (UUID)"
// @Success 200 {object} SwaggerRoleDetailResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Router /api/v1/roles/:id [put]
func (h *Controller) UpdateRole(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	var in securityreq.UpdateRoleRequestDTO
	if err := c.ShouldBindJSON(&in); err != nil {
		webutils.Fail(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	out, err := h.uc.Security.UpdateRole(c.Request.Context(), actorID(c), id, in)
	if err != nil {
		webutils.Fail(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	webutils.OK(c, out)
}

// @Summary ActivateRole
// @Description Activa un rol.
// @Tags security
// @Produce json
// @Security BearerAuth
// @Param id path string true "Resource ID (UUID)"
// @Success 200 {object} SwaggerOKFlagResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Router /api/v1/roles/:id/activate [patch]
func (h *Controller) ActivateRole(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	if err := h.uc.Security.SetRoleActive(c.Request.Context(), actorID(c), id, true); err != nil {
		webutils.Fail(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	webutils.OK(c, gin.H{"ok": true})
}

// @Summary DeactivateRole
// @Description Desactiva un rol.
// @Tags security
// @Produce json
// @Security BearerAuth
// @Param id path string true "Resource ID (UUID)"
// @Success 200 {object} SwaggerOKFlagResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Router /api/v1/roles/:id/deactivate [patch]
func (h *Controller) DeactivateRole(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	if err := h.uc.Security.SetRoleActive(c.Request.Context(), actorID(c), id, false); err != nil {
		webutils.Fail(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	webutils.OK(c, gin.H{"ok": true})
}

// @Summary GetRolePermissions
// @Description Obtiene los permisos asignados a un rol.
// @Tags security
// @Produce json
// @Security BearerAuth
// @Param id path string true "Resource ID (UUID)"
// @Success 200 {object} SwaggerRolePermissionsResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Router /api/v1/roles/:id/permissions [get]
func (h *Controller) GetRolePermissions(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	out, err := h.uc.Security.RolePermissions(c.Request.Context(), id)
	if err != nil {
		webutils.Fail(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	webutils.OK(c, out)
}

// @Summary SetRolePermissions
// @Description Reemplaza los permisos asignados a un rol.
// @Tags security
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body SwaggerReplaceRolePermissionsRequest true "Role permissions payload"
// @Param id path string true "Resource ID (UUID)"
// @Success 200 {object} SwaggerOKFlagResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Router /api/v1/roles/:id/permissions [put]
func (h *Controller) SetRolePermissions(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	var in securityreq.ReplaceRolePermissionsRequestDTO
	if err := c.ShouldBindJSON(&in); err != nil {
		webutils.Fail(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	items := make([]entities.RoleSubMenuPermission, 0)
	for _, it := range in.Items {
		sid, _ := uuid.Parse(it.SubMenuID)
		for _, pidRaw := range it.PermissionIDs {
			pid, _ := uuid.Parse(pidRaw)
			items = append(items, entities.RoleSubMenuPermission{RoleID: id, SubMenuID: sid, PermissionID: pid})
		}
	}
	if err := h.uc.Security.ReplaceRolePermissions(c.Request.Context(), actorID(c), id, items); err != nil {
		webutils.Fail(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	webutils.OK(c, gin.H{"ok": true})
}

// @Summary ListMenus
// @Description Lista paginada de menús con búsqueda y orden opcional.
// @Tags security
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page" default(1)
// @Param pageSize query int false "Page size" default(25)
// @Param q query string false "Search text"
// @Param sort query string false "Sort field"
// @Param dir query string false "Sort direction" Enums(asc,desc)
// @Success 200 {object} SwaggerMenuListResponse
// @Failure 500 {object} SwaggerErrorResponse
// @Router /api/v1/menus [get]
func (h *Controller) ListMenus(c *gin.Context) {
	items, total, err := h.uc.Security.ListMenus(c.Request.Context(), parseQuery(c))
	if err != nil {
		webutils.Fail(c, http.StatusInternalServerError, "INTERNAL", err.Error())
		return
	}
	webutils.OKMeta(c, items, gin.H{"total": total})
}

// @Summary CreateMenu
// @Description Crea un nuevo menú.
// @Tags security
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body SwaggerMenuRequest true "Menu payload"
// @Success 201 {object} SwaggerMenuResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Router /api/v1/menus [post]
func (h *Controller) CreateMenu(c *gin.Context) {
	var item entities.Menu
	if err := c.ShouldBindJSON(&item); err != nil {
		webutils.Fail(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	if err := h.uc.Security.CreateMenu(c.Request.Context(), actorID(c), &item); err != nil {
		webutils.Fail(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	webutils.Created(c, item)
}

// @Summary GetMenu
// @Description Obtiene un menú por identificador.
// @Tags security
// @Produce json
// @Security BearerAuth
// @Param id path string true "Resource ID (UUID)"
// @Success 200 {object} SwaggerMenuResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Failure 404 {object} SwaggerErrorResponse
// @Router /api/v1/menus/:id [get]
func (h *Controller) GetMenu(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	out, err := h.uc.Security.GetMenu(c.Request.Context(), id)
	if err != nil {
		webutils.Fail(c, http.StatusNotFound, "NOT_FOUND", err.Error())
		return
	}
	webutils.OK(c, out)
}

// @Summary UpdateMenu
// @Description Actualiza un menú existente.
// @Tags security
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body SwaggerMenuRequest true "Menu payload"
// @Param id path string true "Resource ID (UUID)"
// @Success 200 {object} SwaggerMenuResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Router /api/v1/menus/:id [put]
func (h *Controller) UpdateMenu(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	var item entities.Menu
	if err := c.ShouldBindJSON(&item); err != nil {
		webutils.Fail(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	item.ID = id
	if err := h.uc.Security.UpdateMenu(c.Request.Context(), actorID(c), &item); err != nil {
		webutils.Fail(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	webutils.OK(c, item)
}

// @Summary ActivateMenu
// @Description Activa un menú.
// @Tags security
// @Produce json
// @Security BearerAuth
// @Param id path string true "Resource ID (UUID)"
// @Success 200 {object} SwaggerOKFlagResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Router /api/v1/menus/:id/activate [patch]
func (h *Controller) ActivateMenu(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	if err := h.uc.Security.SetMenuActive(c.Request.Context(), actorID(c), id, true); err != nil {
		webutils.Fail(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	webutils.OK(c, gin.H{"ok": true})
}

// @Summary DeactivateMenu
// @Description Desactiva un menú.
// @Tags security
// @Produce json
// @Security BearerAuth
// @Param id path string true "Resource ID (UUID)"
// @Success 200 {object} SwaggerOKFlagResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Router /api/v1/menus/:id/deactivate [patch]
func (h *Controller) DeactivateMenu(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	if err := h.uc.Security.SetMenuActive(c.Request.Context(), actorID(c), id, false); err != nil {
		webutils.Fail(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	webutils.OK(c, gin.H{"ok": true})
}

// @Summary ListSubMenus
// @Description Lista paginada de submenús con búsqueda y orden opcional.
// @Tags security
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page" default(1)
// @Param pageSize query int false "Page size" default(25)
// @Param q query string false "Search text"
// @Param sort query string false "Sort field"
// @Param dir query string false "Sort direction" Enums(asc,desc)
// @Success 200 {object} SwaggerSubMenuListResponse
// @Failure 500 {object} SwaggerErrorResponse
// @Router /api/v1/submenus [get]
func (h *Controller) ListSubMenus(c *gin.Context) {
	items, total, err := h.uc.Security.ListSubMenus(c.Request.Context(), parseQuery(c))
	if err != nil {
		webutils.Fail(c, http.StatusInternalServerError, "INTERNAL", err.Error())
		return
	}
	webutils.OKMeta(c, items, gin.H{"total": total})
}

// @Summary CreateSubMenu
// @Description Crea un nuevo submenú.
// @Tags security
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body SwaggerSubMenuRequest true "Submenu payload"
// @Success 201 {object} SwaggerSubMenuResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Router /api/v1/submenus [post]
func (h *Controller) CreateSubMenu(c *gin.Context) {
	var item entities.SubMenu
	if err := c.ShouldBindJSON(&item); err != nil {
		webutils.Fail(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	if err := h.uc.Security.CreateSubMenu(c.Request.Context(), actorID(c), &item); err != nil {
		webutils.Fail(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	webutils.Created(c, item)
}

// @Summary GetSubMenu
// @Description Obtiene un submenú por identificador.
// @Tags security
// @Produce json
// @Security BearerAuth
// @Param id path string true "Resource ID (UUID)"
// @Success 200 {object} SwaggerSubMenuResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Failure 404 {object} SwaggerErrorResponse
// @Router /api/v1/submenus/:id [get]
func (h *Controller) GetSubMenu(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	out, err := h.uc.Security.GetSubMenu(c.Request.Context(), id)
	if err != nil {
		webutils.Fail(c, http.StatusNotFound, "NOT_FOUND", err.Error())
		return
	}
	webutils.OK(c, out)
}

// @Summary UpdateSubMenu
// @Description Actualiza un submenú existente.
// @Tags security
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body SwaggerSubMenuRequest true "Submenu payload"
// @Param id path string true "Resource ID (UUID)"
// @Success 200 {object} SwaggerSubMenuResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Router /api/v1/submenus/:id [put]
func (h *Controller) UpdateSubMenu(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	var item entities.SubMenu
	if err := c.ShouldBindJSON(&item); err != nil {
		webutils.Fail(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	item.ID = id
	if err := h.uc.Security.UpdateSubMenu(c.Request.Context(), actorID(c), &item); err != nil {
		webutils.Fail(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	webutils.OK(c, item)
}

// @Summary ActivateSubMenu
// @Description Activa un submenú.
// @Tags security
// @Produce json
// @Security BearerAuth
// @Param id path string true "Resource ID (UUID)"
// @Success 200 {object} SwaggerOKFlagResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Router /api/v1/submenus/:id/activate [patch]
func (h *Controller) ActivateSubMenu(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	if err := h.uc.Security.SetSubMenuActive(c.Request.Context(), actorID(c), id, true); err != nil {
		webutils.Fail(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	webutils.OK(c, gin.H{"ok": true})
}

// @Summary DeactivateSubMenu
// @Description Desactiva un submenú.
// @Tags security
// @Produce json
// @Security BearerAuth
// @Param id path string true "Resource ID (UUID)"
// @Success 200 {object} SwaggerOKFlagResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Router /api/v1/submenus/:id/deactivate [patch]
func (h *Controller) DeactivateSubMenu(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	if err := h.uc.Security.SetSubMenuActive(c.Request.Context(), actorID(c), id, false); err != nil {
		webutils.Fail(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	webutils.OK(c, gin.H{"ok": true})
}

// @Summary ListPermissions
// @Description Lista paginada de permisos con búsqueda y orden opcional.
// @Tags security
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page" default(1)
// @Param pageSize query int false "Page size" default(25)
// @Param q query string false "Search text"
// @Param sort query string false "Sort field"
// @Param dir query string false "Sort direction" Enums(asc,desc)
// @Success 200 {object} SwaggerPermissionListResponse
// @Failure 500 {object} SwaggerErrorResponse
// @Router /api/v1/permissions [get]
func (h *Controller) ListPermissions(c *gin.Context) {
	items, total, err := h.uc.Security.ListPermissions(c.Request.Context(), parseQuery(c))
	if err != nil {
		webutils.Fail(c, http.StatusInternalServerError, "INTERNAL", err.Error())
		return
	}
	webutils.OKMeta(c, items, gin.H{"total": total})
}

// @Summary CreatePermission
// @Description Crea un nuevo permiso.
// @Tags security
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body SwaggerPermissionRequest true "Permission payload"
// @Success 201 {object} SwaggerPermissionResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Router /api/v1/permissions [post]
func (h *Controller) CreatePermission(c *gin.Context) {
	var item entities.Permission
	if err := c.ShouldBindJSON(&item); err != nil {
		webutils.Fail(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	if err := h.uc.Security.CreatePermission(c.Request.Context(), actorID(c), &item); err != nil {
		webutils.Fail(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	webutils.Created(c, item)
}

// @Summary GetPermission
// @Description Obtiene un permiso por identificador.
// @Tags security
// @Produce json
// @Security BearerAuth
// @Param id path string true "Resource ID (UUID)"
// @Success 200 {object} SwaggerPermissionResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Failure 404 {object} SwaggerErrorResponse
// @Router /api/v1/permissions/:id [get]
func (h *Controller) GetPermission(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	out, err := h.uc.Security.GetPermission(c.Request.Context(), id)
	if err != nil {
		webutils.Fail(c, http.StatusNotFound, "NOT_FOUND", err.Error())
		return
	}
	webutils.OK(c, out)
}

// @Summary UpdatePermission
// @Description Actualiza un permiso existente.
// @Tags security
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body SwaggerPermissionRequest true "Permission payload"
// @Param id path string true "Resource ID (UUID)"
// @Success 200 {object} SwaggerPermissionResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Router /api/v1/permissions/:id [put]
func (h *Controller) UpdatePermission(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	var item entities.Permission
	if err := c.ShouldBindJSON(&item); err != nil {
		webutils.Fail(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	item.ID = id
	if err := h.uc.Security.UpdatePermission(c.Request.Context(), actorID(c), &item); err != nil {
		webutils.Fail(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	webutils.OK(c, item)
}
