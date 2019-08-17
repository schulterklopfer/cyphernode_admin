package main

import "github.com/schulterklopfer/cyphernode_admin/app"

func main() {
  appConfig := new( app.Config )
  appConfig.DatabaseFile = "/tmp/test.sqlite3"
  app := app.NewApp( appConfig )
  app.Init()
  app.Start()
}