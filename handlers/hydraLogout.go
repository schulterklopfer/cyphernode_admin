package handlers

import (
  "github.com/gin-gonic/gin"
)

func HydraLogoutGet( c *gin.Context ) {

  /*
  challenge, _ := c.GetQuery( "logout_challenge" )

  getLogoutRequestParams := admin.NewGetLogoutRequestParams()
  getLogoutRequestParams.LogoutChallenge = challenge

  getLogoutRequestResponse, err := hydraAPI.GetBackendClient().Admin.GetLogoutRequest(getLogoutRequestParams)

  if err != nil {
    c.Header("X-Status-Reason", err.Error() )
    c.Status(http.StatusBadRequest )
    return
  }

  acceptLogoutRequestParams := admin.NewAcceptLogoutRequestParams()
  acceptLogoutRequestParams.LogoutChallenge = challenge
  */

}

func HydraLogoutPost( c *gin.Context ) {

}
