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

func GetApp(c *gin.Context) {
  // param 0 is first param in url pattern
  id, err := strconv.Atoi(c.Params[0].Value)
  if err != nil {
    c.Status(http.StatusNotFound )
    return
  }

  var app models.AppModel

  err = queries.Get( &app, uint(id), true )

  if err != nil {
    c.Status(http.StatusInternalServerError)
    return
  }

  if app.ID == 0 {
    c.Status(http.StatusNotFound )
    return
  }

  var transformedApp transforms.AppV0

  if transforms.Transform( &app, &transformedApp ) {
    c.JSON(http.StatusOK, transformedApp )
    return
  }
}

func CreateApp(c *gin.Context) {

  input := new( map[string]interface{} )
  err := c.Bind( input )

  if err != nil {
    c.Status(http.StatusInternalServerError)
    return
  }

  var app models.AppModel
  helpers.SetByJsonTag( &app, input )

  // just to be sure
  err = queries.Create(&app)

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

  var transformedApp transforms.AppV0
  transforms.Transform( &app, &transformedApp )
  c.JSON( http.StatusCreated, &transformedApp )

}

func UpdateApp(c *gin.Context) {
  id, err := strconv.Atoi(c.Params[0].Value)
  if err != nil {
    c.Status(http.StatusNotFound )
    return
  }

  var app models.AppModel

  err = queries.Get( &app, uint(id), true )

  if err != nil {
    c.Status(http.StatusInternalServerError)
    return
  }

  if app.ID == 0 {
    c.Status(http.StatusNotFound )
    return
  }

  input := new( map[string]interface{} )

  err = c.Bind( &input )

  var newUser models.AppModel

  helpers.SetByJsonTag( &newUser, input )

  newUser.ID = app.ID
  err = queries.Update( &newUser )

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

  queries.LoadRoles( &newUser )

  var transformedApp transforms.AppV0
  transforms.Transform( &newUser, &transformedApp)
  c.JSON( http.StatusOK, &transformedApp)

}

func PatchApp(c *gin.Context) {
  id, err := strconv.Atoi(c.Params[0].Value)
  if err != nil {
    c.Status(http.StatusNotFound )
    return
  }

  var app models.AppModel

  err = queries.Get( &app, uint(id), true )

  if err != nil {
    c.Status(http.StatusInternalServerError)
    return
  }

  if app.ID == 0 {
    c.Status(http.StatusNotFound )
    return
  }

  input := new( map[string]interface{} )

  err = c.Bind( &input )

  helpers.SetByJsonTag( &app, input )

  err = queries.Update( &app )

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

  queries.LoadRoles( &app )

  var transformedApp transforms.AppV0
  transforms.Transform( &app, &transformedApp)
  c.JSON( http.StatusOK, &transformedApp)
}

func DeleteApp(c *gin.Context) {
  id, err := strconv.Atoi(c.Params[0].Value)
  if err != nil {
    c.Status(http.StatusNotFound )
    return
  }

  var app models.AppModel

  err = queries.Get( &app, uint(id), true )

  if err != nil {
    c.Status(http.StatusInternalServerError)
    return
  }

  if app.ID == 0 {
    c.Status(http.StatusNotFound )
    return
  }

  err = queries.DeleteApp( app.ID )

  if err != nil {
    c.Status(http.StatusInternalServerError)
    return
  }

  c.Status(http.StatusNoContent)

}

func FindApps(c *gin.Context) {
  var appQuery transforms.UserV0
  var paging FindParams

  where := make( []interface{}, 0 )
  order := ""

  if c.Bind(&appQuery) == nil {
    fields := make( []string, 0 )
    args := make( []interface{}, 0 )
    if appQuery.Name != "" {
      fields = append( fields, "name LIKE ?" )
      args = append( args, appQuery.Name+"%" )
    }

    if appQuery.Login != "" {
      fields = append( fields, "login LIKE ?" )
      args = append( args, appQuery.Login+"%" )
    }

    if appQuery.EmailAddress != "" {
      fields = append( fields, "email_address LIKE ?" )
      args = append( args, appQuery.EmailAddress+"%" )
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

    // is Sort empty or not in ALLOWED_APP_PROPERTIES?
    if helpers.SliceIndex( len(ALLOWED_APP_PROPERTIES), func(i int) bool {
      return ALLOWED_APP_PROPERTIES[i] == paging.Sort
    } ) == -1 {
      order = "name"
    } else {
      order = paging.Sort
    }

    if ( paging.Order != "ASC" && paging.Order != "DESC" ) {
      order = order + " asc"
    } else {
      order = order + " "+strings.ToLower(paging.Order)
    }
  }

  var apps []*models.AppModel

  err := queries.Find( &apps, where, order, -1, 0, true )

  if err != nil {
    c.Status(http.StatusInternalServerError)
    return
  }

  appCount := len(apps)
  transformedApps := make( []*transforms.AppV0, appCount )

  for i:=0; i<appCount; i++ {
    transformedApps[i] = new( transforms.AppV0 )
    transforms.Transform( apps[i], transformedApps[i] )
  }

  pagedResult := new(Result)
  pagedResult.Results = transformedApps

  c.JSON(http.StatusOK, pagedResult)

}

// Roles

func AppAddRoles(c *gin.Context) {
  id, err := strconv.Atoi(c.Params[0].Value)
  if err != nil {
    c.Status(http.StatusNotFound )
    return
  }

  var app models.AppModel

  err = queries.Get( &app, uint(id), true )

  if err != nil {
    c.Status(http.StatusInternalServerError)
    return
  }

  if app.ID == 0 {
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
    roleInputs[i].AppId = app.ID

    err = queries.CreateRoleForApp( &app, &roleInputs[i] )

    if err != nil {
      switch err {
      case cnaErrors.ErrCannotAddExistingRole:
        c.Header("X-Status-Reason", "Trying to add an existing role" )
        c.Status(http.StatusBadRequest)
      default:
        c.Status(http.StatusInternalServerError)
      }
      return
    }
  }

  var transformedApp transforms.AppV0

  if transforms.Transform( &app, &transformedApp ) {
    c.JSON(http.StatusOK, transformedApp )
  } else {
    c.Status(http.StatusInternalServerError)
  }
}

func AppRemoveRole(c *gin.Context) {
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

  var app models.AppModel

  err = queries.Get( &app, uint(id), true )

  if err != nil {
    c.Status(http.StatusInternalServerError)
    return
  }

  for i:=0; i<len( app.AvailableRoles ); i++ {
    if app.AvailableRoles[i].ID == uint(roleId) {
      // found role
      // TODO: remove
      err := queries.RemoveRoleFromApp( &app, uint(roleId) )
      if err != nil {
        switch err {
        case cnaErrors.ErrNoSuchRole:
          c.Header("X-Status-Reason", "App does not have that role" )
          c.Status(http.StatusBadRequest)
        default:
          c.Status(http.StatusInternalServerError)
        }
        return
      }
      c.Status(http.StatusNoContent)
      return
    }
  }
  c.Header("X-Status-Reason", "App does not have that role" )
  c.Status(http.StatusBadRequest)
  return
}
