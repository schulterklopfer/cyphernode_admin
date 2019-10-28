package models

import (
  "github.com/jinzhu/gorm"
  "github.com/schulterklopfer/cyphernode_admin/cnaErrors"
)

type AppModel struct {
  gorm.Model
  ClientSecret   string       `json:"clientSecret" gorm:"type:varchar(32);unique_index;not null" validate:"regexp=^[a-fA-F0-9]{32}$"`
  ClientID       string       `json:"clientID" gorm:"type:varchar(32);unique_index;not null" validate:"regexp=^[a-fA-F0-9]{32}$"`
  Name           string       `json:"name" gorm:"type:varchar(30);not null" validate:"min=3,max=30,regexp=^[a-zA-Z0-9_\\- ]+$"`
  Description    string       `json:"description" gorm:"type:varchar(255)"`
  AvailableRoles []*RoleModel `json:"availableRoles" gorm:"foreignkey:AppId;preload"`
  HydraClientID  uint         `json:"-" gorm:"DEFAULT:0" form:"-" validate:"-"`
  CallbackURL    string       `json:"callbackURL" gorm:"type:varchar(255);not null" form:"-" validate:"required"`
  PostLogoutCallbackURL string `json:"postLogoutCallbackURL" gorm:"type:varchar(255);not null" form:"-" validate:"required"`
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

func ( app *AppModel ) BeforeDelete( tx *gorm.DB ) (err error) {
  // very important. if no check, will delete all users if ID == 0
  if app.ID == 0 {
    err = cnaErrors.ErrNoSuchApp
    return
  }
  return
}