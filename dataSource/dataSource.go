package dataSource

import (
  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/sqlite"
  "github.com/schulterklopfer/cyphernode_admin/logwrapper"
  "github.com/schulterklopfer/cyphernode_admin/models"
  "sync"
)

var db *gorm.DB
var once sync.Once

func GetDB() *gorm.DB {
  return db
}

func Init( databaseFile string ) {
  once.Do(func() {
    var err error
    logwrapper.StandardLogger.Info( "Opening database")
    db, err = gorm.Open("sqlite3", databaseFile )
    if err != nil {
      logwrapper.StandardLogger.Panic("failed to connect to database" )
    }
    AutoMigrate()
  })
}

func AutoMigrate() {
  if db == nil {
    return
  }
  logwrapper.StandardLogger.Info( "Migrating database")
  db.AutoMigrate(&models.UserModel{},&models.AppModel{},&models.RoleModel{})
}