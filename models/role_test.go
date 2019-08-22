package models_test

import (
  "github.com/go-validator/validator"
  "github.com/schulterklopfer/cyphernode_admin/models"
  "testing"
)

func TestRoleValidation(t *testing.T) {
  role := new( models.RoleModel )
  err := validator.Validate(role)

  if err == nil {
    t.Error("Should not validate" )
  }

  role.Name = "Login"
  err = validator.Validate(role)

  if err == nil {
    t.Error("Should not validate" )
  }

  role.AppId = 0
  err = validator.Validate(role)

  if err == nil {
    t.Error("Should not validate" )
  }

  role.AppId = 1
  err = validator.Validate(role)

  if err != nil {
    t.Error("Should validate" )
  }

}
