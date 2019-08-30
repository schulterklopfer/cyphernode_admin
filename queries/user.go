package queries

import (
  "github.com/schulterklopfer/cyphernode_admin/dataSource"
  "github.com/schulterklopfer/cyphernode_admin/models"
)

func DeleteUser( id uint ) error {
  if id == 0 {
    return models.ErrNoSuchUser
  }
  db := dataSource.GetDB()
  var user models.UserModel
  db.Take( &user, id )
  return db.Unscoped().Delete( &user ).Error
}

