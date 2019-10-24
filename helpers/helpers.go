package helpers

import (
	"crypto/rand"
	"encoding/base32"
	"github.com/markbates/goth"
	"github.com/schulterklopfer/cyphernode_admin/globals"
	"io"
	"regexp"
	"strings"
	"time"
)

func SliceIndex(limit int, predicate func(i int) bool) int {
  for i := 0; i < limit; i++ {
    if predicate(i) {
      return i
    }
  }
  return -1
}

func SetInterval(someFunc func(), milliseconds int, async bool) chan bool {

  // How often to fire the passed in function
  // in milliseconds
  interval := time.Duration(milliseconds) * time.Millisecond

  // Setup the ticket and the channel to signal
  // the ending of the interval
  ticker := time.NewTicker(interval)
  clear := make(chan bool)

  // Put the selection in a go routine
  // so that the for loop is none blocking
  go func() {
    for {

      select {
      case <-ticker.C:
        if async {
          // This won't block
          go someFunc()
        } else {
          // This will block
          someFunc()
        }
      case <-clear:
        ticker.Stop()
        return
      }

    }
  }()

  // We return the channel so we can pass in
  // a value to it to clear the interval
  return clear

}

func EndpointIsPublic( endpoint string ) bool {
  for i:=0; i<len( globals.ENDPOINTS_PUBLIC_PATTERNS); i++ {
    pattern := globals.ENDPOINTS_PUBLIC_PATTERNS[i]
    matches, err := regexp.MatchString( pattern, endpoint )
    if matches && err == nil {
      return true
    }
  }
  return false
}

func UserIsAdmin( user *goth.User ) bool {
	if user == nil {
		return false
	}
	if val, ok := user.RawData["roles"]; ok {
		if roles, rolesAssertionOk := val.([]interface{}); rolesAssertionOk {
			for i:=0; i<len(roles); i++ {
				if role, roleAssertionOk := roles[i].(map[string]interface{}); roleAssertionOk {
					if roleName, ok := role["name"]; ok && roleName.(string) == globals.ROLES_ADMIN_ROLE {
						return true
					}
				}
			}
		}
	}
	return false
}

func RandomString(length int) string {
  randomBytes := make([]byte, length)
  if _, err := io.ReadFull(rand.Reader, randomBytes); err != nil {
    return ""
  }
  return strings.TrimRight( base32.StdEncoding.EncodeToString( randomBytes), "=" )
}

