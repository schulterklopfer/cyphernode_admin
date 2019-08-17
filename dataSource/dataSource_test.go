package dataSource_test

import (
  "github.com/jinzhu/gorm"
  "github.com/schulterklopfer/cyphernode_admin/dataSource"
  "github.com/schulterklopfer/cyphernode_admin/dataSource/models"
  "os"
  "testing"
)

var db *gorm.DB

func TestModels(t *testing.T) {
  dbFile := "/tmp/tests.sqlite3"
  os.Remove(dbFile)
  dataSource.Init(dbFile )
  db = dataSource.GetDB()
  t.Run("Create app", createApp )
  t.Run("Load app", loadApp )
  t.Run("Load role", loadRole )
  t.Run("Create user", createUser )
  t.Run("Load user", loadUser )

}

func createApp( t *testing.T ) {

  app1 := new(models.AppModel)
  app1.Hash = "hash1"
  app1.Name = "app1"
  app1.Description = "description"

  role1 := new( models.RoleModel);
  role1.Name = "role1"
  role1.Description = "description"

  role2 := new( models.RoleModel);
  role2.Name = "role2"
  role2.Description = "description"

  roles1 := [2]*models.RoleModel{role1,role2}

  app1.AvailableRoles = roles1[:]

  db.Create(app1)

  if db.NewRecord(app1) || db.NewRecord(role1) || db.NewRecord( role2 ) {
    t.Error("Failed to insert app")
  }

}

func loadApp( t *testing.T ) {
  var app models.AppModel
  db.First(&app, 1)

  var roles []*models.RoleModel

  db.Model(&app).Association("AvailableRoles").Find(&roles)

  if app.Name != "app1" || app.ID != 1 || roles[0].ID != 1 || roles[1].ID != 2 {
    t.Error("Failed to load app")
  }
}

func loadRole( t *testing.T ) {
  var role models.RoleModel
  db.First(&role, 1)

  var app models.AppModel

  db.First( &app, role.AppId )

  if app.Name != "app1" || app.ID != 1  {
    t.Error("Failed to load role")
  }
}

func createUser(t *testing.T) {
  user := new(models.UserModel)
  user.Login = "login"
  user.Name = "Test user"
  user.EmailAddress = "user@email.com"

  var app models.AppModel

  db.First( &app, 1 )

  var roles []*models.RoleModel

  db.Model(&app).Association("AvailableRoles").Find(&roles)

  user.Roles = roles[0:1]

  db.Create(user)

  if db.NewRecord(user) {
    t.Error("Failed to insert user")
  }
}

func loadUser( t *testing.T ) {
  var user models.UserModel
  db.First(&user, 1)

  var roles []*models.RoleModel

  db.Model(&user).Association("Roles").Find(&roles)

  if user.Login != "login" || user.ID != 1 || roles[0].ID != 1 {
    t.Error("Failed to load user")
  }
}