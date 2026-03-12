package http

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	securityreq "photogallery/api_go/internal/application/dtos/request/security"
	"photogallery/api_go/internal/application/mappers"
	"photogallery/api_go/internal/domain/entities"
	webutils "photogallery/api_go/internal/web/utils"
)

// @Summary ListUsers
// @Description Lista paginada de usuarios con búsqueda y orden opcional.
// @Tags security
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page" default(1)
// @Param pageSize query int false "Page size" default(25)
// @Param q query string false "Search text"
// @Param sort query string false "Sort field"
// @Param dir query string false "Sort direction" Enums(asc,desc)
// @Success 200 {object} SwaggerUserListResponse
// @Failure 500 {object} SwaggerErrorResponse
// @Router /api/v1/users [get]
func (h *Controller) ListUsers(c *gin.Context) {
	items, total, err := h.uc.Security.ListUsers(c.Request.Context(), parseQuery(c))
	if err != nil {
		webutils.Fail(c, http.StatusInternalServerError, "INTERNAL", err.Error())
		return
	}
	webutils.OKMeta(c, mappers.ToUserResponseList(items), gin.H{"total": total})
}

// @Summary CreateUser
// @Description Crea un nuevo usuario.
// @Tags security
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body SwaggerCreateUserRequest true "Create user payload"
// @Success 201 {object} SwaggerUserResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Router /api/v1/users [post]
func (h *Controller) CreateUser(c *gin.Context) {
	var in securityreq.CreateUserRequestDTO
	if err := c.ShouldBindJSON(&in); err != nil {
		webutils.Fail(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	out, err := h.uc.Security.CreateUser(c.Request.Context(), actorID(c), in)
	if err != nil {
		webutils.Fail(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	webutils.Created(c, mappers.ToUserResponse(out))
}

// @Summary GetUser
// @Description Obtiene un usuario por identificador.
// @Tags security
// @Produce json
// @Security BearerAuth
// @Param id path string true "Resource ID (UUID)"
// @Success 200 {object} SwaggerUserResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Failure 404 {object} SwaggerErrorResponse
// @Router /api/v1/users/:id [get]
func (h *Controller) GetUser(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	out, err := h.uc.Security.GetUser(c.Request.Context(), id)
	if err != nil {
		webutils.Fail(c, http.StatusNotFound, "NOT_FOUND", err.Error())
		return
	}
	webutils.OK(c, out)
}

// @Summary UpdateUser
// @Description Actualiza un usuario existente.
// @Tags security
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body SwaggerUpdateUserRequest true "Update user payload"
// @Param id path string true "Resource ID (UUID)"
// @Success 200 {object} SwaggerUserResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Router /api/v1/users/:id [put]
func (h *Controller) UpdateUser(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	var in securityreq.UpdateUserRequestDTO
	if err := c.ShouldBindJSON(&in); err != nil {
		webutils.Fail(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	out, err := h.uc.Security.UpdateUser(c.Request.Context(), actorID(c), id, in)
	if err != nil {
		webutils.Fail(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	webutils.OK(c, mappers.ToUserResponse(out))
}

// @Summary ActivateUser
// @Description Activa un usuario.
// @Tags security
// @Produce json
// @Security BearerAuth
// @Param id path string true "Resource ID (UUID)"
// @Success 200 {object} SwaggerOKFlagResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Router /api/v1/users/:id/activate [patch]
func (h *Controller) ActivateUser(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	if err := h.uc.Security.SetUserActive(c.Request.Context(), actorID(c), id, true); err != nil {
		webutils.Fail(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	webutils.OK(c, gin.H{"ok": true})
}

// @Summary DeactivateUser
// @Description Desactiva un usuario.
// @Tags security
// @Produce json
// @Security BearerAuth
// @Param id path string true "Resource ID (UUID)"
// @Success 200 {object} SwaggerOKFlagResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Router /api/v1/users/:id/deactivate [patch]
func (h *Controller) DeactivateUser(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	if err := h.uc.Security.SetUserActive(c.Request.Context(), actorID(c), id, false); err != nil {
		webutils.Fail(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	webutils.OK(c, gin.H{"ok": true})
}

// @Summary SetUserPassword
// @Description Actualiza la contraseña de un usuario.
// @Tags security
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body SwaggerSetPasswordRequest true "Set password payload"
// @Param id path string true "Resource ID (UUID)"
// @Success 200 {object} SwaggerOKFlagResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Router /api/v1/users/:id/password [patch]
func (h *Controller) SetUserPassword(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	var in securityreq.SetPasswordRequestDTO
	if err := c.ShouldBindJSON(&in); err != nil {
		webutils.Fail(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	if err := h.uc.Security.SetAdminPassword(c.Request.Context(), actorID(c), id, in.NewPassword); err != nil {
		webutils.Fail(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	webutils.OK(c, gin.H{"ok": true})
}

// @Summary GetUserRoles
// @Description Obtiene los roles asignados a un usuario.
// @Tags security
// @Produce json
// @Security BearerAuth
// @Param id path string true "Resource ID (UUID)"
// @Success 200 {object} SwaggerUserRolesResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Router /api/v1/users/:id/roles [get]
func (h *Controller) GetUserRoles(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	out, err := h.uc.Security.GetUserRoles(c.Request.Context(), id)
	if err != nil {
		webutils.Fail(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	webutils.OK(c, out)
}

// @Summary SetUserRoles
// @Description Reemplaza los roles asignados a un usuario.
// @Tags security
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body SwaggerSetPrimaryRoleRequest true "Set primary role payload"
// @Param id path string true "Resource ID (UUID)"
// @Success 200 {object} SwaggerOKFlagResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Router /api/v1/users/:id/roles [put]
func (h *Controller) SetUserRoles(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	var in securityreq.SetPrimaryRoleRequestDTO
	if err := c.ShouldBindJSON(&in); err != nil {
		webutils.Fail(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	rid, err := uuid.Parse(in.PrimaryRoleID)
	if err != nil {
		webutils.Fail(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	if err := h.uc.Security.SetUserRole(c.Request.Context(), actorID(c), id, rid); err != nil {
		webutils.Fail(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	webutils.OK(c, gin.H{"ok": true})
}

// @Summary ListUserOverrides
// @Description Lista paginada de overrides de usuario con búsqueda y orden opcional.
// @Tags security
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page" default(1)
// @Param pageSize query int false "Page size" default(25)
// @Param q query string false "Search text"
// @Param sort query string false "Sort field"
// @Param dir query string false "Sort direction" Enums(asc,desc)
// @Param id path string true "Resource ID (UUID)"
// @Success 200 {object} SwaggerOverrideListResponse
// @Failure 500 {object} SwaggerErrorResponse
// @Router /api/v1/users/:id/overrides [get]
func (h *Controller) ListUserOverrides(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	out, err := h.uc.Security.ListOverrides(c.Request.Context(), id)
	if err != nil {
		webutils.Fail(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	webutils.OK(c, out)
}

// @Summary ReplaceUserOverrides
// @Description Reemplaza en bloque los overrides de un usuario.
// @Tags security
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body SwaggerReplaceOverridesRequest true "Replace overrides payload"
// @Param id path string true "Resource ID (UUID)"
// @Success 200 {object} SwaggerOKFlagResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Router /api/v1/users/:id/overrides [put]
func (h *Controller) ReplaceUserOverrides(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	var in securityreq.ReplaceOverridesRequestDTO
	if err := c.ShouldBindJSON(&in); err != nil {
		webutils.Fail(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	items := make([]entities.UserPermissionOverride, 0, len(in.Items))
	for _, it := range in.Items {
		sid, _ := uuid.Parse(it.SubMenuID)
		pid, _ := uuid.Parse(it.PermissionID)
		items = append(items, entities.UserPermissionOverride{UserID: id, SubMenuID: sid, PermissionID: pid, Mode: it.Mode})
	}
	if err := h.uc.Security.ReplaceOverrides(c.Request.Context(), actorID(c), id, items); err != nil {
		webutils.Fail(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	webutils.OK(c, gin.H{"ok": true})
}

// @Summary CreateUserOverride
// @Description Crea un override para un usuario.
// @Tags security
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body SwaggerOverrideItemRequest true "Create override payload"
// @Param id path string true "Resource ID (UUID)"
// @Success 201 {object} SwaggerOverrideResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Router /api/v1/users/:id/overrides [post]
func (h *Controller) CreateUserOverride(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	var in securityreq.OverrideItemDTO
	if err := c.ShouldBindJSON(&in); err != nil {
		webutils.Fail(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	sid, _ := uuid.Parse(in.SubMenuID)
	pid, _ := uuid.Parse(in.PermissionID)
	item := &entities.UserPermissionOverride{UserID: id, SubMenuID: sid, PermissionID: pid, Mode: in.Mode}
	if err := h.uc.Security.CreateOverride(c.Request.Context(), actorID(c), item); err != nil {
		webutils.Fail(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	webutils.Created(c, item)
}

// @Summary DeleteUserOverride
// @Description Elimina un override de usuario por identificador.
// @Tags security
// @Produce json
// @Security BearerAuth
// @Param id path string true "Resource ID (UUID)"
// @Param overrideId path string true "Override ID (UUID)"
// @Success 200 {object} SwaggerOKFlagResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Router /api/v1/users/:id/overrides/:overrideId [delete]
func (h *Controller) DeleteUserOverride(c *gin.Context) {
	uid, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	oid, ok := parseUUID(c, "overrideId")
	if !ok {
		return
	}
	if err := h.uc.Security.DeleteOverride(c.Request.Context(), actorID(c), uid, oid); err != nil {
		webutils.Fail(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	webutils.OK(c, gin.H{"ok": true})
}
