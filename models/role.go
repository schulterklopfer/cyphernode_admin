package models

import "github.com/jinzhu/gorm"

type RoleModel struct {
  gorm.Model
  Name string `json:"name" gorm:"type:varchar(30)" validate:"min=3,max=30,regexp=^[a-zA-Z0-9_-]+$"`
  Description string `json:"description" gorm:"type:varchar(255)"`
  AutoAssign bool `json:"autoAssign" gorm:"default false"`
  AppId uint `json:"appId" validate:"nonzero"`
}
