package models

import (
  "github.com/jinzhu/gorm"
  "github.com/schulterklopfer/cyphernode_admin/cnaErrors"
)

type SessionModel struct {
  gorm.Model
  SessionID string `json:"sessionID" gorm:"type:varchar(48)" validate:"regexp=^[a-zA-Z0-9_-]+$"`
  Values    string `json:"values" gorm:"type:text"`
}

func ( session *SessionModel ) BeforeDelete( tx *gorm.DB ) (err error) {
  // very important. if no check, will delete all users if ID == 0
  if session.ID == 0 {
    err = cnaErrors.ErrNoSuchSession
    return
  }
  return
}
