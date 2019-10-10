package cyphernodeAdmin

import (
	"github.com/gin-gonic/gin"
	"github.com/schulterklopfer/cyphernode_admin/dataSource"
	"github.com/schulterklopfer/cyphernode_admin/globals"
	"github.com/schulterklopfer/cyphernode_admin/helpers"
	"github.com/schulterklopfer/cyphernode_admin/hydraAPI"
	"os"
)

const ADMIN_APP_NAME string = "Cyphernode Admin"
const ADMIN_APP_DESCRIPTION string = "Manage your cyphernode"
const ADMIN_APP_HASH string = "00000000000000000000000000000000"

const ADMIN_APP_ADMIN_ROLE_NAME string = "admin"
const ADMIN_APP_ADMIN_ROLE_DESCRIPTION string = "Main admin with god mode"

type Config struct {
  DatabaseFile string
  InitialAdminLogin string
  InitialAdminPassword string
  InitialAdminName string
  InitialAdminEmailAddress string
}

type CyphernodeAdmin struct {
  config       *Config
  engine       *gin.Engine
  routerGroups map[string]*gin.RouterGroup
}

func NewCyphernodeAdmin(config *Config) *CyphernodeAdmin {
  cyphernodeAdmin := new(CyphernodeAdmin)
  cyphernodeAdmin.config = config
  return cyphernodeAdmin
}

func (cyphernodeAdmin *CyphernodeAdmin) Init() {
	dataSource.Init(cyphernodeAdmin.config.DatabaseFile)
	hydraAPI.Init()
  cyphernodeAdmin.routerGroups = make(map[string]*gin.RouterGroup)
  cyphernodeAdmin.migrate()
  cyphernodeAdmin.engine = gin.Default()
  cyphernodeAdmin.engine.LoadHTMLGlob("templates/**/*.tmpl")
	cyphernodeAdmin.initUsersHandlers()
  cyphernodeAdmin.initAppsHandlers()
  cyphernodeAdmin.initRolesHandlers()
  cyphernodeAdmin.initHydraHandlers()
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
