package handlers

import (
  "github.com/gin-gonic/gin"
  "github.com/schulterklopfer/cyphernode_admin/helpers"
  "github.com/schulterklopfer/cyphernode_admin/models"
  "github.com/schulterklopfer/cyphernode_admin/queries"
  "github.com/schulterklopfer/cyphernode_admin/transforms"
  "net/http"
  "strconv"
  "strings"
)

var ALLOWED_USER_PROPERTIES = [4]string{ "id", "name","login","emailAddress" }

func GetUser(c *gin.Context) {
  // param 0 is first param in url pattern
  id, err := strconv.Atoi(c.Params[0].Value)

  if err != nil {
    c.Status(http.StatusNotFound )
    return
  }

  var user models.UserModel

  err = queries.Get( &user, uint(id), true )

  if err != nil {
    c.Status(http.StatusInternalServerError)
    return
  }

  if user.ID == 0 {
    c.Status(http.StatusNotFound )
    return
  }

  var transformedUser transforms.UserV0

  if transforms.Transform( &user, &transformedUser ) {
    c.JSON(http.StatusOK, transformedUser )
    return
  }

  c.Status(http.StatusNotFound )
}

type Paging struct {
  Page uint `form:"_page"`
  Limit int `form:"_limit"`
  Sort string `form:"_sort"`
  Order string `form:"_order"`
}

func FindUsers(c *gin.Context) {
  var userQuery transforms.UserV0
  var paging Paging

  where := make( []interface{}, 0 )
  order := ""
  offset := uint(0)
  limit := -1

  if c.Bind(&userQuery) == nil {
    fields := make( []string, 0 )
    args := make( []interface{}, 0 )
    if userQuery.Name != "" {
      fields = append( fields, "name LIKE ?" )
      args = append( args, userQuery.Name+"%" )
    }

    if userQuery.Login != "" {
      fields = append( fields, "login LIKE ?" )
      args = append( args, userQuery.Login+"%" )
    }

    if userQuery.EmailAddress != "" {
      fields = append( fields, "email_address LIKE ?" )
      args = append( args, userQuery.EmailAddress+"%" )
    }

    if len(fields) > 0 {
      where = append( where, strings.Join( fields, " AND ") )
      where = append( where, args...)
    }

  }

  if c.Bind(&paging) == nil {

    // is Sort empty or not in ALLOWED_USER_PROPERTIES?
    if paging.Sort != "" ||
       helpers.SliceIndex( len(ALLOWED_USER_PROPERTIES), func(i int) bool { return ALLOWED_USER_PROPERTIES[i] == paging.Sort } ) == -1 {
      order = "name"
    } else {
      order = paging.Sort
    }

    if paging.Order != "" || ( paging.Order != "ASC" && paging.Order != "DESC" ) {
      order = order + " asc"
    } else {
      order = order + " "+strings.ToLower(paging.Order)
    }
  }

  // makes no sense to request 0 users
  // we assume user wants no limit
  if paging.Limit > 0 {
    limit = paging.Limit
  }

  if paging.Page > 0 && limit > 0 {
    offset = (paging.Page-1)*uint(limit)
  }

  var users []*models.UserModel

  err := queries.Find( &users, where, order, limit, offset, true )

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

}
