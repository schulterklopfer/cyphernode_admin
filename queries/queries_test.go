package queries_test

import (
  "github.com/schulterklopfer/cyphernode_admin/dataSource"
  "github.com/schulterklopfer/cyphernode_admin/queries"
  "os"
  "testing"
)

func TestModels(t *testing.T) {
  dbFile := "/tmp/tests.sqlite3"
  os.Remove(dbFile)
  dataSource.Init(dbFile)
  t.Run("Get user", getUser )
}


func getUser( t *testing.T ) {
  user := queries.GetUser(1, true)

  if user == nil {
    t.Error( "unable to load user")
  }

}
