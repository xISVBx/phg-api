package web

import (
	"github.com/gin-gonic/gin"

	httpctl "photogallery/api_go/internal/web/controllers/http"
)

func registerSecurityRoutes(secure *gin.RouterGroup, c *httpctl.Controller) {
	secure.GET("/users", c.ListUsers)
	secure.POST("/users", c.CreateUser)
	secure.GET("/users/:id", c.GetUser)
	secure.PUT("/users/:id", c.UpdateUser)
	secure.PATCH("/users/:id/activate", c.ActivateUser)
	secure.PATCH("/users/:id/deactivate", c.DeactivateUser)
	secure.PATCH("/users/:id/password", c.SetUserPassword)
	secure.GET("/users/:id/roles", c.GetUserRoles)
	secure.PUT("/users/:id/roles", c.SetUserRoles)
	secure.GET("/users/:id/overrides", c.ListUserOverrides)
	secure.PUT("/users/:id/overrides", c.ReplaceUserOverrides)
	secure.POST("/users/:id/overrides", c.CreateUserOverride)
	secure.DELETE("/users/:id/overrides/:overrideId", c.DeleteUserOverride)

	secure.GET("/roles", c.ListRoles)
	secure.POST("/roles", c.CreateRole)
	secure.GET("/roles/:id", c.GetRole)
	secure.PUT("/roles/:id", c.UpdateRole)
	secure.PATCH("/roles/:id/activate", c.ActivateRole)
	secure.PATCH("/roles/:id/deactivate", c.DeactivateRole)
	secure.GET("/roles/:id/permissions", c.GetRolePermissions)
	secure.PUT("/roles/:id/permissions", c.SetRolePermissions)

	secure.GET("/menus", c.ListMenus)
	secure.POST("/menus", c.CreateMenu)
	secure.GET("/menus/:id", c.GetMenu)
	secure.PUT("/menus/:id", c.UpdateMenu)
	secure.PATCH("/menus/:id/activate", c.ActivateMenu)
	secure.PATCH("/menus/:id/deactivate", c.DeactivateMenu)

	secure.GET("/submenus", c.ListSubMenus)
	secure.POST("/submenus", c.CreateSubMenu)
	secure.GET("/submenus/:id", c.GetSubMenu)
	secure.PUT("/submenus/:id", c.UpdateSubMenu)
	secure.PATCH("/submenus/:id/activate", c.ActivateSubMenu)
	secure.PATCH("/submenus/:id/deactivate", c.DeactivateSubMenu)

	secure.GET("/permissions", c.ListPermissions)
	secure.POST("/permissions", c.CreatePermission)
	secure.GET("/permissions/:id", c.GetPermission)
	secure.PUT("/permissions/:id", c.UpdatePermission)
}
