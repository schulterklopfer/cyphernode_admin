package handlers

import (
  "github.com/gin-gonic/gin"
  "net/http"
)

func GetRolesForApp( c *gin.Context ) {

  // Resource endpoint. Need to check for tokens
  // and scopes here

  clientID := c.Params[0].Value
  username := c.Params[1].Value
  if clientID == "" || username == "" {
    c.Status( http.StatusNotFound )
    return
  }

  // TODO: check if "roles" scope is set
  // TODO: return roles of user with username in app with given hash

  //session := sessions.Default(c)
  //accessToken := session.Get( globals.HYDRA_ACCESS_TOKEN_SESSION_KEY )

  //introspectOAuth2TokenParams := admin.NewIntrospectOAuth2TokenParams()
  //introspectOAuth2TokenParams.Context = hydraAPI

  //hydraAPI.GetBackendClient().Admin.IntrospectOAuth2Token()

  println( clientID, username )
}
