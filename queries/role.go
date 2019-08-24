package queries

import (
  "errors"
  "github.com/go-validator/validator"
  "github.com/schulterklopfer/cyphernode_admin/dataSource"
  "github.com/schulterklopfer/cyphernode_admin/models"
)

func CreateRole( role *models.RoleModel ) error {

  if role.ID != 0 {
    // role must not have any ID possibly existing in DB
    return errors.New( "Role ID must be 0" )
  }

  db := dataSource.GetDB()
  err := validator.Validate(role)
  if err != nil {
    return err
  }
  db.Create(role)
  return nil
}

func DeleteRole( id uint ) error {
  if id == 0 {
    return errors.New("No such role")
  }
  db := dataSource.GetDB()
  var role models.RoleModel
  db.Take( &role, id )
  if role.ID == 0 {
    return errors.New("No such role")
  }
  db.Unscoped().Delete( &role)
  role.ID = 0
  return nil
}
