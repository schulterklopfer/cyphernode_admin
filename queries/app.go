package queries

import (
  "errors"
  "github.com/go-validator/validator"
  "github.com/schulterklopfer/cyphernode_admin/cnaErrors"
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

func RemoveRoleFromApp(  app *models.AppModel, roleId uint ) error {
  db := dataSource.GetDB()

  var role models.RoleModel

  err := Get( &role, roleId, false )

  if err != nil {
    return err
  }

  if role.ID == 0 || role.AppId != app.ID {
    return cnaErrors.ErrNoSuchRole
  }

  db.Model(app).Association("AvailableRoles").Delete( role )
  return db.Error
}

func CreateRoleForApp( app *models.AppModel, role *models.RoleModel ) error {
  db := dataSource.GetDB()

  if role.ID != 0 {
    return cnaErrors.ErrCannotAddExistingRole
  }

  db.Model(app).Association("AvailableRoles").Append( role )
  return db.Error
}
