package queries

import (
  "errors"
  "github.com/go-validator/validator"
  "github.com/schulterklopfer/cyphernode_admin/dataSource"
  "github.com/schulterklopfer/cyphernode_admin/models"
)

func CreateUser( user *models.UserModel ) error {
  // Check if user with same login exists
  // Create user.
  // if user contains roles, check if roles exist

  if user.ID != 0 {
    // user must not have any ID possibly existing in DB
    return errors.New( "user ID must be 0" )
  }

  db := dataSource.GetDB()

  var existingUsers []models.UserModel
  db.Limit(1).Find( &existingUsers, models.UserModel{Login: user.Login} )

  if len(existingUsers) > 0 {
    return errors.New( "user with same login already exists" )
  }

  err := validator.Validate( user )
  if err != nil {
    return err
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
      return errors.New( "cannot create user with unknown role" )
    }
    db.Take( &role, user.Roles[i].ID )
    if role.ID !=  user.Roles[i].ID {
      return errors.New( "cannot create user with unknown role" )
    }
  }

  db.Create( user )

  return nil
}

func DeleteUser( id uint ) error {
  if id == 0 {
    return errors.New("no such user")
  }
  db := dataSource.GetDB()
  var user models.UserModel
  db.Take( &user, id )
  if user.ID == 0 {
    return errors.New("no such user")
  }
  db.Unscoped().Delete( &user )
  user.ID = 0
  return nil
}

