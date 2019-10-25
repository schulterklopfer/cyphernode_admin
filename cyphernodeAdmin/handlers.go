package cyphernodeAdmin

import (
  "github.com/schulterklopfer/cyphernode_admin/globals"
  "github.com/schulterklopfer/cyphernode_admin/handlers"
)

func (cyphernodeAdmin *CyphernodeAdmin) initDefaultHandlers() {
  cyphernodeAdmin.routerGroups["default"] = cyphernodeAdmin.engine.Group("/")
  {
    // TODO: csrf protection
    //cyphernodeAdmin.routerGroups["oidc"].StaticFile(".well-known/openid-configuration", "./static/well-known/openid-configuration" )
    cyphernodeAdmin.routerGroups["default"].GET( globals.ENDPOINTS_LOGIN, handlers.DefaultLogin )
    cyphernodeAdmin.routerGroups["default"].GET("/logout", handlers.DefaultLogout )
    cyphernodeAdmin.routerGroups["default"].GET( globals.ENDPOINTS_CALLBACK, handlers.DefaultCallbackGet)
    cyphernodeAdmin.routerGroups["default"].GET( globals.ENDPOINTS_HOME, handlers.DefaultHome )

  }
}

/*
func (cyphernodeAdmin *CyphernodeAdmin) initOIDCHandlers() {
  cyphernodeAdmin.routerGroups["oidc"] = cyphernodeAdmin.engine.Group("/oidc")
  {
    cyphernodeAdmin.routerGroups["oidc"].StaticFile("/jwks.json", "./static/jwks.json" )
    cyphernodeAdmin.routerGroups["oidc"].GET("/me", handlers.GetMe )
  }
}
*/

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
    cyphernodeAdmin.routerGroups["roles"].GET("/:clientID/:username", handlers.GetRolesForApp)
  }
}

func (cyphernodeAdmin *CyphernodeAdmin) initHydraHandlers() {
  cyphernodeAdmin.routerGroups["hydra"] = cyphernodeAdmin.engine.Group("/hydra")
  {
    // TODO: csrf protection
    cyphernodeAdmin.routerGroups["hydra"].GET("/consent", handlers.HydraConsentGet)
    cyphernodeAdmin.routerGroups["hydra"].GET("/login", handlers.HydraLoginGet)
    cyphernodeAdmin.routerGroups["hydra"].GET("/logout", handlers.HydraLogoutGet)
    cyphernodeAdmin.routerGroups["hydra"].POST("/consent", handlers.HydraConsentPost)
    cyphernodeAdmin.routerGroups["hydra"].POST("/login", handlers.HydraLoginPost)
    cyphernodeAdmin.routerGroups["hydra"].POST("/logout", handlers.HydraLogoutPost)
  }
}

func (cyphernodeAdmin *CyphernodeAdmin) initSessionHandlers() {
  cyphernodeAdmin.routerGroups["sessions"] = cyphernodeAdmin.engine.Group("/sessions" )
  {
    cyphernodeAdmin.routerGroups["sessions"].GET("/:sessionID", handlers.GetSession )
    cyphernodeAdmin.routerGroups["sessions"].PATCH("/:sessionID", handlers.PatchSession )
    cyphernodeAdmin.routerGroups["sessions"].DELETE("/:sessionID", handlers.DeleteSession )
    cyphernodeAdmin.routerGroups["sessions"].POST("/", handlers.CreateSession )
  }
}


