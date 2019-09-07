package queries

import (
  "github.com/schulterklopfer/cyphernode_admin/cnaErrors"
  "github.com/schulterklopfer/cyphernode_admin/dataSource"
  "github.com/schulterklopfer/cyphernode_admin/models"
)

func DeleteUser( id uint ) error {
  if id == 0 {
    return cnaErrors.ErrNoSuchUser
  }
  db := dataSource.GetDB()
  var user models.UserModel
  db.Take( &user, id )
  return db.Unscoped().Delete( &user ).Error
}

func RemoveRoleFromUser(  user *models.UserModel, roleId uint ) error {
  db := dataSource.GetDB()

  var role models.RoleModel

  err := Get( &role, roleId, false )

  if err != nil {
    return err
  }

  if role.ID == 0 {
    return cnaErrors.ErrNoSuchRole
  }

  db.Model(user).Association("Roles").Delete( role )
  return db.Error
}

func AddRoleToUser( user *models.UserModel, roleId uint ) error {
  db := dataSource.GetDB()

  var role models.RoleModel

  err := Get( &role, roleId, false )

  if err != nil {
    return err
  }

  if role.ID == 0 {
    return cnaErrors.ErrNoSuchRole
  }

  for i:=0; i<len( user.Roles ); i++ {
    if user.Roles[i].ID == roleId {
      return cnaErrors.ErrUserAlreadyHasRole
    }
  }

  db.Model(user).Association("Roles").Append( role )
  return db.Error
}