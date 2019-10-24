package dataSource

import (
  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/sqlite"
  "github.com/schulterklopfer/cyphernode_admin/logwrapper"
  "github.com/schulterklopfer/cyphernode_admin/models"
)

var db *gorm.DB

func GetDB() *gorm.DB {
  return db
}

func Init( databaseFile string ) {
  if db != nil {
    return
  }
  var err error
  logwrapper.Logger().Info( "Opening database")
  db, err = gorm.Open("sqlite3", databaseFile )
  if err != nil {
    logwrapper.Logger().Panic("failed to connect to database" )
  }
  AutoMigrate()
}

func Close() {
  if db == nil {
    return
  }
  db.Close()
  db = nil
}

func AutoMigrate() {
  if db == nil {
    return
  }
  logwrapper.Logger().Info( "Migrating database")
  db.AutoMigrate(
    &models.UserModel{},
    &models.AppModel{},
    &models.RoleModel{},
    &models.HydraClientModel{},
    &models.SessionModel{} )
}