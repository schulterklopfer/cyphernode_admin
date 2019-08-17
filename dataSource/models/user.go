package models

import "github.com/jinzhu/gorm"

type UserModel struct {
  gorm.Model
  Login string `json:"login" gorm:"type:varchar(30);unique_index;not null"`
  Name string `json:"name"`
  EmailAddress string `json:"emailAddress" gorm:"type:varchar(100);not null"`
  Roles []*RoleModel `json:"roles" gorm:"many2many:user_roles;"`
}

