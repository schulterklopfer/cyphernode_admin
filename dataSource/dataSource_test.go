package dataSource_test

import (
  "github.com/schulterklopfer/cyphernode_admin/dataSource"
  "github.com/schulterklopfer/cyphernode_admin/logwrapper"
  "github.com/schulterklopfer/cyphernode_admin/models"
  "github.com/sirupsen/logrus"
  "math/rand"
  "os"
  "strconv"
  "testing"
  "time"
)

func TestDataSource(t *testing.T) {
  r := rand.New(rand.NewSource(time.Now().UnixNano()))
  logwrapper.Logger().SetLevel( logrus.PanicLevel )
  dbFile := "/tmp/tests_"+strconv.Itoa(r.Intn(1000000 ))+".sqlite3"
  dataSource.Init(dbFile)

  t.Run("testCreateApp", testCreateApp )
  t.Run("testLoadApp", testLoadApp )
  t.Run("testLoadRole", testLoadRole )
  t.Run("testCreateUser", testCreateUser )
  t.Run("testLoadUser", testLoadUser )

  dataSource.Close()

  os.Remove(dbFile)
}

func testCreateApp(t *testing.T) {

	app1 := new(models.AppModel)
	app1.Hash = "hash1"
	app1.Name = "app1"
	app1.Description = "description"

	role1 := new(models.RoleModel)
	role1.Name = "role1"
	role1.Description = "description"

	role2 := new(models.RoleModel)
	role2.Name = "role2"
	role2.Description = "description"

	roles1 := [2]*models.RoleModel{role1, role2}

	app1.AvailableRoles = roles1[:]

	db := dataSource.GetDB()
	db.Create(app1)

	if db.NewRecord(app1) || db.NewRecord(role1) || db.NewRecord(role2) {
		t.Error("Failed to insert app")
	}

}

func testLoadApp(t *testing.T) {
	var app models.AppModel
  db := dataSource.GetDB()
  db.First(&app, 1)

	var roles []*models.RoleModel

	db.Model(&app).Association("AvailableRoles").Find(&roles)

	if app.Name != "app1" || app.ID != 1 || roles[0].ID != 1 || roles[1].ID != 2 {
		t.Error("Failed to load app")
	}
}

func testLoadRole(t *testing.T) {
	var role models.RoleModel
  db := dataSource.GetDB()
  db.First(&role, 1)

	var app models.AppModel

	db.First(&app, role.AppId)

	if app.Name != "app1" || app.ID != 1 {
		t.Error("Failed to load role")
	}
}

func testCreateUser(t *testing.T) {
	user := new(models.UserModel)
	user.Login = "login"
	user.Name = "Test user"
	user.EmailAddress = "user@email.com"

	var app models.AppModel

  db := dataSource.GetDB()
	db.First(&app, 1)

	var roles []*models.RoleModel

	db.Model(&app).Association("AvailableRoles").Find(&roles)

	user.Roles = roles[0:1]

	db.Create(user)

	if db.NewRecord(user) {
		t.Error("Failed to insert user")
	}
}

func testLoadUser(t *testing.T) {
	var user models.UserModel

	db := dataSource.GetDB()
  db.First(&user, 1)

	var roles []*models.RoleModel

	db.Model(&user).Association("Roles").Find(&roles)

	if user.Login != "login" || user.ID != 1 || roles[0].ID != 1 {
		t.Error("Failed to load user")
	}
}
