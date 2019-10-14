package handlers

import (
  "github.com/gin-gonic/gin"
  "net/http"
)

func GetRolesForApp( c *gin.Context ) {
  hash := c.Params[0].Value
  username := c.Params[1].Value
  if hash == "" || username == "" {
    c.Status( http.StatusNotFound )
    return
  }

  // TODO: check if "roles" scope is set
  // TODO: return roles of user with username in app with given hash

  println( hash, username )
}
