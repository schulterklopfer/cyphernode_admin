package handlers

import (
  "github.com/gin-gonic/gin"
  "github.com/schulterklopfer/cyphernode_admin/globals"
  "github.com/schulterklopfer/cyphernode_admin/helpers"
  "github.com/schulterklopfer/cyphernode_admin/models"
  "github.com/schulterklopfer/cyphernode_admin/queries"
  "github.com/schulterklopfer/cyphernode_admin/transforms"
  "net/http"
)


func CheckSession(c *gin.Context) {

  sessionid, _ := c.Cookie(globals.SESSION_COOKIE_NAME)

  println( sessionid )

  //user, err := cnaOIDC.GetUser(c.Writer, c.Request)
  //if err != nil  {
  //  c.Header("X-Status-Reason", err.Error() )
  //  c.Status(http.StatusUnauthorized)
  //  return
  //}

  //c.JSON(200, map[string]interface{}{
  //  "subject" : user.UserID,
  //  "extra" : user,
  //})
}

func GetSession(c *gin.Context) {
  // param 0 is first param in url pattern
  sessionID := c.Params[0].Value

  var sessions []models.SessionModel

  err := queries.Find( &sessions,  []interface{}{"session_id = ?", sessionID }, "", 1,0,false)

  if err != nil {
    c.Status(http.StatusInternalServerError)
    return
  }

  if len(sessions) == 0 {
    c.Status(http.StatusNotFound )
    return
  }

  session := sessions[0]

  if session.ID == 0 {
    c.Status(http.StatusNotFound )
    return
  }

  var transformedSession transforms.SessionV0

  if transforms.Transform( &session, &transformedSession ) {
    c.JSON(http.StatusOK, transformedSession )
    return
  }
}

func CreateSession(c *gin.Context) {

  input := new( map[string]interface{} )

  err := c.Bind( &input )

  var session models.SessionModel

  helpers.SetByJsonTag( &session, input )

  err = queries.Update( &session )

  if err != nil {
    c.Header("X-Status-Reason", err.Error() )
    c.Status(http.StatusBadRequest)
    return
  }

  var transformedSession transforms.SessionV0
  transforms.Transform( &session, &transformedSession )
  c.JSON( http.StatusCreated, &transformedSession )

}

func PatchSession(c *gin.Context) {
  sessionID := c.Params[0].Value

  var sessions []models.SessionModel

  err := queries.Find( &sessions,  []interface{}{"session_id = ?", sessionID}, "", 1,0,false)

  if err != nil {
    c.Status(http.StatusInternalServerError)
    return
  }

  if len(sessions) == 0 {
    c.Status(http.StatusNotFound )
    return
  }

  session := sessions[0]

  if session.ID == 0 {
    c.Status(http.StatusNotFound )
    return
  }

  input := new( map[string]interface{} )

  err = c.Bind( &input )

  helpers.SetByJsonTag( &session, input )

  err = queries.Update( &session )

  if err != nil {
    c.Header("X-Status-Reason", err.Error() )
    c.Status(http.StatusBadRequest)
    return
  }

  var transformedSession transforms.SessionV0
  transforms.Transform( &session, &transformedSession)
  c.JSON( http.StatusOK, &transformedSession)
}

func DeleteSession(c *gin.Context) {
  sessionID := c.Params[0].Value

  err := queries.DeleteSession( sessionID )

  if err != nil {
    c.Status(http.StatusInternalServerError)
    return
  }

  c.Status(http.StatusNoContent)

}

