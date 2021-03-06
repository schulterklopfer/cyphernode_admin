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
  "encoding/json"
  "github.com/SatoshiPortal/cam/cyphernodeInfo"
  "github.com/SatoshiPortal/cam/utils"
  "github.com/gin-contrib/cors"
  "github.com/gin-gonic/gin"
  "github.com/schulterklopfer/cyphernode_admin/cyphernodeApi"
  "github.com/schulterklopfer/cyphernode_admin/cyphernodeKeys"
  "github.com/schulterklopfer/cyphernode_admin/cyphernodeState"
  "github.com/schulterklopfer/cyphernode_admin/dockerApi"
  "github.com/schulterklopfer/cyphernode_admin/globals"
  "github.com/schulterklopfer/cyphernode_fauth/dataSource"
  "github.com/schulterklopfer/cyphernode_fauth/helpers"
  "github.com/schulterklopfer/cyphernode_fauth/logwrapper"
  "golang.org/x/sync/errgroup"
  "io/ioutil"
  "strconv"
)

const ADMIN_APP_NAME string = "Cyphernode Admin"
const ADMIN_APP_DESCRIPTION string = "Manage your cyphernode"

const ADMIN_APP_ADMIN_ROLE_NAME string = "admin"
const ADMIN_APP_ADMIN_ROLE_DESCRIPTION string = "Main admin with god mode"

const ADMIN_APP_USER_ROLE_NAME string = "user"
const ADMIN_APP_USER_ROLE_DESCRIPTION string = "Regular user"


type Config struct {
  DatabaseDsn string
  InitialAdminLogin string
  InitialAdminPassword string
  InitialAdminName string
  InitialAdminEmailAddress string
  DisableAuth bool
}

type CyphernodeAdmin struct {
  Config         *Config
  //engineInternal *gin.Engine
  engineExternal *gin.Engine
  routerGroups   map[string]*gin.RouterGroup
  ClientID       string
  Secret         string
}

var instance *CyphernodeAdmin

func NewCyphernodeAdmin(config *Config) *CyphernodeAdmin {
  instance = new(CyphernodeAdmin)
  instance.Config = config
  return instance
}

func Get() *CyphernodeAdmin {
  return instance
}

func (cyphernodeAdmin *CyphernodeAdmin) Init() error {

  err := dataSource.Init(cyphernodeAdmin.Config.DatabaseDsn)
  if err != nil {
    logwrapper.Logger().Error("Failed to create database" )
    return err
  }

  err = cyphernodeKeys.Init(
    helpers.GetenvOrDefault( globals.KEYS_FILE_ENV_KEY ),
    helpers.GetenvOrDefault( globals.ACTIONS_FILE_ENV_KEY ),
  )

  if err != nil {
    logwrapper.Logger().Error("Failed to load cyphernode keys and api info" )
    return err
  }

  cyphernodeAdmin.routerGroups = make(map[string]*gin.RouterGroup)

  //cyphernodeAdmin.engineInternal = gin.New()
  cyphernodeAdmin.engineExternal = gin.New()

  cyphernodeAdmin.engineExternal.Use(cors.New(cors.Config{
    AllowMethods: []string{"POST", "GET", "OPTIONS", "PATCH", "DELETE"},
    AllowAllOrigins: true,
    AllowHeaders: []string{"Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization"},
  }))

  // create handlers for public and private endpoints
  cyphernodeAdmin.initDockerHandlers()
  cyphernodeAdmin.initBlocksHandlers()
  cyphernodeAdmin.initPublicHandlers()

  return nil
}

func (cyphernodeAdmin *CyphernodeAdmin) Engine() *gin.Engine {
  return cyphernodeAdmin.engineExternal
}

func (cyphernodeAdmin *CyphernodeAdmin) Start() {

  var g errgroup.Group

  // external interface behind treaefik
  g.Go(func() error {
    return cyphernodeAdmin.engineExternal.Run(":3030")
  })

  var cyphernodeInfo cyphernodeInfo.CyphernodeInfo

  cyphernodeInfoJsonBytes, err := ioutil.ReadFile( utils.GetCyphernodeInfoFilePath() )
  if err != nil {
    logwrapper.Logger().Fatal(err)
    return
  }

  err = json.Unmarshal( cyphernodeInfoJsonBytes, &cyphernodeInfo )
  if err != nil {
    logwrapper.Logger().Fatal(err)
    return
  }

  port, err := strconv.Atoi(helpers.GetenvOrDefault(globals.GATEKEEPER_PORT_ENV_KEY))

  if err != nil {
    logwrapper.Logger().Fatal(err)
    return
  }

  err = dockerApi.Init()

  if err != nil {
    logwrapper.Logger().Fatal(err)
    return
  }

  err = cyphernodeApi.Init( &cyphernodeApi.CyphernodeApiConfig{
    Version: cyphernodeInfo.ApiVersions[len(cyphernodeInfo.ApiVersions)-1],
    Host: helpers.GetenvOrDefault(globals.GATEKEEPER_HOST_ENV_KEY),
    Port: port,
  })

  if err != nil {
    logwrapper.Logger().Fatal(err)
  }

  err = cyphernodeState.Init( &cyphernodeInfo )

  if err != nil {
    logwrapper.Logger().Fatal(err)
    return
  }


  if err := g.Wait(); err != nil {
    logwrapper.Logger().Fatal(err)
  }

}
