package models

import (
  "github.com/jinzhu/gorm"
)

type RoleModel struct {
  gorm.Model
  Name string `json:"name" gorm:"type:varchar(30)" validate:"min=3,max=30,regexp=^[a-zA-Z0-9_-]+$"`
  Description string `json:"description" gorm:"type:varchar(255)"`
  AutoAssign bool `json:"autoAssign" gorm:"default false"`
  AppId uint `json:"appId" validate:"nonzero"`
  Users []*UserModel `json:"users" gorm:"many2many:user_roles;"`
}

func ( role *RoleModel ) AfterDelete( tx *gorm.DB ) {
  tx.Model(role).Association("Users" ).Clear()
}

func ( role *RoleModel) AfterCreate( tx *gorm.DB )  {
  if !role.AutoAssign {
    return
  }

  var allUsers []*UserModel
  tx.Find( &allUsers )
  for i:=0; i< len(allUsers); i++ {
    tx.Model(allUsers[i]).Association("Roles").Append( role )
  }

}

