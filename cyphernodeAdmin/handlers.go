package cyphernodeAdmin

import (
  "github.com/schulterklopfer/cyphernode_admin/globals"
  "github.com/schulterklopfer/cyphernode_admin/handlers"
)

func (cyphernodeAdmin *CyphernodeAdmin) createRouterGroups() {
  for i:=0; i<len( globals.ROUTER_GROUPS ); i++ {
    cyphernodeAdmin.routerGroups[globals.ROUTER_GROUPS[i]] = cyphernodeAdmin.engineExternal.Group(globals.ROUTER_GROUPS_BASE_ENDPOINTS[i])
  }
}

func (cyphernodeAdmin *CyphernodeAdmin) initInternalHandlers() {
  cyphernodeAdmin.engineInternal.POST( globals.INTERNAL_ENDPOINTS_REGISTER_APP, handlers.InternalRegisterApp )
  cyphernodeAdmin.engineInternal.GET( globals.ROUTER_GROUPS_BASE_ENDPOINT_SESSIONS+"/:sessionID", handlers.GetSession )
  cyphernodeAdmin.engineInternal.PATCH( globals.ROUTER_GROUPS_BASE_ENDPOINT_SESSIONS+"/:sessionID", handlers.PatchSession )
  cyphernodeAdmin.engineInternal.DELETE( globals.ROUTER_GROUPS_BASE_ENDPOINT_SESSIONS+"/:sessionID", handlers.DeleteSession )
  cyphernodeAdmin.engineInternal.POST( globals.ROUTER_GROUPS_BASE_ENDPOINT_SESSIONS+"/", handlers.CreateSession )
}

func (cyphernodeAdmin *CyphernodeAdmin) initPublicHandlers() {
  cyphernodeAdmin.routerGroups[globals.ROUTER_GROUPS_PUBLIC].GET("/", handlers.DefaultRoot )
  cyphernodeAdmin.routerGroups[globals.ROUTER_GROUPS_PUBLIC].GET( globals.PUBLIC_ENDPOINTS_LOGIN, handlers.DefaultLogin )
  cyphernodeAdmin.routerGroups[globals.ROUTER_GROUPS_PUBLIC].GET( globals.PUBLIC_ENDPOINTS_CALLBACK, handlers.DefaultCallback)
  cyphernodeAdmin.routerGroups[globals.ROUTER_GROUPS_PUBLIC].GET( globals.PUBLIC_ENDPOINTS_BYEBYE, handlers.DefaultByeBye )

}

func (cyphernodeAdmin *CyphernodeAdmin) initPrivateHandlers() {
  cyphernodeAdmin.routerGroups[globals.ROUTER_GROUPS_PRIVATE].GET("/logout", handlers.DefaultLogout )
  cyphernodeAdmin.routerGroups[globals.ROUTER_GROUPS_PRIVATE].GET( globals.PRIVATE_ENDPOINTS_HOME, handlers.DefaultHome )
}

func (cyphernodeAdmin *CyphernodeAdmin) initUsersHandlers() {
  cyphernodeAdmin.routerGroups[globals.ROUTER_GROUPS_USERS].GET("/", handlers.FindUsers)
  cyphernodeAdmin.routerGroups[globals.ROUTER_GROUPS_USERS].POST("/", handlers.CreateUser)
  cyphernodeAdmin.routerGroups[globals.ROUTER_GROUPS_USERS].GET("/:id", handlers.GetUser)
  cyphernodeAdmin.routerGroups[globals.ROUTER_GROUPS_USERS].PUT("/:id", handlers.UpdateUser )
  cyphernodeAdmin.routerGroups[globals.ROUTER_GROUPS_USERS].PATCH("/:id", handlers.PatchUser )
  cyphernodeAdmin.routerGroups[globals.ROUTER_GROUPS_USERS].DELETE("/:id", handlers.DeleteUser )
  cyphernodeAdmin.routerGroups[globals.ROUTER_GROUPS_USERS].POST("/:id/roles", handlers.UserAddRoles )
  cyphernodeAdmin.routerGroups[globals.ROUTER_GROUPS_USERS].DELETE("/:id/roles/:roleId", handlers.UserRemoveRole )
}

func (cyphernodeAdmin *CyphernodeAdmin) initAppsHandlers() {
  cyphernodeAdmin.routerGroups[globals.ROUTER_GROUPS_APPS].GET("/", handlers.FindApps)
  cyphernodeAdmin.routerGroups[globals.ROUTER_GROUPS_APPS].POST("/", handlers.CreateApp)
  cyphernodeAdmin.routerGroups[globals.ROUTER_GROUPS_APPS].GET("/:id", handlers.GetApp)
  cyphernodeAdmin.routerGroups[globals.ROUTER_GROUPS_APPS].PUT("/:id", handlers.UpdateApp )
  cyphernodeAdmin.routerGroups[globals.ROUTER_GROUPS_APPS].PATCH("/:id", handlers.PatchApp )
  cyphernodeAdmin.routerGroups[globals.ROUTER_GROUPS_APPS].DELETE("/:id", handlers.DeleteApp )
  cyphernodeAdmin.routerGroups[globals.ROUTER_GROUPS_APPS].POST("/:id/roles", handlers.AppAddRoles )
  cyphernodeAdmin.routerGroups[globals.ROUTER_GROUPS_APPS].DELETE("/:id/roles/:roleId", handlers.AppRemoveRole )
}

func (cyphernodeAdmin *CyphernodeAdmin) initHydraHandlers() {
  // TODO: csrf protection
  cyphernodeAdmin.routerGroups[globals.ROUTER_GROUPS_HYDRA].GET("/consent", handlers.HydraConsentGet)
  cyphernodeAdmin.routerGroups[globals.ROUTER_GROUPS_HYDRA].GET("/login", handlers.HydraLoginGet)
  cyphernodeAdmin.routerGroups[globals.ROUTER_GROUPS_HYDRA].POST("/consent", handlers.HydraConsentPost)
  cyphernodeAdmin.routerGroups[globals.ROUTER_GROUPS_HYDRA].POST("/login", handlers.HydraLoginPost)
}
