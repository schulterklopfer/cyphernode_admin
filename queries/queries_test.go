package queries_test

import (
  "github.com/schulterklopfer/cyphernode_admin/dataSource"
  "github.com/schulterklopfer/cyphernode_admin/models"
  "github.com/schulterklopfer/cyphernode_admin/queries"
  "os"
  "testing"
)

func TestModels(t *testing.T) {
  dbFile := "/tmp/tests.sqlite3"
  os.Remove(dbFile)
  dataSource.Init(dbFile)

  role := new(models.RoleModel)

  role.ID = 1
  role.Name = "role"
  role.Description = "role description"
  role.AutoAssign = false
  role.AppId = 99999

  dataSource.GetDB().Create(role)

  t.Run( "Create user", createUser )
  t.Run("Get user", getUser )
  t.Run("Find users", findUsers )

}

func createUser( t *testing.T ) {

  db := dataSource.GetDB()
  role := new(models.RoleModel)

  db.Take( &role, 1)

  roles := make( []*models.RoleModel, 1 )
  roles[0]=role

  var user *models.UserModel

  user = new(models.UserModel)
  user.Login = "login1"
  user.Password = "password1"
  user.Name ="name1"
  user.EmailAddress = "email1@email.rocks"
  user.Roles = roles

  user, err := queries.CreateUser( user )
  _, err = queries.CreateUser( user )

  if err == nil {
    t.Error( "Create user with user id" )
  }

  user.ID = 0
  _, err = queries.CreateUser( user )

  if err == nil {
    t.Error( "Created same user twice" )
  }



  user = new(models.UserModel)
  user.Login = "login2"
  user.Password = "password2"
  user.Name ="name2"
  user.EmailAddress = "email2@email.rocks"
  user.Roles = roles

  queries.CreateUser( user )

}

func getUser( t *testing.T ) {
  user, _ := queries.GetUser(1, true)

  if user == nil {
    t.Error( "unable to load user")
  }

}

func findUsers( t *testing.T ) {
  var users []*models.UserModel

  users, _ = queries.FindUsers( nil, "", -1,0, true)

  if users == nil {
    t.Error( "unable to load all users")
  }

  users, _ = queries.FindUsers( nil, "", 0,0,true)

  if users == nil || len(users) != 0 {
    t.Error( "unable to load 0 users")
  }

  users, _ = queries.FindUsers( nil, "", 1,0,true)

  if users == nil || len(users) != 1 || users[0].Name != "name1" || len(users[0].Roles) != 1 {
    t.Error( "unable to load 1 user")
  }

  users, _ = queries.FindUsers( nil, "", 1,1,true)

  if users == nil || len(users) != 1 || users[0].Name != "name2" {
    t.Error( "unable to load 1 user with offset 1")
  }

  users, _ = queries.FindUsers( nil, "id desc", 1,1,true)

  if users == nil || len(users) != 1 || users[0].Name != "name1" {
    t.Error( "unable to load 1 user with offset 1 order by id desc")
  }

  users, _ = queries.FindUsers( []interface{}{"name LIKE ?", "name%"}, "name", -1,0,true)

  if users == nil || len(users) != 2 || users[0].Name != "name1" || users[1].Name != "name2" {
    t.Error( "unable to load 2 users with order by name")
  }

}

