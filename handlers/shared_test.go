package handlers_test

import (
  "github.com/schulterklopfer/cyphernode_admin/cyphernodeAdmin"
  "github.com/schulterklopfer/cyphernode_admin/dataSource"
  "github.com/schulterklopfer/cyphernode_admin/logwrapper"
  "github.com/schulterklopfer/cyphernode_admin/models"
  "github.com/schulterklopfer/cyphernode_admin/queries"
  "github.com/sirupsen/logrus"
  "math/rand"
  "net/http/httptest"
  "os"
  "strconv"
  "testing"
  "time"
)

var testServer *httptest.Server

func TestMain( m *testing.M ) {

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

  code := m.Run()

  testServer.Close()
  dataSource.Close()
  os.Remove(config.DatabaseFile)

  os.Exit(code)
}

func TestEveryting( t *testing.T ) {
  t.Run("apps", testAppHandlers )
  t.Run("users", testUserHandlers )
}