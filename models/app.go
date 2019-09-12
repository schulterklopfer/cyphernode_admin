package models

import (
  "github.com/jinzhu/gorm"
)

type AppModel struct {
  gorm.Model
  Hash string `json:"hash" gorm:"type:varchar(32);unique_index;not null" validate:"regexp=^[a-fA-F0-9]{32}$"`
  Name string `json:"name" gorm:"type:varchar(30);not null" validate:"min=3,max=30,regexp=^[a-zA-Z0-9_\\- ]+$"`
  Description string `json:"description" gorm:"type:varchar(255)"`
  AvailableRoles []*RoleModel `json:"availableRoles" gorm:"foreignkey:AppId;preload"`
}

func ( app *AppModel ) AfterDelete( tx *gorm.DB ) {
	var roles []RoleModel
	tx.Model(app).Association("AvailableRoles" ).Find(&roles)
	for i:=0; i< len(roles); i++ {
		tx.Delete( roles[i] )
		// Why do I have to call this manually?
		roles[i].AfterDelete( tx )
	}
}