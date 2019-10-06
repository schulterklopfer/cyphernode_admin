package cyphernodeAdmin

import (
	"github.com/schulterklopfer/cyphernode_admin/handlers"
)

func (cyphernodeAdmin *CyphernodeAdmin) initUsersHandlers() {
	cyphernodeAdmin.routerGroups["users"] = cyphernodeAdmin.engine.Group("/api/v0/users")
	{
		cyphernodeAdmin.routerGroups["users"].GET("/", handlers.FindUsers)
		cyphernodeAdmin.routerGroups["users"].POST("/", handlers.CreateUser)
		cyphernodeAdmin.routerGroups["users"].GET("/:id", handlers.GetUser)
		cyphernodeAdmin.routerGroups["users"].PUT("/:id", handlers.UpdateUser )
		cyphernodeAdmin.routerGroups["users"].PATCH("/:id", handlers.PatchUser )
		cyphernodeAdmin.routerGroups["users"].DELETE("/:id", handlers.DeleteUser )
		cyphernodeAdmin.routerGroups["users"].POST("/:id/roles", handlers.UserAddRoles )
		cyphernodeAdmin.routerGroups["users"].DELETE("/:id/roles/:roleId", handlers.UserRemoveRole )
	}
}

func (cyphernodeAdmin *CyphernodeAdmin) initAppsHandlers() {
	cyphernodeAdmin.routerGroups["apps"] = cyphernodeAdmin.engine.Group("/api/v0/apps")
	{
		cyphernodeAdmin.routerGroups["apps"].GET("/", handlers.FindApps)
		cyphernodeAdmin.routerGroups["apps"].POST("/", handlers.CreateApp)
		cyphernodeAdmin.routerGroups["apps"].GET("/:id", handlers.GetApp)
		cyphernodeAdmin.routerGroups["apps"].PUT("/:id", handlers.UpdateApp )
		cyphernodeAdmin.routerGroups["apps"].PATCH("/:id", handlers.PatchApp )
		cyphernodeAdmin.routerGroups["apps"].DELETE("/:id", handlers.DeleteApp )
		cyphernodeAdmin.routerGroups["apps"].POST("/:id/roles", handlers.AppAddRoles )
		cyphernodeAdmin.routerGroups["apps"].DELETE("/:id/roles/:roleId", handlers.AppRemoveRole )
	}
}

func (cyphernodeAdmin *CyphernodeAdmin) initRolesHandlers() {
	cyphernodeAdmin.routerGroups["roles"] = cyphernodeAdmin.engine.Group("/api/v0/roles")
	{
		cyphernodeAdmin.routerGroups["roles"].GET("/:hash/:username", handlers.GetRolesForApp)
	}
}

func (cyphernodeAdmin *CyphernodeAdmin) initHydraHandlers() {
	cyphernodeAdmin.routerGroups["hydra"] = cyphernodeAdmin.engine.Group("/hydra")
	{
		// TODO: csrf protection
		cyphernodeAdmin.routerGroups["hydra"].GET("/consent", handlers.GetHydraConsent )
		cyphernodeAdmin.routerGroups["hydra"].GET("/login", handlers.GetHydraLogin )
		cyphernodeAdmin.routerGroups["hydra"].GET("/logout", handlers.GetHydraLogout )
		cyphernodeAdmin.routerGroups["hydra"].POST("/consent", handlers.PostHydraConsent )
		cyphernodeAdmin.routerGroups["hydra"].POST("/login", handlers.PostHydraLogin )
		cyphernodeAdmin.routerGroups["hydra"].POST("/logout", handlers.PostHydraLogout )
	}
}
