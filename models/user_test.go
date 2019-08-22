package models_test

import (
  "github.com/go-validator/validator"
  "github.com/schulterklopfer/cyphernode_admin/models"
  "testing"
)

func TestUserValidation(t *testing.T) {
  user := new( models.UserModel )
  err := validator.Validate( user )

  if err == nil {
    t.Error("Should not validate" )
  }

  user.Login = "Login"
  err = validator.Validate( user )

  if err == nil {
    t.Error("Should not validate" )
  }

  user.EmailAddress = "foobar"
  err = validator.Validate( user )

  if err == nil {
    t.Error("Should not validate" )
  }

  user.Password = "test123"
  err = validator.Validate( user )

  if err == nil {
    t.Error("Should not validate" )
  }

  user.EmailAddress = "foo@bar.de"
  err = validator.Validate( user )

  if err != nil {
    t.Error("Should validate" )
  }

}
