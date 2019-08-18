package handlers

import (
  "github.com/gin-gonic/gin"
  "github.com/schulterklopfer/cyphernode_admin/queries"
  "github.com/schulterklopfer/cyphernode_admin/transforms"
  "net/http"
  "strconv"
  "strings"
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

type Paging struct {
  Offset uint `form:"offset"`
  Limit int `form:"limit"`
  Order []string `form:"order"`
}

func FindUsers(c *gin.Context) {
  var userQuery transforms.UserV0
  var paging Paging

  where := make( []interface{}, 0 )
  order := ""

  if c.Bind(&userQuery) == nil {
    fields := make( []string, 0 )
    args := make( []interface{}, 0 )


    if userQuery.ID != 0 {
      fields = append( fields, "id = ?" )
      args = append( args, userQuery.ID )
    }

    if userQuery.Name != "" {
      fields = append( fields, "name LIKE ?" )
      args = append( args, userQuery.Name+"%" )
    }

    if userQuery.Login != "" {
      fields = append( fields, "login LIKE ?" )
      args = append( args, userQuery.Login+"%" )
    }

    if userQuery.EmailAddress != "" {
      fields = append( fields, "emailAddress LIKE ?" )
      args = append( args, userQuery.EmailAddress+"%" )
    }

    if len(fields) > 0 {
      where = append( where, strings.Join( fields, " AND ") )
      where = append( where, args...)
    }

  }

  c.Bind(&paging)

  // makes no sense to request 0 users
  // we assume user wants no limit
  if paging.Limit == 0 {
    paging.Limit=-1
  }

  users, err := queries.FindUsers( where, order, paging.Limit, paging.Offset, true )

  if err != nil {
    c.Status(http.StatusInternalServerError)
    return
  }

  userCount := len(users)
  transformedUsers := make( []*transforms.UserV0, userCount )

  for i:=0; i<userCount; i++ {
    transformedUsers[i] = new( transforms.UserV0 )
    transforms.Transform( users[i], transformedUsers[i] )
  }

  c.JSON(http.StatusOK, transformedUsers)
  return

}
