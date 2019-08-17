package app

import (
	"github.com/gin-gonic/gin"
	"github.com/schulterklopfer/cyphernode_admin/dataSource"
	"github.com/schulterklopfer/cyphernode_admin/handlers"
  "github.com/schulterklopfer/cyphernode_admin/models"
)

type Config struct {
	DatabaseFile         string
	InitialAdminUsername string
	InitialAdminPassword string
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
		app.routerGroups["users"].GET("/:id", handlers.GetUser)

		/*v1.POST("/", createTodo)
		  v1.GET("/", fetchAllTodo)
		  v1.GET("/:id", fetchSingleTodo)
		  v1.PUT("/:id", updateTodo)
		  v1.DELETE("/:id", deleteTodo)
		*/
	}
}

func (self *App) initContent() {
	// Create user id=1, app id=1, role id=1
  role := new(models.RoleModel)
  app := new(models.AppModel)
  user := new(models.UserModel)

  db := dataSource.GetDB()

  db.First( role, 1 )
  db.First( app, 1 )
  db.First( user, 1 )

  if !db.NewRecord( role ) || !db.NewRecord( app ) || !db.NewRecord( user ) {
    return
  }

  role.ID = 1
  role.Name = "admin"
  role.Description = "Main admin with god mode"
  role.AutoAssign = false
  role.AppId = 1

  roles := make( []*models.RoleModel, 1 )
  roles[0]=role

  app.ID = 1
  app.Name = "Cyphernode admin"
  app.Description = "Manage your cyphernode"
  app.Hash = "adminapphash" // change me
  app.AvailableRoles = roles

  user.ID = 1
  user.Login = self.config.InitialAdminUsername
  user.Password = self.config.InitialAdminPassword
  user.Name = "Adnim Administer"
  user.EmailAddress = "admin@admin.rocks"
  user.Roles = roles


  db.Create( role )
  db.Create( app )
  db.Create( user )

}

func (app *App) Start() {
	app.engine.Run()
}
