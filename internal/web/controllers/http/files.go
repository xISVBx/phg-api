package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
	filesreq "photogallery/api_go/internal/application/dtos/request/files"
	webutils "photogallery/api_go/internal/web/utils"
)

// @Summary ListFiles
// @Description Lista paginada de archivos con búsqueda y orden opcional.
// @Tags files
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page" default(1)
// @Param pageSize query int false "Page size" default(25)
// @Param q query string false "Search text"
// @Param sort query string false "Sort field"
// @Param dir query string false "Sort direction" Enums(asc,desc)
// @Success 200 {object} SwaggerFileListResponse
// @Failure 500 {object} SwaggerErrorResponse
// @Router /api/v1/files [get]
func (h *Controller) ListFiles(c *gin.Context) {
	items, total, err := h.uc.Files.List(c.Request.Context(), parseQuery(c))
	if err != nil {
		webutils.Fail(c, http.StatusInternalServerError, "INTERNAL", err.Error())
		return
	}
	webutils.OKMeta(c, items, gin.H{"total": total})
}

// @Summary UploadFile
// @Description Sube un archivo y lo vincula a una entidad (Sale, WorkOrder, Appointment, etc.).
// @Tags files
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "File"
// @Param entityType formData string true "Entity type"
// @Param entityId formData string true "Entity ID (UUID)"
// @Param customerId formData string false "Customer ID (UUID)"
// @Param notes formData string false "Notes"
// @Success 201 {object} SwaggerFileResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Router /api/v1/files/upload [post]
func (h *Controller) UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		webutils.Fail(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	var link filesreq.FileUploadLinkDTO
	if err := c.ShouldBind(&link); err != nil {
		webutils.Fail(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	out, err := h.uc.Files.Upload(c.Request.Context(), actorID(c), file, link)
	if err != nil {
		webutils.Fail(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	webutils.Created(c, out)
}

// @Summary DownloadFile
// @Description Descarga un archivo previamente cargado por identificador.
// @Tags files
// @Produce application/octet-stream
// @Security BearerAuth
// @Param id path string true "Resource ID (UUID)"
// @Success 200 {file} file
// @Failure 400 {object} SwaggerErrorResponse
// @Failure 404 {object} SwaggerErrorResponse
// @Router /api/v1/files/:id/download [get]
func (h *Controller) DownloadFile(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	path, file, err := h.uc.Files.ResolveDownloadPath(c.Request.Context(), id)
	if err != nil {
		webutils.Fail(c, http.StatusNotFound, "NOT_FOUND", err.Error())
		return
	}
	c.FileAttachment(path, file.OriginalName)
}

// @Summary DeleteFile
// @Description Elimina un archivo por identificador.
// @Tags files
// @Produce json
// @Security BearerAuth
// @Param id path string true "Resource ID (UUID)"
// @Success 200 {object} SwaggerOKFlagResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Router /api/v1/files/:id [delete]
func (h *Controller) DeleteFile(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	if err := h.uc.Files.Delete(c.Request.Context(), actorID(c), id); err != nil {
		webutils.Fail(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	webutils.OK(c, gin.H{"ok": true})
}
