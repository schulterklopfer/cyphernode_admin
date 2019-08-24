package models

import (
  "github.com/jinzhu/gorm"
)

type UserModel struct {
  gorm.Model
  Login string `json:"login" gorm:"type:varchar(30);unique_index;not null" validate:"min=3,max=30,regexp=^[a-zA-Z0-9_\\-]+$"`
  Name string `json:"name"` // optional
  Password string `json:"password" gorm:"type:varchar(128);not null" validate:"nonzero"`
  EmailAddress string `json:"emailAddress" gorm:"type:varchar(100)" validate:"max=100,regexp=(^$|^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\\.[a-zA-Z0-9-.]+$)"`
  Roles []*RoleModel `json:"roles" gorm:"many2many:user_roles;"`
}

func (user *UserModel) AfterCreate()  {
  // TODO: assign all auto assignable roles to this user
}

func (user *UserModel) AfterDelete()  {
  // TODO: remove all role associations with this user
}