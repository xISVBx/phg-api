package mappers

import (
	"strings"

	"photogallery/api_go/internal/domain/entities"
)

func CustomerKey(c *entities.Customer) string {
	key := strings.TrimSpace(strings.ToLower(c.Email))
	if key == "" {
		key = strings.TrimSpace(c.CustomerCode)
	}
	key = strings.ReplaceAll(key, "@", "_")
	key = strings.ReplaceAll(key, ".", "_")
	key = strings.ReplaceAll(key, " ", "_")
	if key == "" {
		key = "unknown_customer"
	}
	return key
}
