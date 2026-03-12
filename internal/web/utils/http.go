package utils

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type ProblemDetails struct {
	Type     string   `json:"type"`
	Title    string   `json:"title"`
	Status   int      `json:"status"`
	Detail   string   `json:"detail"`
	Instance string   `json:"instance,omitempty"`
	Code     string   `json:"code,omitempty"`
	Errors   []string `json:"errors,omitempty"`
}

func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func OKMeta(c *gin.Context, data any, meta any) {
	c.JSON(http.StatusOK, gin.H{"data": data, "meta": meta})
}

func Created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, gin.H{"data": data})
}

func Fail(c *gin.Context, status int, code, msg string, details ...string) {
	pd := ProblemDetails{
		Type:     problemType(status, code),
		Title:    problemTitle(status, code),
		Status:   status,
		Detail:   msg,
		Instance: c.Request.URL.Path,
		Code:     code,
		Errors:   details,
	}
	c.Header("Content-Type", "application/problem+json")
	c.AbortWithStatusJSON(status, pd)
}

func problemType(status int, code string) string {
	if strings.TrimSpace(code) != "" {
		return "https://api.photogallery.local/problems/" + strings.ToLower(strings.ReplaceAll(code, "_", "-"))
	}
	return "about:blank"
}

func problemTitle(status int, code string) string {
	if strings.TrimSpace(code) == "" {
		return http.StatusText(status)
	}
	switch strings.ToUpper(code) {
	case "VALIDATION_ERROR":
		return "Validation Error"
	case "BAD_REQUEST":
		return "Bad Request"
	case "UNAUTHORIZED":
		return "Unauthorized"
	case "FORBIDDEN":
		return "Forbidden"
	case "NOT_FOUND":
		return "Not Found"
	case "INTERNAL":
		return "Internal Server Error"
	default:
		return http.StatusText(status)
	}
}
