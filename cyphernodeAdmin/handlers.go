/*
 * MIT License
 *
 * Copyright (c) 2021 schulterklopfer/__escapee__
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILIT * Y, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package cyphernodeAdmin

import (
  "github.com/gin-gonic/gin"
  "github.com/schulterklopfer/cyphernode_admin/globals"
  "github.com/schulterklopfer/cyphernode_admin/handlers"
  "github.com/schulterklopfer/cyphernode_admin/helpers"
)

func (cyphernodeAdmin *CyphernodeAdmin) initAuthHandlers() {
  cyphernodeAdmin.engineAuth.GET( globals.FORWARD_AUTH_ENDPOINTS_AUTH, handlers.ForwardAuthAuth )
  cyphernodeAdmin.engineAuth.GET( globals.PROXY_GATEKEEPER_ENDPOINTS_AUTH, handlers.ProxyGatekeeperAuth )
}

func (cyphernodeAdmin *CyphernodeAdmin) initInternalHandlers() {
  cyphernodeAdmin.engineInternal.GET( globals.BASE_ENDPOINT_SESSIONS+"/:sessionID", handlers.GetSession )
  cyphernodeAdmin.engineInternal.PATCH( globals.BASE_ENDPOINT_SESSIONS+"/:sessionID", handlers.PatchSession )
  cyphernodeAdmin.engineInternal.DELETE( globals.BASE_ENDPOINT_SESSIONS+"/:sessionID", handlers.DeleteSession )
  cyphernodeAdmin.engineInternal.POST( globals.BASE_ENDPOINT_SESSIONS+"", handlers.CreateSession )
}

func (cyphernodeAdmin *CyphernodeAdmin) initPublicHandlers() {
  cyphernodeAdmin.engineExternal.POST( globals.PUBLIC_ENDPOINTS_LOGIN, handlers.DefaultLogin )
  cyphernodeAdmin.engineExternal.Static( "/_", helpers.GetenvOrDefault(globals.CNA_STATIC_FILE_DIR_ENV_KEY) )
  cyphernodeAdmin.engineExternal.GET(globals.BASE_ENDPOINT_STATUS+"", handlers.GetStatus )
  cyphernodeAdmin.engineExternal.Use(func(c *gin.Context) {
    if c.Request.URL.Path == "/" || c.Request.URL.Path == "/_/" {
      c.Redirect( 307, "/"+globals.BASE_ADMIN_MOUNTPOINT+"/_/index.html" )
      return
    }
    c.Next()
  })
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

func (cyphernodeAdmin *CyphernodeAdmin) initDockerHandlers() {
  cyphernodeAdmin.engineExternal.GET(globals.BASE_ENDPOINT_DOCKER+"/image/:image", handlers.GetContainerByBas64Image)
  cyphernodeAdmin.engineExternal.GET(globals.BASE_ENDPOINT_DOCKER+"/name/:name", handlers.GetContainerByName)
  cyphernodeAdmin.engineExternal.GET(globals.BASE_ENDPOINT_DOCKER+"/logs/:containerId", handlers.WSLogsByContainerId)
}

func (cyphernodeAdmin *CyphernodeAdmin) initBlocksHandlers() {
  cyphernodeAdmin.engineExternal.GET(globals.BASE_ENDPOINT_BLOCKS+"/latest", handlers.GetLatestBlocks )
}