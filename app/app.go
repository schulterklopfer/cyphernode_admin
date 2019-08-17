package app

import (
	"github.com/gin-gonic/gin"
	"github.com/schulterklopfer/cyphernode_admin/dataSource"
	"github.com/schulterklopfer/cyphernode_admin/handlers"
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

func (app *App) initContent() {
	// Create user id=1, app id=1, role id=1

}

func (app *App) Start() {
	app.engine.Run()
}
