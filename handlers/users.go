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

package handlers

import (
  "encoding/json"
  "github.com/dgrijalva/jwt-go"
  "github.com/gin-gonic/gin"
  "github.com/schulterklopfer/cyphernode_admin/cnaErrors"
  "github.com/schulterklopfer/cyphernode_admin/helpers"
  "github.com/schulterklopfer/cyphernode_admin/models"
  "github.com/schulterklopfer/cyphernode_admin/queries"
  "github.com/schulterklopfer/cyphernode_admin/transforms"
  "gopkg.in/validator.v2"
  "net/http"
  "strconv"
  "strings"
)

func GetUser(c *gin.Context) {
  // param 0 is first param in url pattern

  idString := c.Params[0].Value
  var id int

  if idString == "me" {
    // my user based on header
    claimsBas64String := c.Request.Header.Get("x-auth-user-claims" )
    claimsJsonString, err := jwt.DecodeSegment( claimsBas64String )

    if err != nil {
      println( err.Error() )
      c.Status(http.StatusBadRequest)
      return
    }

    var claims map[string]interface{}
    err = json.Unmarshal( []byte(claimsJsonString), &claims )

    if err != nil {
      c.Status(http.StatusBadRequest)
      return
    }

    idInterface, exists := claims["id"]

    if !exists {
      c.Status(http.StatusBadRequest)
      return
    }

    floatId, ok := idInterface.(float64)

    if !ok {
      c.Status(http.StatusBadRequest)
      return
    }

    id = int(floatId)

  } else {
    var err error
    id, err = strconv.Atoi(idString)
    if err != nil {
      c.Status(http.StatusBadRequest )
      return
    }
  }

  var user models.UserModel

  err := queries.Get( &user, uint(id), true )

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
}

func CreateUser(c *gin.Context) {

  input := new( map[string]interface{} )
  err := c.Bind( input )

  if err != nil {
    c.Status(http.StatusInternalServerError)
    return
  }

  var user models.UserModel

  helpers.SetByJsonTag( &user, input )

  // just to be sure
  err = queries.CreateUser(&user)

  if err != nil {
    switch err {
    case cnaErrors.ErrDuplicateUser:
      c.Header("X-Status-Reason", err.Error() )
      c.Status(http.StatusConflict )
      return
    case cnaErrors.ErrUserHasUnknownRole:
      c.Header("X-Status-Reason", err.Error() )
      c.Status(http.StatusBadRequest )
      return
    default:
      switch err.(type) {
      case validator.ErrorMap:
        c.Header("X-Status-Reason", err.Error() )
        c.Status(http.StatusBadRequest)
        return
      }
    }
    c.Status(http.StatusInternalServerError)
    return
  }

  queries.LoadRoles( &user )

  var transformedUser transforms.UserV0
  transforms.Transform( &user, &transformedUser )
  c.JSON( http.StatusCreated, &transformedUser )

}

func PatchUser(c *gin.Context) {
  id, err := strconv.Atoi(c.Params[0].Value)
  if err != nil {
    c.Status(http.StatusNotFound )
    return
  }

  var user models.UserModel

  err = queries.Get( &user, uint(id), true )

  if err != nil {
    c.Header("X-Status-Reason", err.Error() )
    c.Status(http.StatusInternalServerError)
    return
  }

  if user.ID == 0 {
    c.Status(http.StatusNotFound )
    return
  }

  input := new( map[string]interface{} )

  err = c.Bind( &input )

  if err != nil {
    c.Header("X-Status-Reason", err.Error() )
    c.Status(http.StatusInternalServerError)
    return
  }

  helpers.SetByJsonTag( &user, input )
  err = queries.UpdateUser( &user )

  if err != nil {
    switch err {
    case cnaErrors.ErrDuplicateUser:
      c.Header("X-Status-Reason", err.Error() )
      c.Status(http.StatusConflict )
      return
    case cnaErrors.ErrUserHasUnknownRole:
      c.Header("X-Status-Reason", err.Error() )
      c.Status(http.StatusBadRequest )
      return
    default:
      switch err.(type) {
      case validator.ErrorMap:
        c.Header("X-Status-Reason", err.Error() )
        c.Status(http.StatusBadRequest)
        return
      }
    }
    c.Status(http.StatusInternalServerError)
    return
  }

  err = queries.LoadRoles( &user )
  if err != nil {
    c.Header("X-Status-Reason", err.Error() )
    c.Status(http.StatusInternalServerError)
    return
  }

  var transformedUser transforms.UserV0
  transforms.Transform( &user, &transformedUser )
  c.JSON( http.StatusOK, &transformedUser )
}

