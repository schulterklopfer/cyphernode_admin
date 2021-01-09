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
  config.DisableAuth = true

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
  t.Run("sessions", testSessionHandlers )
}