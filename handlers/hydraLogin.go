package handlers

import (
	"github.com/gin-gonic/gin"
	hydraAdmin "github.com/ory/hydra/sdk/go/hydra/client/admin"
	hydraModels "github.com/ory/hydra/sdk/go/hydra/models"
	"github.com/schulterklopfer/cyphernode_admin/authentication"
	"github.com/schulterklopfer/cyphernode_admin/hydraAPI"
	"net/http"
)

func GetHydraLogin( c *gin.Context ) {
	challenge, _ := c.GetQuery( "login_challenge" )

	if challenge == "" {
		// no challenge ... bad
		return
	}
	getLoginRequestParams := hydraAdmin.NewGetLoginRequestParams()
	getLoginRequestParams.LoginChallenge = challenge

	getLoginResponse, err := hydraAPI.GetBackendClient().Admin.GetLoginRequest(getLoginRequestParams)

	if err != nil {
		// err ... bad
		println( err )
		return
	}

	if getLoginResponse.GetPayload().Skip {
		// You can apply logic here, for example grant another scope, or do whatever...
		// ...

		// Now it's time to grant the login request. You could also deny the request if something went terribly wrong
		acceptLoginRequestParams := hydraAdmin.NewAcceptLoginRequestParams()
		var handledLoginRequest hydraModels.HandledLoginRequest
		acceptLoginRequestParams.LoginChallenge = challenge
		acceptLoginRequestParams.Body = &handledLoginRequest
		acceptLoginRequestParams.Body.Subject = &getLoginResponse.GetPayload().Subject

		acceptLoginResponse, err := hydraAPI.GetBackendClient().Admin.AcceptLoginRequest(acceptLoginRequestParams)

		if err != nil {
			// something is wrong
			println( err )
			return
		}

		// All we need to do now is to redirect the user back to hydra!
		c.Redirect(http.StatusTemporaryRedirect, acceptLoginResponse.GetPayload().RedirectTo )
	} else {
		// If consent can't be skipped we MUST show the consent UI.
		// TODO: render login UI here
		c.HTML(http.StatusOK, "hydra/login", gin.H{
			"title": "login",
			"csrfToken": "",
			"challenge": challenge,
		})
	}
}

func PostHydraLogin( c *gin.Context ) {
	challenge, _ := c.GetPostForm( "challenge" )

	if challenge == "" {
		// no challenge ... bad
		return
	}

	login, _ := c.GetPostForm( "login" )
	password, _ := c.GetPostForm( "password" )
	rememberValue, _ := c.GetPostForm("remember" )
	remember := rememberValue == "1"

	err := authentication.CheckUserPassword( login, password )

	if err != nil {
		// stuff is wrong
		// show login ui again
		c.HTML(http.StatusOK, "hydra/login", gin.H{
			"title": "login",
			"csrfToken": "",
			"challenge": challenge,
			"error": err.Error(),
		})
		return
	}

	acceptLoginRequestParams := hydraAdmin.NewAcceptLoginRequestParams()
	var handledLoginRequest hydraModels.HandledLoginRequest
	acceptLoginRequestParams.Body = &handledLoginRequest
	acceptLoginRequestParams.LoginChallenge = challenge

	// Subject is an alias for user ID. A subject can be a random string, a UUID, an email address, ....
	acceptLoginRequestParams.Body.Subject = &login

	// This tells hydra to remember the browser and automatically authenticate the user in future requests. This will
	// set the "skip" parameter in the other route to true on subsequent requests!
	acceptLoginRequestParams.Body.Remember = remember

	// When the session expires, in seconds. Set this to 0 so it will never expire.
	acceptLoginRequestParams.Body.RememberFor = 3600

	// Sets which "level" (e.g. 2-factor authentication) of authentication the user has. The value is really arbitrary
	// and optional. In the context of OpenID Connect, a value of 0 indicates the lowest authorization level.
	acceptLoginRequestParams.Body.ACR = "0"

	// Seems like the user authenticated! Let's tell hydra...
	acceptLoginResponse, err := hydraAPI.GetBackendClient().Admin.AcceptLoginRequest(acceptLoginRequestParams)

	if err != nil {
		// something is wrong
		println( err )
		return
	}

	// All we need to do now is to redirect the browser back to hydra!
	c.Redirect(http.StatusTemporaryRedirect, acceptLoginResponse.Payload.RedirectTo )
}

// You could also deny the login request which tells hydra that no one authenticated!
// hydra.rejectLoginRequest(challenge, {
//   error: 'invalid_request',
//   error_description: 'The user did something stupid...'
// })
//   .then(function (response) {
//     // All we need to do now is to redirect the browser back to hydra!
//     res.redirect(response.redirect_to);
//   })
//   // This will handle any error that happens when making HTTP calls to hydra
//   .catch(function (error) {
//     next(error);
//   });