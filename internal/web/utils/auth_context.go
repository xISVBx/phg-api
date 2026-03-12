package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const userIDKey = "auth_user_id"

func SetUserID(c *gin.Context, id uuid.UUID) {
	c.Set(userIDKey, id)
}

func MustUserID(c *gin.Context) (uuid.UUID, bool) {
	v, ok := c.Get(userIDKey)
	if !ok {
		return uuid.Nil, false
	}
	id, ok := v.(uuid.UUID)
	return id, ok
}
