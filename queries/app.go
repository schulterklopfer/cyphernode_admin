package queries

import (
  "errors"
  "github.com/go-validator/validator"
  "github.com/schulterklopfer/cyphernode_admin/dataSource"
  "github.com/schulterklopfer/cyphernode_admin/models"
)

func CreateApp( app *models.AppModel ) error {
  if app.ID != 0 {
    // app must not have any ID possibly existing in DB
    return errors.New( "app ID must be 0" )
  }
  db := dataSource.GetDB()

  var existingApps []models.AppModel
  db.Limit(1).Find( &existingApps, models.AppModel{Hash: app.Hash} )

  if len(existingApps) > 0 {
    return errors.New( "app with same hash already exists" )
  }

  err := validator.Validate(app)
  if err != nil {
    return err
  }
  db.Create(app)
  return nil
}

func DeleteApp( id uint ) error {
  if id == 0 {
    return errors.New("no such app")
  }
  db := dataSource.GetDB()
  var app models.AppModel
  db.Take( &app, id )
  if app.ID == 0 {
    return errors.New("no such app")
  }
  db.Unscoped().Delete( &app )
  return nil
}


