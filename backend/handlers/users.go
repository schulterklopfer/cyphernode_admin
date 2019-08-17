package handlers

import (
  "encoding/json"
  "github.com/gin-gonic/gin"
  "github.com/schulterklopfer/cyphernode_admin/queries"
  "net/http"
  "strconv"
)

func GetUser(c *gin.Context) {
  // param 0 is first param in url pattern
  id, err := strconv.Atoi(c.Params[0].Value)

  if err != nil {
    if gin.DebugMode == "debug" {
      c.JSON(http.StatusNotFound, gin.H{ "message": err.Error() } )
    } else {
      c.Status(http.StatusNotFound )
    }
    return
  }

  user := queries.GetUser( uint(id), true )

  if user == nil {
    c.Status(http.StatusNotFound )
    return
  }

  userJson, err := json.Marshal(user)

  if err != nil {
    if gin.DebugMode == "debug" {
      c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
    } else {
      c.Status(http.StatusInternalServerError )
    }
    return
  }

  c.JSON(http.StatusCreated, userJson )
}
