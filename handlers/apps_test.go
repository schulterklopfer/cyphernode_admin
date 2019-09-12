package handlers_test

import (
	"bytes"
	"encoding/json"
	"github.com/schulterklopfer/cyphernode_admin/cyphernodeAdmin"
	"github.com/schulterklopfer/cyphernode_admin/transforms"
	"io/ioutil"
	"net/http"
	"testing"
)


func testAppHandlers( t *testing.T ) {
  t.Run( "testGetApp", testGetApp )
  t.Run( "testCreateApp", testCreateApp )
  t.Run( "testAppAddRole", testAppAddRole )
  t.Run( "testAppRemoveRole", testAppRemoveRole )
  t.Run( "testPatchApp", testPatchApp )
  t.Run( "testFindApp", testFindApp )
  t.Run( "testDeleteApp", testDeleteApp )
}

func testGetApp(t *testing.T) {

  res, err := testServer.Client().Get( testServer.URL+"/api/v0/apps/1" )

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

  if app.Name != cyphernodeAdmin.ADMIN_APP_NAME ||
    app.Description != cyphernodeAdmin.ADMIN_APP_DESCRIPTION ||
    app.Hash != cyphernodeAdmin.ADMIN_APP_HASH ||
  	app.ID != 1 ||

    len(app.AvailableRoles) != 1 {
    t.Error( "error in get app" )
  }
}

func testCreateApp(t *testing.T) {

  jsonInput := `
{
  "name": "testApp",
  "description": "test app",
  "hash": "01234567890123456789012345678912",
  "availableRoles": [
    {
			"name": "admin",
			"description": "admin",
			"autoAssign": false
		},
		{
			"name": "user",
			"description": "user",
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
      app.Hash != "01234567890123456789012345678912" ||
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
  res, err := testServer.Client().Post( testServer.URL+"/api/v0/apps/2/roles", "application/json", bytes.NewBuffer([]byte(jsonInput)) )

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
			app.Hash != "01234567890123456789012345678912" ||
			app.Description != "test app" ||
			len(app.AvailableRoles) != 3 {
		t.Error( "error in app add role" )
	}
}

func testAppRemoveRole( t *testing.T ) {
  req, err := http.NewRequest("DELETE", testServer.URL+"/api/v0/apps/2/roles/3",nil)

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
  req, err := http.NewRequest("PATCH", testServer.URL+"/api/v0/apps/2",bytes.NewBuffer([]byte(jsonInput)))
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
			app.Hash != "01234567890123456789012345678912" ||
			app.Description != "test app" ||
			len(app.AvailableRoles) != 2 {
		t.Error( "error in patch app" )
	}
}

type pagedApps struct {
  Page int `json:"page"`
  Limit int` json:"limit"`
  Sort string `json:"sort"`
  Order string `json:"order"`
  Total int `json:"total"`
  Data []*transforms.AppV0 `json:"data"`
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

  var pagedResult pagedApps
  json.Unmarshal(body, &pagedResult)

  if pagedResult.Total != 2 ||
    pagedResult.Page != 0 ||
    len(pagedResult.Data) != 1 ||
    pagedResult.Data[0].Name != "testApp2" ||
    pagedResult.Data[0].ID != 2 ||
    len(pagedResult.Data[0].AvailableRoles) != 2 {
    t.Error( "error in find app" )
  }
}


func testDeleteApp(t *testing.T) {

  req, err := http.NewRequest("DELETE", testServer.URL+"/api/v0/apps/2",nil)
  //req.Header.Set("Content-Type", "application/json")

  res, err := testServer.Client().Do(req)

  if err != nil {
    t.Error(err)
  }

  if res.StatusCode != http.StatusNoContent {
    t.Error( "could not delete")
  }
}


