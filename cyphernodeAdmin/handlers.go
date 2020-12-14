package cyphernodeAdmin

import (
  "github.com/schulterklopfer/cyphernode_admin/globals"
  "github.com/schulterklopfer/cyphernode_admin/handlers"
)

func (cyphernodeAdmin *CyphernodeAdmin) initForwardAuthHandlers() {
  cyphernodeAdmin.engineForwardAuth.GET( globals.FORWARD_AUTH_ENDPOINTS_AUTH, handlers.ForwardAuthAuth )
}

func (cyphernodeAdmin *CyphernodeAdmin) initInternalHandlers() {
  cyphernodeAdmin.engineInternal.GET( globals.BASE_ENDPOINT_SESSIONS+"/:sessionID", handlers.GetSession )
  cyphernodeAdmin.engineInternal.PATCH( globals.BASE_ENDPOINT_SESSIONS+"/:sessionID", handlers.PatchSession )
  cyphernodeAdmin.engineInternal.DELETE( globals.BASE_ENDPOINT_SESSIONS+"/:sessionID", handlers.DeleteSession )
  cyphernodeAdmin.engineInternal.POST( globals.BASE_ENDPOINT_SESSIONS+"", handlers.CreateSession )
}

func (cyphernodeAdmin *CyphernodeAdmin) initPublicHandlers() {
  cyphernodeAdmin.engineExternal.POST( globals.PUBLIC_ENDPOINTS_LOGIN, handlers.DefaultLogin )
}

func (cyphernodeAdmin *CyphernodeAdmin) initPrivateHandlers() {
  cyphernodeAdmin.engineExternal.GET(globals.PRIVATE_ENDPOINTS_LOGOUT, handlers.DefaultLogout )
}

func (cyphernodeAdmin *CyphernodeAdmin) initUsersHandlers() {
  cyphernodeAdmin.engineExternal.GET(globals.BASE_ENDPOINT_USERS+"", handlers.FindUsers)
  cyphernodeAdmin.engineExternal.POST(globals.BASE_ENDPOINT_USERS+"", handlers.CreateUser)
  cyphernodeAdmin.engineExternal.GET(globals.BASE_ENDPOINT_USERS+"/:id", handlers.GetUser)
  cyphernodeAdmin.engineExternal.PATCH(globals.BASE_ENDPOINT_USERS+"/:id", handlers.PatchUser )
  cyphernodeAdmin.engineExternal.DELETE(globals.BASE_ENDPOINT_USERS+"/:id", handlers.DeleteUser )
  cyphernodeAdmin.engineExternal.POST(globals.BASE_ENDPOINT_USERS+"/:id/roles", handlers.UserAddRoles )
  cyphernodeAdmin.engineExternal.PATCH(globals.BASE_ENDPOINT_USERS+"/:id/roles", handlers.UserPatchRoles )
  cyphernodeAdmin.engineExternal.DELETE(globals.BASE_ENDPOINT_USERS+"/:id/roles/:roleId", handlers.UserRemoveRole )
}

func (cyphernodeAdmin *CyphernodeAdmin) initAppsHandlers() {
  cyphernodeAdmin.engineExternal.GET(globals.BASE_ENDPOINT_APPS+"", handlers.FindApps)
  cyphernodeAdmin.engineExternal.POST(globals.BASE_ENDPOINT_APPS+"", handlers.CreateApp)
  cyphernodeAdmin.engineExternal.GET(globals.BASE_ENDPOINT_APPS+"/:id", handlers.GetApp)
  cyphernodeAdmin.engineExternal.PUT(globals.BASE_ENDPOINT_APPS+"/:id", handlers.UpdateApp )
  cyphernodeAdmin.engineExternal.PATCH(globals.BASE_ENDPOINT_APPS+"/:id", handlers.PatchApp )
  cyphernodeAdmin.engineExternal.DELETE(globals.BASE_ENDPOINT_APPS+"/:id", handlers.DeleteApp )
  cyphernodeAdmin.engineExternal.POST(globals.BASE_ENDPOINT_APPS+"/:id/roles", handlers.AppAddRoles )
  cyphernodeAdmin.engineExternal.DELETE(globals.BASE_ENDPOINT_APPS+"/:id/roles/:roleId", handlers.AppRemoveRole )
}
