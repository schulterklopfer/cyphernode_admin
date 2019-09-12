package cyphernodeAdmin

import (
  "github.com/gin-gonic/gin"
  "github.com/schulterklopfer/cyphernode_admin/dataSource"
  "github.com/schulterklopfer/cyphernode_admin/handlers"
  "github.com/schulterklopfer/cyphernode_admin/models"
  "github.com/schulterklopfer/cyphernode_admin/password"
)

const ADMIN_APP_NAME string = "Cyphernode admin cyphernodeAdmin"
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
  cyphernodeAdmin.routerGroups = make(map[string]*gin.RouterGroup)
  dataSource.Init(cyphernodeAdmin.config.DatabaseFile)
  cyphernodeAdmin.initContent()
  cyphernodeAdmin.engine = gin.Default()
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

func (cyphernodeAdmin *CyphernodeAdmin) initContent() error {

  // Create adminUser id=1, cyphernodeAdmin id=1, adminRole id=1
  adminRole := new(models.RoleModel)
  adminApp := new(models.AppModel)
  adminUser := new(models.UserModel)

  db := dataSource.GetDB()

  db.Take(adminRole, 1 )
  db.Take(adminApp, 1 )
  db.Take(adminUser, 1 )

  if !db.NewRecord(adminRole) || !db.NewRecord(cyphernodeAdmin) || !db.NewRecord(adminUser) {
    return nil
  }

  hashedPassword, err := password.HashPassword( cyphernodeAdmin.config.InitialAdminPassword )

  if err != nil {
    return err
  }

  adminRole.ID = 1
  adminRole.Name = ADMIN_APP_ADMIN_ROLE_NAME
  adminRole.Description = ADMIN_APP_ADMIN_ROLE_DESCRIPTION
  adminRole.AutoAssign = false
  adminRole.AppId = 1

  roles := make( []*models.RoleModel, 1 )
  roles[0]= adminRole

  adminApp.ID = 1
  adminApp.Name = ADMIN_APP_NAME
  adminApp.Description = ADMIN_APP_DESCRIPTION
  adminApp.Hash = ADMIN_APP_HASH
  adminApp.AvailableRoles = roles

  adminUser.ID = 1
  adminUser.Login = cyphernodeAdmin.config.InitialAdminLogin
  adminUser.Password = hashedPassword
  adminUser.Name = cyphernodeAdmin.config.InitialAdminName
  adminUser.EmailAddress = cyphernodeAdmin.config.InitialAdminEmailAddress
  adminUser.Roles = roles

  tx := db.Begin()
  tx.Create(adminRole)
  tx.Create(adminApp)
  tx.Create(adminUser)
  return tx.Commit().Error

}

func (cyphernodeAdmin *CyphernodeAdmin) Engine() *gin.Engine {
  return cyphernodeAdmin.engine
}

func (cyphernodeAdmin *CyphernodeAdmin) Start() {
  cyphernodeAdmin.engine.Run("localhost:8080")
}
