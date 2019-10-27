package helpers_test

import (
  "github.com/schulterklopfer/cyphernode_admin/globals"
  "github.com/schulterklopfer/cyphernode_admin/helpers"
  "os"
  "testing"
)

type testStruct struct {
  Aint int `json:"aint"`
  Bint32 int32 `json:"bint32"`
  Cint64 int64 `json:"cint64"`
  Dstring string `json:"dstring"`
  Ebool bool `json:"ebool"`
  Ffloat32 float32 `json:"ffloat32"`
  Gfloat64 float64 `json:"gfloat64"`
}

func TestSetByJsonTag(t *testing.T) {

  target := testStruct{ 1,1,1,"foo", false, 1.0, 1.0 }

  newValues := map[string]interface{}{
    "aint": 2,
    "bint32": int32(3),
    "cint64": int64(4),
    "dstring": "bar",
    "ebool": true,
    "ffloat32": float32(2.0),
    "gfloat64": float64(3.0),
  }

  helpers.SetByJsonTag(  &target, &newValues )

  if target.Aint != 2 ||
      target.Bint32 != 3 ||
      target.Cint64 != 4 ||
      target.Dstring != "bar" ||
      target.Ebool != true ||
      target.Ffloat32 != 2.0 ||
      target.Gfloat64 != 3.0 {
    t.Error( "Set value failed")
  }

}

func TestAbsoluteURL( t *testing.T ) {

  _ = os.Setenv( globals.BASE_URL_ENV_KEY, "http://www.foo.com")

  a := helpers.AbsoluteURL( "bar" )

  if a != "http://www.foo.com/bar" {
    t.Errorf( "%s should be %s", a, "http://www.foo.com/bar" )
  }

  a = helpers.AbsoluteURL( "/bar" )

  if a != "http://www.foo.com/bar" {
    t.Errorf( "%s should be %s", a, "http://www.foo.com/bar" )
  }

  a = helpers.AbsoluteURL( "//bar" )

  if a != "http://www.foo.com/bar" {
    t.Errorf( "%s should be %s", a, "http://www.foo.com/bar" )
  }

  _ = os.Setenv( globals.BASE_URL_ENV_KEY, "http://www.foo.com/")

  a = helpers.AbsoluteURL( "bar" )

  if a != "http://www.foo.com/bar" {
    t.Errorf( "%s should be %s", a, "http://www.foo.com/bar" )
  }

  a = helpers.AbsoluteURL( "/bar" )

  if a != "http://www.foo.com/bar" {
    t.Errorf( "%s should be %s", a, "http://www.foo.com/bar" )
  }

  a = helpers.AbsoluteURL( "//bar" )

  if a != "http://www.foo.com/bar" {
    t.Errorf( "%s should be %s", a, "http://www.foo.com/bar" )
  }

  _ = os.Setenv( globals.BASE_URL_ENV_KEY, "http://www.foo.com//")

  a = helpers.AbsoluteURL( "bar" )

  if a != "http://www.foo.com/bar" {
    t.Errorf( "%s should be %s", a, "http://www.foo.com/bar" )
  }

  a = helpers.AbsoluteURL( "/bar" )

  if a != "http://www.foo.com/bar" {
    t.Errorf( "%s should be %s", a, "http://www.foo.com/bar" )
  }

  a = helpers.AbsoluteURL( "//bar" )

  if a != "http://www.foo.com/bar" {
    t.Errorf( "%s should be %s", a, "http://www.foo.com/bar" )
  }

}