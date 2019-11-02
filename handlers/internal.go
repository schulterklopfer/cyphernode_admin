package handlers

import (
  "github.com/gin-gonic/gin"
  "net/http"
)

func InternalRegisterApp( c *gin.Context ) {

  // get client id and check whitelist for id.
  // if client id is in whiteliste, register the app
  // if not, deny registration

  c.JSON( http.StatusOK, map[string]string{ "registered": "true" })
}
