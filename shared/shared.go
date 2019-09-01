package shared

import (
  "github.com/schulterklopfer/cyphernode_admin/password"
  "reflect"
)

func SetByJsonTag( obj interface{}, values *map[string]interface{} ) {
  structType := reflect.TypeOf(obj).Elem()
  mutableObject := reflect.ValueOf(obj).Elem()

  for jsonFieldName, jsonFieldValue := range *values {
    for i := 0; i < structType.NumField(); i++ {
      field := structType.Field(i)
      fieldValue := mutableObject.FieldByName(field.Name)
      jsonTag, hasJsonTag := field.Tag.Lookup("json")

      if hasJsonTag && jsonTag == jsonFieldName {
        if fieldValue.Type() == reflect.TypeOf(jsonFieldValue) {
          switch jsonFieldValue.(type) {
          case string:
            sbjtTag, hasSbjtTag := field.Tag.Lookup("sbjt")
            v := jsonFieldValue.(string)
            if hasSbjtTag && sbjtTag == "hashPassword" {
              v, _ = password.HashPassword( v )
            }
            fieldValue.SetString(v)
          case int64:
            fieldValue.SetInt(jsonFieldValue.(int64))
          case int32:
            fieldValue.SetInt(int64(jsonFieldValue.(int32)))
          case int:
            fieldValue.SetInt(int64(jsonFieldValue.(int)))
          case bool:
            fieldValue.SetBool(jsonFieldValue.(bool))
          case float64:
            fieldValue.SetFloat(jsonFieldValue.(float64))
          case float32:
            fieldValue.SetFloat(float64(jsonFieldValue.(float32)))
          }
        }
      }
    }
  }
}