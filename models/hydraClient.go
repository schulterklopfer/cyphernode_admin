package models

import (
  "github.com/jinzhu/gorm"
  "github.com/schulterklopfer/cyphernode_admin/cnaErrors"
)

type HydraClientModel struct {
  gorm.Model
  App *AppModel `json:"-" gorm:"foreignkey:AppID" form:"-" validate:"-"`
  AppID uint `json:"-" gorm:"DEFAULT:0" form:"-" validate:"-"`
  ClientID string `json:"-" gorm:"type:varchar(100)" form:"-" validate:"-"`
  Secret string `json:"-" gorm:"type:varchar(32)" form:"-" validate:"-"`
  CallbackURL string `json:"-" gorm:"type:varchar(255)" form:"-" validate:"-"`
  PostLogoutCallbackURL string `json:"-" gorm:"type:varchar(255)" form:"-" validate:"-"`
  Synced bool `json:"-" gorm:"type:boolean;default:false" form:"-" validate:"-"`
}

func ( hydraClient *HydraClientModel ) BeforeDelete( tx *gorm.DB ) (err error) {
  // very important. if no check, will delete all users if ID == 0
  if hydraClient.ID == 0 {
    err = cnaErrors.ErrNoSuchHydraClient
    return
  }
  return
}