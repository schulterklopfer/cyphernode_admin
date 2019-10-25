package cyphernodeAdmin

import (
  "github.com/gin-contrib/sessions"
  "github.com/gin-gonic/gin"
  "github.com/markbates/goth"
  "github.com/markbates/goth/providers/openidConnect"
  "github.com/schulterklopfer/cyphernode_admin/cnaErrors"
  "github.com/schulterklopfer/cyphernode_admin/dataSource"
  "github.com/schulterklopfer/cyphernode_admin/globals"
  "github.com/schulterklopfer/cyphernode_admin/helpers"
  "github.com/schulterklopfer/cyphernode_admin/hydraAPI"
  "github.com/schulterklopfer/cyphernode_admin/models"
  "github.com/schulterklopfer/cyphernode_admin/oidc"
  "github.com/schulterklopfer/cyphernode_admin/queries"
  "github.com/schulterklopfer/cyphernode_admin/sqliteStore"
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
  sqliteStore := sqliteStore.NewSqliteStore( []byte("secret") )

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

  openidConnect, _ := openidConnect.New(thisApp.ClientID, thisApp.ClientSecret, globals.URLS_CALLBACK, globals.URLS_OIDC_DISCOVERY )
  if openidConnect != nil {
    goth.UseProviders(openidConnect)
  }
  oidc.Store = sqliteStore

  cyphernodeAdmin.engine = gin.Default()
  cyphernodeAdmin.engine.LoadHTMLGlob("templates/**/*.tmpl")
  if !cyphernodeAdmin.config.DisableAuth {
    cyphernodeAdmin.engine.Use(sessions.Sessions("_oidc_session", sqliteStore))
    cyphernodeAdmin.engine.Use(CheckSession())
  }
  cyphernodeAdmin.initDefaultHandlers()
  cyphernodeAdmin.initSessionHandlers()
  cyphernodeAdmin.initUsersHandlers()
  cyphernodeAdmin.initAppsHandlers()
  cyphernodeAdmin.initRolesHandlers()
  cyphernodeAdmin.initHydraHandlers()

  return nil
}

func CheckSession() gin.HandlerFunc {
  return func(c *gin.Context) {
    if !helpers.EndpointIsPublic( c.Request.URL.Path ) {
      // fetch userinfo from hydra
      user, err := oidc.GetUser( c.Writer, c.Request )
      if err != nil {
        c.Redirect( http.StatusTemporaryRedirect, globals.ENDPOINTS_LOGIN)
        return
      }
      // all not public endpoints need the role "admin"
      userIsAdmin := helpers.UserIsAdmin( &user )
      print( userIsAdmin )
      c.Set("user", user )
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
