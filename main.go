package main

import (
  "github.com/schulterklopfer/cyphernode_admin/cyphernodeAdmin"
  "github.com/schulterklopfer/cyphernode_admin/globals"
  "github.com/schulterklopfer/cyphernode_admin/helpers"
  "os"
)
func main() {
  app := cyphernodeAdmin.NewCyphernodeAdmin( &cyphernodeAdmin.Config{
      DatabaseFile: helpers.GetenvOrDefault(globals.CNA_ADMIN_DATABASE_FILE_ENV_KEY, globals.DEFAULTS_CNA_ADMIN_DATABASE_FILE ),
      InitialAdminLogin: helpers.GetenvOrDefault( globals.CNA_ADMIN_LOGIN_ENV_KEY, globals.DEFAULTS_CNA_ADMIN_LOGIN ),
      InitialAdminPassword: helpers.GetenvOrDefault(globals.CNA_ADMIN_PASSWORD_ENV_KEY, globals.DEFAULTS_CNA_ADMIN_PASSWORD ),
      InitialAdminName: helpers.GetenvOrDefault(globals.CNA_ADMIN_NAME_ENV_KEY, globals.DEFAULTS_CNA_ADMIN_NAME ),
      InitialAdminEmailAddress: helpers.GetenvOrDefault(globals.CNA_ADMIN_EMAIL_ADDRESS_ENV_KEY, globals.DEFAULTS_CNA_ADMIN_EMAIL_ADDRESS ),
    },
  )
  err := app.Init()
  if err != nil {
    println("Error in application init: ", err.Error() )
    os.Exit(1)
  }
  app.Start()
}