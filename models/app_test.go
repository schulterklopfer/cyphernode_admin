package models_test

import (
  "github.com/go-validator/validator"
  "github.com/schulterklopfer/cyphernode_admin/models"
  "testing"
)

func TestAppValidation(t *testing.T) {
  app := new( models.AppModel )
  err := validator.Validate(app)

  if err == nil {
    t.Error("Should not validate" )
  }

  app.ClientSecret = "abc"
  err = validator.Validate(app)

  if err == nil {
    t.Error("Should not validate" )
  }

  app.ClientSecret = "000000000011111111112222222222"
  err = validator.Validate(app)

  if err == nil {
    t.Error("Should not validate" )
  }

  app.ClientSecret = "00000000001111111111222222222233"
  err = validator.Validate(app)

  if err == nil {
    t.Error("Should not validate" )
  }

  app.Name = "Test app"
  err = validator.Validate(app)
  
  if err != nil {
    t.Error("Should validate" )
  }

}
