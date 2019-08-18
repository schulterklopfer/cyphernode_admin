package main

import (
  "github.com/schulterklopfer/cyphernode_admin/app"
  "os"
)
func main() {
  appConfig := new( app.Config )
  initialAdminLogin := os.Getenv("CNA_ADMIN_LOGIN")
  initialAdminPassword := os.Getenv("CNA_ADMIN_PASSWORD")
  initialAdminName := os.Getenv("CNA_ADMIN_NAME")
  initialAdminEmailAddress := os.Getenv("CNA_ADMIN_EMAIL_ADDRESS")

  if initialAdminLogin == "" {
    initialAdminLogin = "admin"
  }

  if initialAdminPassword == "" {
    initialAdminPassword = "admin"
  }

  if initialAdminName == "" {
    initialAdminName = "Administrator"
  }

  if initialAdminEmailAddress == "" {
    initialAdminEmailAddress = "admin@admin.rocks"
  }

  appConfig.DatabaseFile = "/tmp/test.sqlite3"
  appConfig.InitialAdminLogin = initialAdminLogin
  appConfig.InitialAdminPassword = initialAdminPassword
  appConfig.InitialAdminName = initialAdminName
  appConfig.InitialAdminEmailAddress = initialAdminEmailAddress
  app := app.NewApp( appConfig )
  app.Init()
  app.Start()
}