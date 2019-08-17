package models

import "github.com/jinzhu/gorm"

type RoleModel struct {
  gorm.Model
  Name string `json:"name"`
  Description string `json:"description" gorm:"type:varchar(255)"`
  AutoAssign bool `json:"autoAssign" gorm:"default false"`
  AppId uint `json:"appId"`
}

