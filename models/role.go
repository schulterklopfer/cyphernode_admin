package models

import (
  "github.com/jinzhu/gorm"
)

type RoleModel struct {
  gorm.Model
  Name string `json:"name" gorm:"type:varchar(30)" validate:"min=3,max=30,regexp=^[a-zA-Z0-9_-]+$"`
  Description string `json:"description" gorm:"type:varchar(255)"`
  AutoAssign bool `json:"autoAssign" gorm:"default false"`
  AppId uint `json:"appId"`
  Users []*UserModel `json:"users" gorm:"many2many:user_roles;"`
}

func ( role *RoleModel ) AfterDelete( tx *gorm.DB ) {
  role.removeFromAllUsers( tx )
}

func ( role *RoleModel ) AfterSave( tx *gorm.DB ) {
  role.AfterUpdate( tx )
}

func ( role *RoleModel ) AfterUpdate( tx *gorm.DB ) {
  if role.AutoAssign {
    role.addToAllUsers( tx )
  } else {
    role.removeFromAllUsers( tx )
  }
}

func ( role *RoleModel) AfterCreate( tx *gorm.DB )  {
  if !role.AutoAssign {
    return
  }
  role.addToAllUsers( tx )
}

func ( role *RoleModel) addToAllUsers( tx *gorm.DB ) {
  var allUsers []*UserModel
  tx.Find( &allUsers )
  for i:=0; i< len(allUsers); i++ {
    tx.Model(allUsers[i]).Association("Roles").Append( role )
  }
}

func ( role *RoleModel) removeFromAllUsers( tx *gorm.DB ) {
  tx.Model(role).Association("Users" ).Clear()
}


