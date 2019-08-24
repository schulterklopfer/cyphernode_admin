package queries

import (
  "errors"
  "github.com/schulterklopfer/cyphernode_admin/dataSource"
  "github.com/schulterklopfer/cyphernode_admin/models"
)

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
