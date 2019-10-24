package queries

import (
  "github.com/schulterklopfer/cyphernode_admin/dataSource"
  "github.com/schulterklopfer/cyphernode_admin/models"
)

func DeleteSession( sessionID string ) error {
  if sessionID == "" {
    return nil
  }

  var sessions []*models.SessionModel
  err := Find( &sessions,  []interface{}{"session_id = ?", sessionID }, "", 1,0,false)

  if err != nil {
    return err
  }

  if len(sessions) == 1 {
    db := dataSource.GetDB()
    return db.Unscoped().Delete( &sessions[0] ).Error
  }

  return nil
}