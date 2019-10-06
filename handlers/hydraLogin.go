package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/schulterklopfer/cyphernode_admin/hydra"
	"net/http"
)

func GetHydraLogin( c *gin.Context ) {
	challenge, _ := c.GetQuery( "login_challenge" )

	if challenge == "" {
		// no challenge ... bad
		return
	}

	getLoginResponse, err := hydra.GetLoginRequest( http.DefaultClient, challenge )

	if err != nil {
		// err ... bad
		println( err )
		return
	}

	if getLoginResponse.Skip {
		// You can apply logic here, for example grant another scope, or do whatever...
		// ...

		// Now it's time to grant the login request. You could also deny the request if something went terribly wrong
		requestBody := new( hydra.RequestBody )
		requestBody.Subject = getLoginResponse.Subject

		acceptLoginResponse, err := hydra.AcceptLoginRequest( http.DefaultClient, challenge, requestBody )

		if err != nil {
			// something is wrong
			println( err )
			return
		}

		// All we need to do now is to redirect the user back to hydra!
		c.Redirect(http.StatusTemporaryRedirect, acceptLoginResponse.RedirectTo )
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

	email, _ := c.GetPostForm( "email" )
	password, _ := c.GetPostForm( "password" )
	rememberValue, _ := c.GetPostForm("remember" )
	remember := rememberValue == "1"

	if email != "foo@bar.com" || password != "foobar" {
		// stuff is wrong
		// show login ui again
		c.HTML(http.StatusOK, "hydra/login", gin.H{
			"title": "login",
			"csrfToken": "",
			"challenge": challenge,
			"error": "user or password wrong",
		})
		return
	}

	requestBody := new( hydra.RequestBody )

	// Subject is an alias for user ID. A subject can be a random string, a UUID, an email address, ....
	requestBody.Subject = "foo@bar.com"

	// This tells hydra to remember the browser and automatically authenticate the user in future requests. This will
	// set the "skip" parameter in the other route to true on subsequent requests!
	requestBody.Remember = remember

	// When the session expires, in seconds. Set this to 0 so it will never expire.
	requestBody.RememberFor = 3600

	// Sets which "level" (e.g. 2-factor authentication) of authentication the user has. The value is really arbitrary
	// and optional. In the context of OpenID Connect, a value of 0 indicates the lowest authorization level.
	requestBody.Acr = 0

	// Seems like the user authenticated! Let's tell hydra...
	acceptLoginResponse, err := hydra.AcceptLoginRequest( http.DefaultClient, challenge, requestBody )

	if err != nil {
		// something is wrong
		println( err )
		return
	}

	// All we need to do now is to redirect the browser back to hydra!
	c.Redirect(http.StatusTemporaryRedirect, acceptLoginResponse.RedirectTo )
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