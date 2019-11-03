package handlers

import (
  "github.com/gin-gonic/gin"
  "github.com/schulterklopfer/cyphernode_admin/appWhitelist"
  "github.com/schulterklopfer/cyphernode_admin/models"
  "github.com/schulterklopfer/cyphernode_admin/queries"
  "gopkg.in/validator.v2"
  "net/http"
)

func InternalRegisterApp( c *gin.Context ) {

  // get client id and check whitelist for id.
  // if client id is in whiteliste, register the app
  // if not, deny registration


  var app models.AppModel
  err := c.Bind( &app )


  if err != nil {
    c.Header("X-Status-Reason", err.Error() )
    c.Status(http.StatusBadRequest)
    return
  }

  err = validator.Validate(app)

  if err != nil {
    c.Header("X-Status-Reason", err.Error() )
    c.Status(http.StatusBadRequest)
    return
  }

  if appWhitelist.Get().ContainsClientID( app.ClientID ) {

    var existingApps []*models.AppModel
    err := queries.Find( &existingApps,  []interface{}{"client_id = ?", app.ClientID }, "", 1,0,false)
    if err != nil {
      c.Header("X-Status-Reason", err.Error() )
      c.Status(http.StatusBadRequest)
      return
    }

    if len(existingApps) > 0 {
      c.Header("X-Status-Reason", "app already registered" )
      c.Status(http.StatusBadRequest)
      return
    }

    err = queries.Create( &app )
    if err != nil {
      c.Header("X-Status-Reason", err.Error() )
      c.Status(http.StatusBadRequest)
      return
    }
    c.Status( http.StatusCreated )
  } else {
    c.Header("X-Status-Reason","not in whitelist" )
    c.Status(http.StatusBadRequest)
    return
  }

}
