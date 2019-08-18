package handlers

import (
  "github.com/gin-gonic/gin"
  "github.com/schulterklopfer/cyphernode_admin/queries"
  "github.com/schulterklopfer/cyphernode_admin/transforms"
  "net/http"
  "strconv"
)

func GetUser(c *gin.Context) {
  // param 0 is first param in url pattern
  id, err := strconv.Atoi(c.Params[0].Value)

  if err != nil {
    c.Status(http.StatusNotFound )
    return
  }

  user, err := queries.GetUser( uint(id), true )

  if err != nil {
    c.Status(http.StatusInternalServerError)
    return
  }

  if user == nil {
    c.Status(http.StatusNotFound )
    return
  }

  var transformedUser transforms.UserV0

  if transforms.Transform( user, &transformedUser ) {
    c.JSON(http.StatusOK, transformedUser )
    return
  }

  c.Status(http.StatusNotFound )
}
