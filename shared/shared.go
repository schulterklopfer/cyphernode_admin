package shared

import (
  "encoding/json"
  "github.com/schulterklopfer/cyphernode_admin/password"
  "reflect"
)

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