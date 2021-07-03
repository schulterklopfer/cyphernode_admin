/*
 * MIT License
 *
 * Copyright (c) 2021 schulterklopfer/__escapee__
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILIT * Y, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package handlers_test

import (
  "bytes"
  "github.com/schulterklopfer/cyphernode_admin/cyphernodeAdmin"
  "github.com/schulterklopfer/cyphernode_fauth/dataSource"
  "github.com/schulterklopfer/cyphernode_fauth/logwrapper"
  "github.com/schulterklopfer/cyphernode_fauth/test_helpers"
  "github.com/sirupsen/logrus"
  "net/http/httptest"
  "testing"
  "time"
)

var testServer *httptest.Server

func TestMain( m *testing.M ) {

  logwrapper.Logger().SetLevel( logrus.PanicLevel )

  // startup fresh postgres db
  println( "Preparing test env" )
  containerId, err := test_helpers.StartupPostgres()

  if err != nil {
    return
  }

  println( "Waiting 3 seconds for database" )
  time.Sleep(3 * time.Second)

  var config cyphernodeAdmin.Config
  config.DatabaseDsn = "host=localhost port=5432 user=cnadmin password=cnadmin dbname=cnadmin sslmode=disable"
  config.InitialAdminEmailAddress = "email@email.com"
  config.InitialAdminName = "admin"
  config.InitialAdminLogin = "admin"
  config.InitialAdminPassword = "test123"
  config.DisableAuth = true

  cnAdmin := cyphernodeAdmin.NewCyphernodeAdmin( &config )
  err = cnAdmin.Init()

  if err != nil {
    return
  }

  testServer = httptest.NewServer( cnAdmin.Engine() )

  err = createDummyAdminApp( testServer )

  if err != nil {
    println( "Closing test server" )
    testServer.Close()
    return
  }

  err = createDummyAdminUser( testServer )

  if err != nil {
    println( "Closing test server" )
    testServer.Close()
    return
  }

  defer func() {
    println( "Closing test server" )
    testServer.Close()
    println( "Closing database connections" )
    dataSource.Close()
    println( "Removing database container" )
    _ = test_helpers.CleanupContainer(containerId)
  }()

  m.Run()

}

func createDummyAdminApp( testServer *httptest.Server ) error {
  jsonInput := `
{
  "name": "dummy admin app",
  "description": "dummy admin app",
  "clientSecret": "01234567890123456789012345678912",
  "availableRoles": [
    {
      "name": "admin",
      "description": "god admin",
      "autoAssign": false
    },
    {
      "name": "user",
      "description": "admin app user",
      "autoAssign": true
    }
  ]
}
`
  _, err := testServer.Client().Post( testServer.URL+"/api/v0/apps/", "application/json", bytes.NewBuffer([]byte(jsonInput)) )

  if err != nil {
    return err
  }
  return nil
}

func createDummyAdminUser( testServer *httptest.Server ) error {
  jsonInput := `
{
  "login": "admin",
  "name": "admin",
  "email_address": "admin@admin.com",
  "password": "test123"
}
`
  _, err := testServer.Client().Post( testServer.URL+"/api/v0/users/", "application/json", bytes.NewBuffer([]byte(jsonInput)) )

  if err != nil {
    return err
  }
  return nil
}

func TestEveryting( t *testing.T ) {
  t.Run("apps", testAppHandlers )
  t.Run("users", testUserHandlers )
}