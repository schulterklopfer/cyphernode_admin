package handlers_test

import (
  "bytes"
  "encoding/json"
  "github.com/schulterklopfer/cyphernode_admin/transforms"
  "io/ioutil"
  "net/http"
  "testing"
)


func testSessionHandlers( t *testing.T ) {
  t.Run( "testCreateSession", testCreateSession )
  t.Run( "testGetSession", testGetSession )
  t.Run( "testSessionAddRole", testPatchSession )
  t.Run( "testSessionAddRole", testDeleteSession )
}

var sessionID = "thisisafakesessionid"
var sessionValues1 = "foobar"
var sessionValues2 = "barbaz"

func testGetSession(t *testing.T) {

  res, err := testServer.Client().Get( testServer.URL+"/sessions/"+sessionID )

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

  var session transforms.SessionV0
  _ = json.Unmarshal(body, &session)

  if session.SessionID != sessionID ||
      session.Values != sessionValues1 {
    t.Error( "error in get session" )
  }
}

func testCreateSession(t *testing.T) {

  jsonInput := `
{
  "sessionID": "`+sessionID+`",
  "values": "`+ sessionValues1 +`"
}
`
  res, err := testServer.Client().Post( testServer.URL+"/sessions/", "application/json", bytes.NewBuffer([]byte(jsonInput)) )

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

  var session transforms.SessionV0
  _ = json.Unmarshal(body, &session)

  if session.SessionID != sessionID ||
     session.Values != sessionValues1 {
    t.Error( "error in get session" )
  }
}


func testPatchSession(t *testing.T) {

  jsonInput := `
{
  "values": "`+sessionValues2+`"
}
`
  req, err := http.NewRequest("PATCH", testServer.URL+"/sessions/"+sessionID, bytes.NewBuffer([]byte(jsonInput)))
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

  var session transforms.SessionV0
  _ = json.Unmarshal(body, &session)

  if session.SessionID != sessionID ||
      session.Values != sessionValues2 {
    t.Error( "error in patch session" )
  }
}

func testDeleteSession(t *testing.T) {

  req, err := http.NewRequest("DELETE", testServer.URL+"/sessions/"+sessionID,nil)
  //req.Header.Set("Content-Type", "application/json")

  res, err := testServer.Client().Do(req)

  if err != nil {
    t.Error(err)
  }

  if res.StatusCode != http.StatusNoContent {
    t.Error( "could not delete")
  }
}



