package queries

import (
  "errors"
  "github.com/schulterklopfer/cyphernode_admin/dataSource"
  "github.com/schulterklopfer/cyphernode_admin/models"
)

func DeleteApp( id uint ) error {
  if id == 0 {
    return errors.New("No such app")
  }
  db := dataSource.GetDB()
  var app models.AppModel
  db.Take( &app, id )
  if app.ID == 0 {
    return errors.New("No such app")
  }
  db.Unscoped().Delete( &app )
  return nil
}


