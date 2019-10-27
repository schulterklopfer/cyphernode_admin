package cyphernodeAdmin

import (
  "github.com/gin-gonic/gin"
  "github.com/schulterklopfer/cyphernode_admin/cnaErrors"
  "github.com/schulterklopfer/cyphernode_admin/cnaOIDC"
  "github.com/schulterklopfer/cyphernode_admin/dataSource"
  "github.com/schulterklopfer/cyphernode_admin/globals"
  "github.com/schulterklopfer/cyphernode_admin/helpers"
  "github.com/schulterklopfer/cyphernode_admin/hydraAPI"
  "github.com/schulterklopfer/cyphernode_admin/models"
  "github.com/schulterklopfer/cyphernode_admin/queries"
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
  engine       *gin.Engine
  routerGroups map[string]*gin.RouterGroup
}

var instance *CyphernodeAdmin

func NewCyphernodeAdmin(config *Config) *CyphernodeAdmin {
  cyphernodeAdmin := new(CyphernodeAdmin)
  cyphernodeAdmin.config = config
  return cyphernodeAdmin
}

func (cyphernodeAdmin *CyphernodeAdmin) Init() error {
  //sessionStore := sqliteStore.NewSqliteStore( []byte("secret") )

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
    os.Getenv( globals.OIDC_DISCOVERY_URL_ENV_KEY ),
    helpers.AbsoluteURL( globals.ROUTER_GROUPS_BASE_ENDPOINT_SESSIONS ),
    []byte(os.Getenv( globals.OIDC_SESSION_COOKIE_SECRET_ENV_KEY) ),
    "127.0.0.1") )

  cyphernodeAdmin.engine = gin.Default()
  cyphernodeAdmin.engine.LoadHTMLGlob("templates/**/*.tmpl")
  cyphernodeAdmin.createRouterGroups()

  // add session checks b4 other handlers so they are handled first
  // order is important here
  if !cyphernodeAdmin.config.DisableAuth {
    for i:=0; i<len( globals.PROTECTED_ROUTER_GROUPS_INDICES); i++ {
      cyphernodeAdmin.routerGroups[globals.ROUTER_GROUPS[globals.PROTECTED_ROUTER_GROUPS_INDICES[i]]].Use( CheckSession() )
    }
  }

  // create handlers for public and private endpoints
  cyphernodeAdmin.initPublicHandlers()
  cyphernodeAdmin.initPrivateHandlers()
  cyphernodeAdmin.initSessionHandlers()
  cyphernodeAdmin.initUsersHandlers()
  cyphernodeAdmin.initAppsHandlers()
  cyphernodeAdmin.initHydraHandlers()

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
  return cyphernodeAdmin.engine
}

func (cyphernodeAdmin *CyphernodeAdmin) Start() {
  if os.Getenv(globals.HYDRA_DISABLE_SYNC_ENV_KEY) == "" {
    helpers.SetInterval(cyphernodeAdmin.checkHydraClients, 1000, false)
  }
  cyphernodeAdmin.engine.Run("0.0.0.0:3030")
}
