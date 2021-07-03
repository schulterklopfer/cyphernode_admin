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
  "encoding/json"
  "github.com/schulterklopfer/cyphernode_admin/transforms"
  "io/ioutil"
  "net/http"
  "strconv"
  "testing"
)


func testAppHandlers( t *testing.T ) {
  t.Run( "testCreateApp", testCreateApp )
  t.Run( "testGetApp", testGetApp )
  t.Run( "testAppAddRole", testAppAddRole )
  t.Run( "testAppRemoveRole", testAppRemoveRole )
  t.Run( "testPatchApp", testPatchApp )
  t.Run( "testFindApp", testFindApp )
  t.Run( "testDeleteApp", testDeleteApp )
}

var createdAppId = uint(0)
var createdRoleId = uint(0)

func testCreateApp(t *testing.T) {

  jsonInput := `
{
  "name": "testApp",
  "description": "test app",
  "clientSecret": "01234567890123456789012345678912",
  "availableRoles": [
    {
      "name": "admin",
      "description": "test app admin",
      "autoAssign": false
    },
    {
      "name": "user",
      "description": "test app user",
      "autoAssign": true
    }
  ]
}
`
  res, err := testServer.Client().Post( testServer.URL+"/api/v0/apps/", "application/json", bytes.NewBuffer([]byte(jsonInput)) )

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

  var app transforms.AppV0
  json.Unmarshal(body, &app)

  if app.Name != "testApp" ||
      app.Description != "test app" ||
      len(app.AvailableRoles) != 2 {
    t.Error( "error in create app" )
  }
  createdAppId = app.ID
}

func testGetApp(t *testing.T) {

  res, err := testServer.Client().Get( testServer.URL+"/api/v0/apps/"+strconv.Itoa(int(createdAppId)) )

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

  var app transforms.AppV0
  _ = json.Unmarshal(body, &app)

  if app.Name != "testApp" ||
      app.Description != "test app" ||
      len(app.AvailableRoles) != 2 {
    t.Error( "error in get app" )
  }
}

func testAppAddRole( t *testing.T ) {
  jsonInput := `[
    {
      "name": "additionalRole",
      "description": "additional role",
      "autoAssign": false
    }
  ]`
  res, err := testServer.Client().Post( testServer.URL+"/api/v0/apps/"+strconv.Itoa(int(createdAppId))+"/roles", "application/json", bytes.NewBuffer([]byte(jsonInput)) )

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

  var app transforms.AppV0
  json.Unmarshal(body, &app)


  if app.Name != "testApp" ||
      app.Description != "test app" ||
      len(app.AvailableRoles) != 3 {
    t.Error( "error in app add role" )
  }

  for _, role := range app.AvailableRoles {
    if role.Name == "additionalRole" {
      createdRoleId = role.ID
      break
    }
  }

}

func testAppRemoveRole( t *testing.T ) {
  req, err := http.NewRequest("DELETE", testServer.URL+"/api/v0/apps/"+strconv.Itoa(int(createdAppId))+"/roles/"+strconv.Itoa(int(createdRoleId)),nil)

  res, err := testServer.Client().Do(req)

  if err != nil {
    t.Error(err)
  }

  if res.StatusCode != http.StatusNoContent {
    t.Error( "could not delete")
  }
}

func testPatchApp(t *testing.T) {

  jsonInput := `
{
  "name": "testApp2"
}
`
  req, err := http.NewRequest("PATCH", testServer.URL+"/api/v0/apps/"+strconv.Itoa(int(createdAppId)),bytes.NewBuffer([]byte(jsonInput)))
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

  var app transforms.AppV0
  json.Unmarshal(body, &app)

  if app.Name != "testApp2" ||
      app.Description != "test app" ||
      len(app.AvailableRoles) != 2 {
    t.Error( "error in patch app" )
  }
}

type FindAppsResult struct {
  Error string `json:"error"`
  Results []*transforms.AppV0 `json:"results"`
}

func testFindApp(t *testing.T) {

  res, err := testServer.Client().Get( testServer.URL+"/api/v0/apps/?name_like=testApp" )

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

  var result FindAppsResult
  json.Unmarshal(body, &result)

  if len(result.Results) != 1 ||
    result.Results[0].Name != "testApp2" ||
    result.Results[0].ID != createdAppId ||
    len(result.Results[0].AvailableRoles) != 2 {
    t.Error( "error in find app" )
  }
}


func testDeleteApp(t *testing.T) {

  req, err := http.NewRequest("DELETE", testServer.URL+"/api/v0/apps/"+strconv.Itoa(int(createdAppId)),nil)
  //req.Header.Set("Content-Type", "application/json")

  res, err := testServer.Client().Do(req)

  if err != nil {
    t.Error(err)
  }

  if res.StatusCode != http.StatusNoContent {
    t.Error( "could not delete")
  }
}



