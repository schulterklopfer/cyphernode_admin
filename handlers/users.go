package handlers

import (
  "github.com/gin-gonic/gin"
  "github.com/go-validator/validator"
  "github.com/schulterklopfer/cyphernode_admin/cnaErrors"
  "github.com/schulterklopfer/cyphernode_admin/helpers"
  "github.com/schulterklopfer/cyphernode_admin/models"
  "github.com/schulterklopfer/cyphernode_admin/queries"
  "github.com/schulterklopfer/cyphernode_admin/shared"
  "github.com/schulterklopfer/cyphernode_admin/transforms"
  "net/http"
  "strconv"
  "strings"
)

var ALLOWED_USER_PROPERTIES = [4]string{ "id", "name","login","email_address" }

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
}

func CreateUser(c *gin.Context) {

  input := new( map[string]interface{} )
  err := c.Bind( input )

  if err != nil {
    c.Status(http.StatusInternalServerError)
    return
  }

  var user models.UserModel
  shared.SetByJsonTag( &user, input )

  updateRoles( input, &user )

  // just to be sure
  err = queries.Create(&user)

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

func UpdateUser(c *gin.Context) {
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

  input := new( map[string]interface{} )

  err = c.Bind( &input )

  var newUser models.UserModel

  shared.SetByJsonTag( &newUser, input )

  newUser.ID = user.ID
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

  var transformedUser transforms.UserV0
  transforms.Transform( &newUser, &transformedUser )
  c.JSON( http.StatusOK, &transformedUser )

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
    c.Status(http.StatusInternalServerError)
    return
  }

  if user.ID == 0 {
    c.Status(http.StatusNotFound )
    return
  }

  input := new( map[string]interface{} )

  err = c.Bind( &input )

  shared.SetByJsonTag( &user, input )

  err = queries.Update( &user )

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

func updateRoles( input *map[string]interface{}, user *models.UserModel ) {
  if rolesInput, foundRoles := (*input)["roles"]; foundRoles {
    switch rolesInput.(type) {
    case []interface{}:
      roleCount := len(rolesInput.([]interface{}))
      user.Roles = make( []*models.RoleModel, roleCount )
      if roleCount > 0 {
        user.Roles = make( []*models.RoleModel, roleCount )
        for i:=0; i<roleCount; i++ {
          elem := rolesInput.([]interface{})[i]
          switch elem.(type) {
          case map[string]interface{}:
            if roleId, foundRoleID := elem.(map[string]interface{})["ID"]; foundRoleID {
              role := new( models.RoleModel )
              role.ID = uint( roleId.(float64))
              user.Roles[i] = role
            }
          }
        }
      }
    }
  }
}