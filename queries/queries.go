package queries

import (
  "errors"
  "github.com/go-validator/validator"
  "github.com/schulterklopfer/cyphernode_admin/dataSource"
  "github.com/schulterklopfer/cyphernode_admin/models"
)

func CreateUser( user *models.UserModel ) (*models.UserModel, error) {
  // Check if user with same login exists
  // Create user.
  // if user contains roles, check if roles exist

  if user.ID != 0 {
    // user must not have any ID possibly existing in DB
    return nil, errors.New( "User ID must be 0" )
  }

  db := dataSource.GetDB()

  var existingUsers []models.UserModel
  db.Limit(1).Find( &existingUsers, models.UserModel{Login: user.Login} )

  if len(existingUsers) > 0 {
    return nil, errors.New( "User with same login already exists" )
  }

  err := validator.Validate( user )
  if err != nil {
    return nil, err
  }

  // Validate roles:
  // For each role:
  // Check if role id == 0
  // if yes: do not create user
  // Check if role is in db
  // if not: do not create user

  var role models.RoleModel
  for i:=0; i<len( user.Roles ); i++ {
    if user.Roles[i].ID == 0 {
      return nil, errors.New( "Cannot create user with unknown role" )
    }
    db.Take( &role, user.Roles[i].ID )
    if role.ID !=  user.Roles[i].ID {
      return nil, errors.New( "Cannot create user with unknown role" )
    }
  }

  db.Create( user )

  return user, nil
}

func DeleteUser( id uint ) (*models.UserModel, error) {
  if id == 0 {
    return nil, errors.New("No such user")
  }
  db := dataSource.GetDB()

  var user models.UserModel

  db.Take( &user, id )

  if user.ID == 0 {
    return nil, errors.New("No such user")
  }

  db.Unscoped().Delete( &user )
  user.ID = 0
  return &user, nil
}

func Get( model interface{}, id uint, recursive bool ) error {
  db := dataSource.GetDB()
  db.Take(model, id)
  if recursive {
    err := loadRoles(model)
    if err != nil {
      return err
    }
  }
  return nil
}

func Find( out interface{}, where []interface{}, order string, limit int, offset uint, recursive bool ) error {

  /*
     where == nil -> no where
     order == "" -> no order
     limit == -1 -> no limit
     offset == 0 -> no offset
  */

  db := dataSource.GetDB()

  if len(where) > 0 {
    db = db.Where( where[0].(string), where[1:] )
  }

  if order != "" {
    db = db.Order( order )
  }

  if limit != -1 {
    db = db.Limit( limit )
  }

  if offset > 0 {
    db = db.Offset( offset )
  }

  db.Find( out )

  if recursive {
    switch out.(type) {
    case *[]*models.UserModel:
      users := *out.(*[]*models.UserModel)
      for i:=0; i<len(users); i++ {
        _ = loadRoles(users[i])
      }
    case *[]*models.AppModel:
      apps := *out.(*[]*models.AppModel)
      for i:=0; i<len(apps); i++ {
        _ = loadRoles(apps[i])
      }
    }
  }

  return db.Error

}

func loadRoles( in interface{} ) error {
  db := dataSource.GetDB()
  var roles []*models.RoleModel
  switch in.(type) {
  case *models.UserModel:
    if in.(*models.UserModel).ID > 0 {
      db.Model(in).Association("Roles").Find(&roles)
      in.(*models.UserModel).Roles = roles
    }
  case *models.AppModel:
    if in.(*models.AppModel).ID > 0 {
      db.Model(in).Association("AvailableRoles").Find(&roles)
      in.(*models.AppModel).AvailableRoles = roles
    }
  }
  return db.Error
}