func DeleteUser(c *gin.Context) {
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

  err = queries.DeleteUser( user.ID )

  if err != nil {
    c.Status(http.StatusInternalServerError)
    return
  }

  c.Status(http.StatusNoContent)

}

func FindUsers(c *gin.Context) {
  var userQuery transforms.UserV0
  var paging PagingParams

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


    if paging.Sort == "" {
      paging.Sort = "login"
    }

    if paging.Order == "" {
      paging.Order = "ASC"
    }

    if paging.Limit == 0 {
      paging.Limit = 20
    }

    // is Sort empty or not in ALLOWED_USER_PROPERTIES?
    if helpers.SliceIndex( len(ALLOWED_USER_PROPERTIES), func(i int) bool {
          return ALLOWED_USER_PROPERTIES[i] == paging.Sort
       } ) == -1 {
      order = "login"
    } else {
      order = paging.Sort
    }

    if ( paging.Order != "ASC" && paging.Order != "DESC" ) {
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

  pagedResult := new( PagedResult )

  pagedResult.Page = paging.Page
  pagedResult.Limit = paging.Limit
  pagedResult.Sort = paging.Sort
  pagedResult.Order = paging.Order
  pagedResult.Data = transformedUsers

  _ = queries.TotalCount( &models.UserModel{}, &pagedResult.Total )

  c.JSON(http.StatusOK, pagedResult)

}

// Roles
func UserPatchRoles(c *gin.Context) {
  id, err := strconv.Atoi(c.Params[0].Value)
  if err != nil {
    c.Status(http.StatusNotFound)
    return
  }

  var user models.UserModel

  err = queries.Get(&user, uint(id), true)

  if err != nil {
    c.Status(http.StatusInternalServerError)
    return
  }

  if user.ID == 0 {
    c.Status(http.StatusNotFound)
    return
  }

  var roleInputs []models.RoleModel

  err = c.Bind(&roleInputs)

  if err != nil {
    c.Status(http.StatusInternalServerError)
    return
  }


}

func UserAddRoles(c *gin.Context) {
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

  var roleInputs []models.RoleModel

  err = c.Bind( &roleInputs )

  if err != nil {
    c.Status(http.StatusInternalServerError)
    return
  }

  for i:=0; i<len( roleInputs ); i++ {
    err = queries.AddRoleToUser( &user, uint( roleInputs[i].ID ) )
    if err != nil {
      switch err {
      case cnaErrors.ErrNoSuchRole:
        c.Header("X-Status-Reason", "Role does not exist" )
        c.Status(http.StatusBadRequest)
      case cnaErrors.ErrUserAlreadyHasRole:
        c.Header("X-Status-Reason", "Trying to add role twice" )
        c.Status(http.StatusBadRequest)
      default:
        c.Status(http.StatusInternalServerError)
      }
      return
    }
  }

  var transformedUser transforms.UserV0

  if transforms.Transform( &user, &transformedUser ) {
    c.JSON(http.StatusOK, transformedUser )
  } else {
    c.Status(http.StatusInternalServerError)
  }
}

func UserRemoveRole(c *gin.Context) {
  id, err := strconv.Atoi(c.Params[0].Value)
  if err != nil {
    c.Status(http.StatusNotFound )
    return
  }
  roleId, err := strconv.Atoi(c.Params[1].Value)
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

  for i:=0; i<len( user.Roles ); i++ {
    if user.Roles[i].ID == uint(roleId) {
      // found role
      // TODO: remove
      err := queries.RemoveRoleFromUser( &user, uint(roleId) )
      if err != nil {
        c.Status(http.StatusInternalServerError)
        return
      }
      c.Status(http.StatusNoContent)
      return
    }
  }
  c.Header("X-Status-Reason", "User does not have that role" )
  c.Status(http.StatusBadRequest)
  return
}
