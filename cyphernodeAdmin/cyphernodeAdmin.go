package cyphernodeAdmin

import (
  "github.com/gin-gonic/gin"
  "github.com/schulterklopfer/cyphernode_admin/appWhitelist"
  "github.com/schulterklopfer/cyphernode_admin/cnaErrors"
  "github.com/schulterklopfer/cyphernode_admin/cnaOIDC"
  "github.com/schulterklopfer/cyphernode_admin/dataSource"
  "github.com/schulterklopfer/cyphernode_admin/globals"
  "github.com/schulterklopfer/cyphernode_admin/helpers"
  "github.com/schulterklopfer/cyphernode_admin/hydraAPI"
  "github.com/schulterklopfer/cyphernode_admin/logwrapper"
  "github.com/schulterklopfer/cyphernode_admin/models"
  "github.com/schulterklopfer/cyphernode_admin/queries"
  "golang.org/x/sync/errgroup"
  "net/http"
  "os"
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
  config       *Config
  engineInternal *gin.Engine
  engineExternal *gin.Engine
  routerGroups map[string]*gin.RouterGroup
}

var instance *CyphernodeAdmin

func NewCyphernodeAdmin(config *Config) *CyphernodeAdmin {
  cyphernodeAdmin := new(CyphernodeAdmin)
  cyphernodeAdmin.config = config
  return cyphernodeAdmin
}

func (cyphernodeAdmin *CyphernodeAdmin) Init() error {

  dataSource.Init(cyphernodeAdmin.config.DatabaseFile)
  hydraAPI.Init()

  cyphernodeAdmin.routerGroups = make(map[string]*gin.RouterGroup)
  err := cyphernodeAdmin.migrate()
  if err != nil {
    return err
  }

  var thisApp models.AppModel
  err = queries.Get( &thisApp, 1, true )

  if err != nil {
    return err
  }

  if thisApp.ID == 0 {
    return cnaErrors.ErrNoSuchApp
  }

  cnaOIDC.Init( cnaOIDC.NewInitParams(
    thisApp.ClientID,
    thisApp.ClientSecret,
    helpers.AbsoluteURL(globals.URLS_CALLBACK),
    helpers.GetenvOrDefault( globals.OIDC_DISCOVERY_URL_ENV_KEY, globals.DEFAULTS_OIDC_DISCOVERY_URL ),
    helpers.AbsoluteURL( globals.ROUTER_GROUPS_BASE_ENDPOINT_SESSIONS ),
    []byte(helpers.GetenvOrDefault( globals.OIDC_SESSION_COOKIE_SECRET_ENV_KEY, globals.DEFAULTS_OIDC_SESSION_COOKIE_SECRET ) ),
    helpers.GetenvOrDefault( globals.OIDC_SSO_COOKIE_DOMAIN_ENV_KEY, globals.DEFAULTS_OIDC_SSO_COOKIE_DOMAIN )) )

  cyphernodeAdmin.engineInternal = gin.New()
  cyphernodeAdmin.engineExternal = gin.New()
  cyphernodeAdmin.engineExternal.LoadHTMLGlob("templates/**/*.tmpl")
  cyphernodeAdmin.createRouterGroups()
  // add session checks b4 other handlers so they are handled first
  // order is important here
  if !cyphernodeAdmin.config.DisableAuth {
    for i := 0; i < len(globals.PROTECTED_ROUTER_GROUPS_INDICES); i++ {
      cyphernodeAdmin.routerGroups[globals.ROUTER_GROUPS[globals.PROTECTED_ROUTER_GROUPS_INDICES[i]]].Use(CheckSession())
    }
  }
  // create handlers for public and private endpoints
  cyphernodeAdmin.initInternalHandlers()
  cyphernodeAdmin.initPublicHandlers()
  cyphernodeAdmin.initPrivateHandlers()
  cyphernodeAdmin.initSessionHandlers()
  cyphernodeAdmin.initUsersHandlers()
  cyphernodeAdmin.initAppsHandlers()
  cyphernodeAdmin.initHydraHandlers()
  appWhitelist.Init( helpers.GetenvOrDefault( globals.CNA_ADMIN_APP_WHITELIST_FILE_ENV_KEY, globals.DEFAULTS_CNA_ADMIN_APP_WHITELIST_FILE ) )

  return nil
}

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

func (cyphernodeAdmin *CyphernodeAdmin) Engine() *gin.Engine {
  return cyphernodeAdmin.engineExternal
}

func (cyphernodeAdmin *CyphernodeAdmin) Start() {
  if os.Getenv(globals.HYDRA_DISABLE_SYNC_ENV_KEY) == "" {
    helpers.SetInterval(cyphernodeAdmin.checkHydraClients, 1000, false)
  }

  var g errgroup.Group

  // internal interface, only available to cypherapps
  g.Go(func() error {
    return cyphernodeAdmin.engineInternal.Run(":3031")
  })

  // external interface behind treaefik
  g.Go(func() error {
    return  cyphernodeAdmin.engineExternal.Run(":3030")
  })

  if err := g.Wait(); err != nil {
    logwrapper.Logger().Fatal(err)
  }

}
