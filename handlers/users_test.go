package handlers_test

import (
  "bytes"
  "encoding/json"
  "github.com/schulterklopfer/cyphernode_admin/cyphernodeAdmin"
  "github.com/schulterklopfer/cyphernode_admin/dataSource"
  "github.com/schulterklopfer/cyphernode_admin/logwrapper"
  "github.com/schulterklopfer/cyphernode_admin/models"
  "github.com/schulterklopfer/cyphernode_admin/queries"
  "github.com/schulterklopfer/cyphernode_admin/transforms"
  "github.com/sirupsen/logrus"
  "io/ioutil"
  "math/rand"
  "net/http"
  "net/http/httptest"
  "os"
  "strconv"
  "testing"
  "time"
)

var testServer *httptest.Server

func TestUserHandlers( t *testing.T ) {
  r := rand.New(rand.NewSource(time.Now().UnixNano()))
  logwrapper.Logger().SetLevel( logrus.PanicLevel )

  var config cyphernodeAdmin.Config
  config.DatabaseFile = "/tmp/tests_"+strconv.Itoa(r.Intn(1000000 ))+".sqlite3"
  config.InitialAdminEmailAddress = "email@email.com"
  config.InitialAdminName = "admin"
  config.InitialAdminLogin = "admin"
  config.InitialAdminPassword = "test123"

  cnAdmin := cyphernodeAdmin.NewCyphernodeAdmin( &config )
  cnAdmin.Init()

  testServer = httptest.NewServer( cnAdmin.Engine() )

  var role models.RoleModel

  role.Name = "testRole"
  role.AutoAssign = false
  role.AppId = 99999

  _ = queries.CreateRole(&role)

  t.Run( "testGetUser", testGetUser )
  t.Run( "testCreateUser", testCreateUser )
  t.Run( "testUserAddRole", testUserAddRole )
  t.Run( "testUserRemoveRole", testUserRemoveRole )
  t.Run( "testPatchUser", testPatchUser )
  t.Run( "testFindUser", testFindUser )
  t.Run( "testDeleteUser", testDeleteUser )

  testServer.Close()
  dataSource.Close()
  os.Remove(config.DatabaseFile)
}

func testGetUser(t *testing.T) {

  res, err := testServer.Client().Get( testServer.URL+"/api/v0/users/1" )

  if err != nil {
    t.Error(err)
  }

  body, err := ioutil.ReadAll(res.Body)
  if err != nil {
    t.Error(err)
  }

  if res.StatusCode != http.StatusOK {
    t.Error("wrong status")
  }

  var user transforms.UserV0
  json.Unmarshal(body, &user)

  if user.Login != "admin" ||
    user.Name != "admin" ||
    user.EmailAddress != "email@email.com" ||
    user.ID != 1 ||
    len(user.Roles) != 1 {
    t.Error( "error in get user" )
  }
}

func testCreateUser(t *testing.T) {

  jsonInput := `
{
  "login": "testUser",
  "name": "name",
  "email_address": "email@test.com",
  "password": "test123",
  "roles": [
    { "ID": 1 }
  ]
}
`
  res, err := testServer.Client().Post( testServer.URL+"/api/v0/users/", "application/json", bytes.NewBuffer([]byte(jsonInput)) )

  if err != nil {
    t.Error(err)
  }

  body, err := ioutil.ReadAll(res.Body)
  if err != nil {
    t.Error(err)
  }

  if res.StatusCode != http.StatusCreated {
    t.Error("wrong status")
  }

  var user transforms.UserV0
  json.Unmarshal(body, &user)

  if user.Login != "testUser" ||
      user.Name != "name" ||
      user.EmailAddress != "email@test.com" ||
      user.ID != 2 ||
      len(user.Roles) != 1 {
    t.Error( "error in get user" )
  }
}

func testUserAddRole( t *testing.T ) {
  jsonInput := `{ "ID": 2 }`
  res, err := testServer.Client().Post( testServer.URL+"/api/v0/users/2/roles", "application/json", bytes.NewBuffer([]byte(jsonInput)) )

  if err != nil {
    t.Error(err)
  }

  body, err := ioutil.ReadAll(res.Body)
  if err != nil {
    t.Error(err)
  }

  if res.StatusCode != http.StatusOK {
    t.Error("wrong status")
  }

  var user transforms.UserV0
  json.Unmarshal(body, &user)

  if user.Login != "testUser" ||
      user.Name != "name" ||
      user.EmailAddress != "email@test.com" ||
      user.ID != 2 ||
      len(user.Roles) != 2 {
    t.Error( "error in user add role" )
  }
}

func testUserRemoveRole( t *testing.T ) {
  req, err := http.NewRequest("DELETE", testServer.URL+"/api/v0/users/2/roles/2",nil)

  res, err := testServer.Client().Do(req)

  if err != nil {
    t.Error(err)
  }

  if res.StatusCode != http.StatusNoContent {
    t.Error( "could not delete")
  }
}

func testPatchUser(t *testing.T) {

  jsonInput := `
{
  "login": "testUser2",
  "roles": []
}
`
  req, err := http.NewRequest("PATCH", testServer.URL+"/api/v0/users/2",bytes.NewBuffer([]byte(jsonInput)))
  req.Header.Set("Content-Type", "application/json")

  res, err := testServer.Client().Do(req)

  if err != nil {
    t.Error(err)
  }

  if res.StatusCode != http.StatusOK {
    t.Error("wrong status")
  }

  body, err := ioutil.ReadAll(res.Body)
  if err != nil {
    t.Error(err)
  }

  var user transforms.UserV0
  json.Unmarshal(body, &user)

  if user.Login != "testUser2" ||
      user.Name != "name" ||
      user.EmailAddress != "email@test.com" ||
      user.ID != 2 ||
      len(user.Roles) != 1 {
    t.Error( "error in patch user" )
  }
}

type pagedUsers struct {
  Page int `json:"page"`
  Limit int` json:"limit"`
  Sort string `json:"sort"`
  Order string `json:"order"`
  Total int `json:"total"`
  Data []*transforms.UserV0 `json:"data"`
}

func testFindUser(t *testing.T) {

  res, err := testServer.Client().Get( testServer.URL+"/api/v0/users/?login_like=testUser" )

  if err != nil {
    t.Error(err)
  }

  if res.StatusCode != http.StatusOK {
    t.Error("wrong status")
  }

  body, err := ioutil.ReadAll(res.Body)
  if err != nil {
    t.Error(err)
  }

  var pagedResult pagedUsers
  json.Unmarshal(body, &pagedResult)

  if pagedResult.Total != 2 ||
    pagedResult.Page != 0 ||
    len(pagedResult.Data) != 1 ||
    pagedResult.Data[0].Login != "testUser2" ||
    pagedResult.Data[0].Name != "name" ||
    pagedResult.Data[0].EmailAddress != "email@test.com" ||
    pagedResult.Data[0].ID != 2 ||
    len(pagedResult.Data[0].Roles) != 1 {
    t.Error( "error in find user" )
  }
}


func testDeleteUser(t *testing.T) {

  req, err := http.NewRequest("DELETE", testServer.URL+"/api/v0/users/2",nil)
  //req.Header.Set("Content-Type", "application/json")

  res, err := testServer.Client().Do(req)

  if err != nil {
    t.Error(err)
  }

  if res.StatusCode != http.StatusNoContent {
    t.Error( "could not delete")
  }
}



