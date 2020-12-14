package queries

import (
  "github.com/schulterklopfer/cyphernode_admin/cnaErrors"
  "github.com/schulterklopfer/cyphernode_admin/dataSource"
  "github.com/schulterklopfer/cyphernode_admin/models"
  "gopkg.in/validator.v2"
)

func CreateUser( user *models.UserModel ) error {
  db := dataSource.GetDB()
  err := validator.Validate( user )
  if err != nil {
    return err
  }
  // update associations, but don't upsert roles.
  return db.Create( user ).Error
}

func UpdateUser( user *models.UserModel ) error {
  db := dataSource.GetDB()
  tx := db.Begin()

  err := validator.Validate( user )
  if err != nil {
    tx.Rollback()
    return err
  }
  err = tx.Model(&user).Association("Roles").Replace(user.Roles).Error
  if err != nil {
    tx.Rollback()
    return err
  }
  err = tx.Save( user ).Error
  if err != nil {
    tx.Rollback()
    return err
  }

  tx.Commit()
  return nil
}


func DeleteUser( id uint ) error {
  if id == 0 {
    return cnaErrors.ErrNoSuchUser
  }
  if id == 1 {
    return cnaErrors.ErrActionForbidden
  }
  db := dataSource.GetDB()
  var user models.UserModel
  db.Take( &user, id )
  return db.Unscoped().Delete( &user ).Error
}

func RemoveRoleFromUser(  user *models.UserModel, roleId uint ) error {
  if roleId == 1 && user.ID == 1 {
    return cnaErrors.ErrActionForbidden
  }

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

func GetRolesOfUser( user *models.UserModel ) ( []*models.RoleModel, error) {
  db := dataSource.GetDB()
  var roles []*models.RoleModel
  db.Model( user ).Association( "Roles" ).Find(&roles)
  return roles, db.Error
}
