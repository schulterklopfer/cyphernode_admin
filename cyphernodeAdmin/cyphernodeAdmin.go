package cyphernodeAdmin

import (
  "encoding/json"
  "github.com/SatoshiPortal/cam/cyphernodeInfo"
  "github.com/SatoshiPortal/cam/utils"
  "github.com/gin-contrib/cors"
  "github.com/gin-gonic/gin"
  "github.com/schulterklopfer/cyphernode_admin/appList"
  "github.com/schulterklopfer/cyphernode_admin/cyphernodeApi"
  "github.com/schulterklopfer/cyphernode_admin/cyphernodeState"
  "github.com/schulterklopfer/cyphernode_admin/dataSource"
  "github.com/schulterklopfer/cyphernode_admin/dockerApi"
  "github.com/schulterklopfer/cyphernode_admin/globals"
  "github.com/schulterklopfer/cyphernode_admin/helpers"
  "github.com/schulterklopfer/cyphernode_admin/logwrapper"
  "golang.org/x/sync/errgroup"
  "io/ioutil"
  "strconv"
)

const ADMIN_APP_NAME string = "Cyphernode Admin"
const ADMIN_APP_DESCRIPTION string = "Manage your cyphernode"

const ADMIN_APP_ADMIN_ROLE_NAME string = "admin"
const ADMIN_APP_ADMIN_ROLE_DESCRIPTION string = "Main admin with god mode"

type Config struct {
  DatabaseFile string
  InitialAdminLogin string
  InitialAdminPassword string
  InitialAdminName string
  InitialAdminEmailAddress string
  DisableAuth bool
}

type CyphernodeAdmin struct {
  Config            *Config
  engineInternal    *gin.Engine
  engineExternal    *gin.Engine
  engineForwardAuth *gin.Engine
  routerGroups      map[string]*gin.RouterGroup
  ClientID          string
  Secret            string
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

  var cyphernodeInfo cyphernodeInfo.CyphernodeInfo

  cyphernodeInfoJsonBytes, err := ioutil.ReadFile( utils.GetCyphernodeInfoFilePath() )
  if err != nil {
    return err
  }

  err = json.Unmarshal( cyphernodeInfoJsonBytes, &cyphernodeInfo )
  if err != nil {
    return err
  }

  err = dataSource.Init(cyphernodeAdmin.Config.DatabaseFile)
  if err != nil {
    logwrapper.Logger().Error("Failed to create database" )
    return err
  }

  port, err := strconv.Atoi(helpers.GetenvOrDefault(globals.GATEKEEPER_PORT_ENV_KEY))

  if err != nil {
    logwrapper.Logger().Error("Failed to read gatekeeper port from env" )
    return err
  }

  err = dockerApi.Init()

  if err != nil {
    logwrapper.Logger().Error("Failed to connect to docker daemon" )
    return err
  }

  err = cyphernodeApi.Init( &cyphernodeApi.CyphernodeApiConfig{
    Version: cyphernodeInfo.ApiVersions[len(cyphernodeInfo.ApiVersions)-1],
    Host: helpers.GetenvOrDefault(globals.GATEKEEPER_HOST_ENV_KEY),
    Port: port,
  })

  if err != nil {
    logwrapper.Logger().Error("Failed to init cyphernode api" )
    return err
  }

  err = cyphernodeState.Init( &cyphernodeInfo )

  if err != nil {
    logwrapper.Logger().Error("Failed to init cyphernode state" )
    return err
  }

  cyphernodeAdmin.routerGroups = make(map[string]*gin.RouterGroup)
  err = cyphernodeAdmin.migrate()
  if err != nil {
    logwrapper.Logger().Error("Failed to init database" )
    return err
  }

  cyphernodeAdmin.engineForwardAuth = gin.New()
  cyphernodeAdmin.engineInternal = gin.New()
  cyphernodeAdmin.engineExternal = gin.New()

  cyphernodeAdmin.engineExternal.Use(cors.New(cors.Config{
    AllowMethods: []string{"POST", "GET", "OPTIONS", "PATCH", "DELETE"},
    AllowAllOrigins: true,
    AllowHeaders: []string{"Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization"},
  }))

  // add session checks b4 other handlers so they are handled first
  // order is important here
  /*
  if !cyphernodeAdmin.Config.DisableAuth {
    for i := 0; i < len(globals.PROTECTED_ROUTER_GROUPS_INDICES); i++ {
      cyphernodeAdmin.routerGroups[globals.ROUTER_GROUPS[globals.PROTECTED_ROUTER_GROUPS_INDICES[i]]].Use(CheckSession())
    }
  }
  */
  // create handlers for public and private endpoints
  cyphernodeAdmin.initInternalHandlers()
  cyphernodeAdmin.initPrivateHandlers()
  cyphernodeAdmin.initUsersHandlers()
  cyphernodeAdmin.initAppsHandlers()
  cyphernodeAdmin.initDockerHandlers()
  cyphernodeAdmin.initForwardAuthHandlers()
  cyphernodeAdmin.initPublicHandlers()
  err = appList.Init( helpers.GetenvOrDefault( globals.CYPHERAPPS_INSTALL_DIR_ENV_KEY ) )
  if err != nil {
    logwrapper.Logger().Error("Failed to init applist" )
    return err
  }
  return nil
}

/*
func CheckSession() gin.HandlerFunc {
  return func(c *gin.Context) {
    if !helpers.EndpointIsPublic( c.Request.URL.Path ) {
      // fetch userinfo from hydra
      if _, exists := c.Get("user"); !exists {
        user, err := cnaOIDC.GetUser(c.Writer, c.Request)
        if err != nil {
          c.Redirect(http.StatusTemporaryRedirect, globals.ROUTER_GROUPS_BASE_ENDPOINT_PUBLIC+globals.PUBLIC_ENDPOINTS_LOGIN)
          return
        }
        // put it in gin context for other handlers
        c.Set("user", user)
      }
    }
    c.Next()
  }
}
 */

func (cyphernodeAdmin *CyphernodeAdmin) Engine() *gin.Engine {
  return cyphernodeAdmin.engineExternal
}

func (cyphernodeAdmin *CyphernodeAdmin) Start() {

  var g errgroup.Group

  // oathkeeper session checker
  g.Go(func() error {
    return  cyphernodeAdmin.engineForwardAuth.Run(":3032")
  })

  // internal interface, only available to cypherapps
  g.Go(func() error {
    return cyphernodeAdmin.engineInternal.Run(":3031")
  })

  // external interface behind treaefik
  g.Go(func() error {
    return cyphernodeAdmin.engineExternal.Run(":3030")
  })

  if err := g.Wait(); err != nil {
    logwrapper.Logger().Fatal(err)
  }

}
