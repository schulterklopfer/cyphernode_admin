package models

import "github.com/jinzhu/gorm"

type AppModel struct {
  gorm.Model
  Hash string `json:"hash" gorm:"type:varchar(32);unique_index;not null"`
  Name string `json:"name" gorm:"type:varchar(30);not null"`
  Description string `json:"description" gorm:"type:varchar(255)"`
  AvailableRoles []*RoleModel `json:"availableRoles" gorm:"foreignkey:AppId;preload"`
}



