package cyphernodeAdmin

import (
  "github.com/gin-gonic/gin"
  "github.com/schulterklopfer/cyphernode_admin/dataSource"
  "github.com/schulterklopfer/cyphernode_admin/handlers"
  "github.com/schulterklopfer/cyphernode_admin/models"
  "github.com/schulterklopfer/cyphernode_admin/password"
)

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


    /*v1.POST("/", createTodo)
      v1.GET("/", fetchAllTodo)
      v1.GET("/:id", fetchSingleTodo)
      v1.PUT("/:id", updateTodo)
      v1.DELETE("/:id", deleteTodo)
    */
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
  adminRole.Name = "admin"
  adminRole.Description = "Main admin with god mode"
  adminRole.AutoAssign = false
  adminRole.AppId = 1

  roles := make( []*models.RoleModel, 1 )
  roles[0]= adminRole

  adminApp.ID = 1
  adminApp.Name = "Cyphernode admin cyphernodeAdmin"
  adminApp.Description = "Manage your cyphernode"
  adminApp.Hash = "adminapphash" // change me
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
  cyphernodeAdmin.engine.Run()
}
