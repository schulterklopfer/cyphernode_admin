package handlers_test

import (
  "bytes"
  "encoding/json"
  "github.com/schulterklopfer/cyphernode_admin/cyphernodeAdmin"
  "github.com/schulterklopfer/cyphernode_admin/dataSource"
  "github.com/schulterklopfer/cyphernode_admin/logwrapper"
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

func TestMain(m *testing.M) {
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

  code := m.Run()

  testServer.Close()
  dataSource.Close()
  os.Remove(config.DatabaseFile)
  os.Exit(code)
}

func TestGetUser(t *testing.T) {

  res, err := testServer.Client().Get( testServer.URL+"/api/v0/users/1" )

  if err != nil {
    t.Error(err)
  }

  body, err := ioutil.ReadAll(res.Body)
  if err != nil {
    t.Error(err)
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

func TestCreateUser(t *testing.T) {

  input := map[string]string {
    "login": "testUser",
    "name": "name",
    "email_address": "email@test.com",
    "password": "test123",
  }

  body, err := json.Marshal( input )

  res, err := testServer.Client().Post( testServer.URL+"/api/v0/users/", "application/json", bytes.NewBuffer(body) )

  if err != nil {
    t.Error(err)
  }

  body, err = ioutil.ReadAll(res.Body)
  if err != nil {
    t.Error(err)
  }

  var user transforms.UserV0
  json.Unmarshal(body, &user)

  if user.Login != "testUser" ||
      user.Name != "name" ||
      user.EmailAddress != "email@test.com" ||
      user.ID != 2 ||
      len(user.Roles) != 0 {
    t.Error( "error in get user" )
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

func TestFindUser(t *testing.T) {

  res, err := testServer.Client().Get( testServer.URL+"/api/v0/users/?login_like=testUser" )

  if err != nil {
    t.Error(err)
  }

  body, err := ioutil.ReadAll(res.Body)
  if err != nil {
    t.Error(err)
  }

  var pagedResult pagedUsers
  json.Unmarshal(body, &pagedResult)

  var bodyStr = string(body)

  print( bodyStr )

  if pagedResult.Total != 2 ||
    pagedResult.Page != 0 ||
    len(pagedResult.Data) != 1 ||
    pagedResult.Data[0].Login != "testUser" ||
    pagedResult.Data[0].Name != "name" ||
    pagedResult.Data[0].EmailAddress != "email@test.com" ||
    pagedResult.Data[0].ID != 2 ||
    len(pagedResult.Data[0].Roles) != 0 {
    t.Error( "error in get user" )
  }
}


func TestDeleteUser(t *testing.T) {

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



