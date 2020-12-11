package main

import (
  "github.com/schulterklopfer/cyphernode_admin/cyphernodeAdmin"
  "github.com/schulterklopfer/cyphernode_admin/globals"
  "github.com/schulterklopfer/cyphernode_admin/helpers"
  "github.com/schulterklopfer/cyphernode_admin/logwrapper"
  "github.com/sirupsen/logrus"
  "os"
)
func main() {

  logwrapper.Logger().SetLevel(logrus.TraceLevel)

  app := cyphernodeAdmin.NewCyphernodeAdmin( &cyphernodeAdmin.Config{
      DatabaseFile: helpers.GetenvOrDefault(globals.CNA_ADMIN_DATABASE_FILE_ENV_KEY ),
      InitialAdminLogin: helpers.GetenvOrDefault( globals.CNA_ADMIN_LOGIN_ENV_KEY ),
      InitialAdminPassword: helpers.GetenvOrDefault(globals.CNA_ADMIN_PASSWORD_ENV_KEY ),
      InitialAdminName: helpers.GetenvOrDefault(globals.CNA_ADMIN_NAME_ENV_KEY ),
      InitialAdminEmailAddress: helpers.GetenvOrDefault(globals.CNA_ADMIN_EMAIL_ADDRESS_ENV_KEY ),
    },
  )
  err := app.Init()
  if err != nil {
    println("Error in application init: ", err.Error() )
    os.Exit(1)
  }
  app.Start()
}