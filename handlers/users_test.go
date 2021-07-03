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

func testUserHandlers( t *testing.T ) {
  t.Run( "testCreateUser", testCreateUser )
  t.Run( "testGetUser", testGetUser )
  t.Run( "testUserAddRole", testUserAddRole )
  t.Run( "testUserRemoveRole", testUserRemoveRole )
  t.Run( "testPatchUser", testPatchUser )
  t.Run( "testFindUser", testFindUser )
  t.Run( "testDeleteUser", testDeleteUser )
}

var createdUserId = uint(0)


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
      len(user.Roles) != 1 {
    t.Error( "error in get user" )
  }

  createdUserId = user.ID
}

func testGetUser(t *testing.T) {

  res, err := testServer.Client().Get( testServer.URL+"/api/v0/users/"+strconv.Itoa(int(createdUserId)) )

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
      user.ID != createdUserId ||
      len(user.Roles) != 1 {
    t.Error( "error in get user" )
  }
}

func testUserAddRole( t *testing.T ) {
  jsonInput := `[{ "ID": 2 }]`
  res, err := testServer.Client().Post( testServer.URL+"/api/v0/users/"+strconv.Itoa(int(createdUserId))+"/roles", "application/json", bytes.NewBuffer([]byte(jsonInput)) )

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
  req, err := http.NewRequest("DELETE", testServer.URL+"/api/v0/users/"+strconv.Itoa(int(createdUserId))+"/roles/2",nil)

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
  req, err := http.NewRequest("PATCH", testServer.URL+"/api/v0/users/"+strconv.Itoa(int(createdUserId)),bytes.NewBuffer([]byte(jsonInput)))
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
      user.ID != createdUserId ||
      len(user.Roles) != 0 {
    t.Error( "error in patch user" )
  }
}

type FindUsersResult struct {
  Error string `json:"error"`
  Results []*transforms.UserV0 `json:"results"`
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

  var pagedResult FindUsersResult
  json.Unmarshal(body, &pagedResult)

  if len(pagedResult.Results) != 1 ||
      pagedResult.Results[0].Login != "testUser2" ||
      pagedResult.Results[0].Name != "name" ||
      pagedResult.Results[0].EmailAddress != "email@test.com" ||
      pagedResult.Results[0].ID != createdUserId ||
      len(pagedResult.Results[0].Roles) != 0 {
    t.Error( "error in find user" )
  }
}


func testDeleteUser(t *testing.T) {

  req, err := http.NewRequest("DELETE", testServer.URL+"/api/v0/users/"+strconv.Itoa(int(createdUserId)),nil)
  //req.Header.Set("Content-Type", "application/json")

  res, err := testServer.Client().Do(req)

  if err != nil {
    t.Error(err)
  }

  if res.StatusCode != http.StatusNoContent {
    t.Error( "could not delete")
  }
}



