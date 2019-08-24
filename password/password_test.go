package password_test

import (
  "github.com/schulterklopfer/cyphernode_admin/password"
  "testing"
)

func TestHashing(t *testing.T) {
  hashedPassword, err := password.HashPassword( "test123" )
  match := password.CheckPasswordHash( "test123", hashedPassword )

  if err != nil || !match {
    t.Error( "Failed to hash and verify")
  }
}