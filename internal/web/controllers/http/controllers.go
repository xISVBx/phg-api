package http

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"photogallery/api_go/internal/application/use_cases"
	drepo "photogallery/api_go/internal/domain/repositories"
	webutils "photogallery/api_go/internal/web/utils"
	"strconv"
)

type Controller struct {
	uc *use_cases.UseCases
}

func NewController(uc *use_cases.UseCases) *Controller { return &Controller{uc: uc} }
func parseUUID(c *gin.Context, name string) (uuid.UUID, bool) {
	id, err := uuid.Parse(c.Param(name))
	if err != nil {
		webutils.Fail(c, 400, "VALIDATION_ERROR", "invalid uuid")
		return uuid.Nil, false
	}
	return id, true
}
func parseQuery(c *gin.Context) drepo.QueryOptions {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "25"))
	return drepo.QueryOptions{Page: page, PageSize: pageSize, Q: c.Query("q"), Sort: c.Query("sort"), Dir: c.Query("dir")}
}
func actorID(c *gin.Context) uuid.UUID {
	id, ok := webutils.MustUserID(c)
	if !ok {
		return uuid.Nil
	}
	return id
}
