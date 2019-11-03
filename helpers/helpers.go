package helpers

import (
  "crypto/rand"
  "encoding/json"
  "github.com/schulterklopfer/cyphernode_admin/globals"
  "github.com/schulterklopfer/cyphernode_admin/password"
  "io"
  "os"
  "reflect"
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

func RandomString(length int, encodeToString func([]byte) string ) string {
  randomBytes := make([]byte, length)
  if _, err := io.ReadFull(rand.Reader, randomBytes); err != nil {
    return ""
  }
  return strings.TrimRight( encodeToString( randomBytes), "=" )
}

func AbsoluteURL( path string ) string {
  return AbsoluteURLFromHostEnvKey( globals.BASE_URL_EXTERNAL_ENV_KEY, path )
}

func AbsoluteURLFromHostEnvKey( hostEnvKEy string, path string ) string {
  return AbsoluteURLFromHost( GetenvOrDefault( hostEnvKEy ), path )
}

func AbsoluteURLFromHost( host string, path string ) string {
  for strings.HasSuffix( host,"/") {
    // remove last character
    host = host[:len(host)-1]
  }

  for strings.HasPrefix( path,"/") {
    // remove last character
    path = path[1:len(path)]
  }

  return host+"/"+path
}


func SetByJsonTag( obj interface{}, values *map[string]interface{} ) {

  // evaluate sbjt tag actions like hashing passwords
  structType := reflect.TypeOf(obj).Elem()
  //mutableObject := reflect.ValueOf(obj).Elem()
  for jsonFieldName, jsonFieldValue := range *values {
    for i := 0; i < structType.NumField(); i++ {
      field := structType.Field(i)
      jsonTag, hasJsonTag := field.Tag.Lookup("json")
      sbjtTag, hasSbjtTag := field.Tag.Lookup("sbjt")

      if hasSbjtTag && hasJsonTag && jsonTag == jsonFieldName {
        switch sbjtTag {
        case "hashPassword":
          if reflect.TypeOf(jsonFieldValue).Kind() == reflect.String {
            hashedPassword, _ := password.HashPassword(jsonFieldValue.(string))
            (*values)[jsonFieldName] = hashedPassword
          }
        }
      }
    }
  }

  jsonStringBytes, _ := json.Marshal( values )
  _ = json.Unmarshal( jsonStringBytes, obj )

}

func GetenvOrDefault( key string ) string {
  value := os.Getenv( key )
  if value == "" {
    defaultValue, _ := globals.DEFAULTS[key]
    return defaultValue
  }
  return value
}