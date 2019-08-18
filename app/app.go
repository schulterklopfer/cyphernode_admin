package app

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

type App struct {
  config       *Config
  engine       *gin.Engine
  routerGroups map[string]*gin.RouterGroup
}

func NewApp(config *Config) *App {
  app := new(App)
  app.config = config
  return app
}

func (app *App) Init() {
  app.routerGroups = make(map[string]*gin.RouterGroup)
  dataSource.Init(app.config.DatabaseFile)
  app.initContent()
  app.engine = gin.Default()
  app.routerGroups["users"] = app.engine.Group("/api/v0/users")
  {
    app.routerGroups["users"].GET("/", handlers.FindUsers)
    app.routerGroups["users"].GET("/:id", handlers.GetUser)

    /*v1.POST("/", createTodo)
      v1.GET("/", fetchAllTodo)
      v1.GET("/:id", fetchSingleTodo)
      v1.PUT("/:id", updateTodo)
      v1.DELETE("/:id", deleteTodo)
    */
  }
}

func (app *App) initContent() error {

  // Create adminUser id=1, app id=1, adminRole id=1
  adminRole := new(models.RoleModel)
  adminApp := new(models.AppModel)
  adminUser := new(models.UserModel)

  db := dataSource.GetDB()

  db.Take(adminRole, 1 )
  db.Take(adminApp, 1 )
  db.Take(adminUser, 1 )

  if !db.NewRecord(adminRole) || !db.NewRecord( app ) || !db.NewRecord(adminUser) {
    return nil
  }

  hashedPassword, err := password.HashPassword( app.config.InitialAdminPassword )

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
  adminApp.Name = "Cyphernode admin app"
  adminApp.Description = "Manage your cyphernode"
  adminApp.Hash = "adminapphash" // change me
  adminApp.AvailableRoles = roles

  adminUser.ID = 1
  adminUser.Login = app.config.InitialAdminLogin
  adminUser.Password = hashedPassword
  adminUser.Name = app.config.InitialAdminName
  adminUser.EmailAddress = app.config.InitialAdminEmailAddress
  adminUser.Roles = roles

  tx := db.Begin()
  tx.Create(adminRole)
  tx.Create(adminApp)
  tx.Create(adminUser)
  return tx.Commit().Error

}

func (app *App) Start() {
  app.engine.Run()
}
