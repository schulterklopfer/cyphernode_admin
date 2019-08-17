package dataSource

import (
  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/sqlite"
  "github.com/schulterklopfer/cyphernode_admin/dataSource/models"
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
    db, err = gorm.Open("sqlite3", databaseFile )
    if err != nil {
      panic("failed to connect to database" )
    }
    println("Connected to db" )
    AutoMigrate()
  })
}

func AutoMigrate() {
  if db == nil {
    return
  }
  println("Migrating models" )
  db.AutoMigrate(&models.UserModel{},&models.AppModel{},&models.RoleModel{})
}