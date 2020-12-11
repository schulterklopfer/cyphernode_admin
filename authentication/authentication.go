package authentication

import (
  "github.com/schulterklopfer/cyphernode_admin/cnaErrors"
  "github.com/schulterklopfer/cyphernode_admin/models"
  "github.com/schulterklopfer/cyphernode_admin/password"
  "github.com/schulterklopfer/cyphernode_admin/queries"
)

func CheckUserPassword( login string, pwString string ) (*models.UserModel, error) {
  var users []*models.UserModel
  err := queries.Find( &users,  []interface{}{"login = ?", login }, "", 1,0,true)

  if err != nil {
    return nil,err
  }

  if len(users) != 1 {
    return nil, cnaErrors.ErrNoSuchUser
  }

  if !password.CheckPasswordHash( pwString, users[0].Password ) {
    return nil, cnaErrors.ErrLoginOrPasswordWrong
  }

  return users[0], nil
}
