package models

import (
  "github.com/jinzhu/gorm"
)

type UserModel struct {
  gorm.Model
  Login string `json:"login" gorm:"type:varchar(30);unique_index;not null" validate:"min=3,max=30,regexp=^[a-zA-Z0-9_\\-]+$"`
  Name string `json:"name"` // optional
  Password string `json:"password" gorm:"type:varchar(128);not null" validate:"nonzero"`
  EmailAddress string `json:"email_address" gorm:"type:varchar(100)" validate:"max=100,regexp=(^$|^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\\.[a-zA-Z0-9-.]+$)"`
  Roles []*RoleModel `json:"roles" gorm:"many2many:user_roles;"`
}

func (user *UserModel) AfterCreate( tx *gorm.DB )  {
  var allAutoAssignRoles []*RoleModel
  tx.Where( &RoleModel{ AutoAssign: true }).Find( &allAutoAssignRoles )
  for i:=0; i< len(allAutoAssignRoles); i++ {
    tx.Model(user).Association("Roles").Append( allAutoAssignRoles[i] )
  }
}

func (user *UserModel) AfterDelete( tx *gorm.DB )  {
  tx.Model(user).Association("Roles").Clear()
}