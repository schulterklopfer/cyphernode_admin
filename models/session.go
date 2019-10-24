package models

import (
  "github.com/jinzhu/gorm"
  "github.com/schulterklopfer/cyphernode_admin/cnaErrors"
)

type SessionModel struct {
  gorm.Model
  SessionID string `json:"-" gorm:"type:varchar(48)"`
  Values    string `json:"-" gorm:"type:text"`
}

func ( session *SessionModel ) BeforeDelete( tx *gorm.DB ) (err error) {
  // very important. if no check, will delete all users if ID == 0
  if session.ID == 0 {
    err = cnaErrors.ErrNoSuchSession
    return
  }
  return
}
