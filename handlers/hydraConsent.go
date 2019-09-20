package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/schulterklopfer/cyphernode_admin/hydra"
	"net/http"
)

func GetHydraConsent( c *gin.Context ) {
	challenge, _ := c.GetQuery( "consent_challenge" )

	if challenge == "" {
		// no challenge ... bad
		return
	}

	getConsentResponse, err := hydra.GetConsentRequest( http.DefaultClient, challenge )

	if err == nil {
		// err ... bad
		return
	}

	if getConsentResponse.Skip {
		// You can apply logic here, for example grant another scope, or do whatever...
		// ...

		// Now it's time to grant the consent request. You could also deny the request if something went terribly wrong
		requestBody := new( hydra.RequestBody )

		// We can grant all scopes that have been requested - hydra already checked for us that no additional scopes
		// are requested accidentally.
		requestBody.GrantScope = getConsentResponse.RequestedScope

		// ORY Hydra checks if requested audiences are allowed by the client, so we can simply echo this.
		requestBody.GrantAccessTokenAudience = getConsentResponse.RequestedAccessTokenAudience

		// This data will be available when introspecting the token. Try to avoid sensitive information here,
		// unless you limit who can introspect tokens.
		// access_token: { foo: 'bar' },

		// This data will be available in the ID token.
		// id_token: { baz: 'bar' },
		requestBody.Session = new(hydra.Session)

		acceptConsentResponse, err := hydra.AcceptConsentRequest( http.DefaultClient, challenge, requestBody )

		if err != nil {
			// something is wrong
			return
		}

		// All we need to do now is to redirect the user back to hydra!
		c.Redirect(http.StatusTemporaryRedirect, acceptConsentResponse.RedirectTo )
	} else {
		// If consent can't be skipped we MUST show the consent UI.
		// TODO: render consent UI here

		// If consent can't be skipped we MUST show the consent UI.
		c.HTML(http.StatusOK, "hydra/consent", gin.H{
			"title": "consent",
			"csrfToken": "",
			"challenge": challenge,
			"requestedScope": getConsentResponse.RequestedScope,
			"user": getConsentResponse.Subject,
			"client": getConsentResponse.Client,
		})

	}
}

func PostHydraConsent( c *gin.Context ) {
	challenge, _ := c.GetPostForm( "challenge" )

	if challenge == "" {
		// no challenge ... bad
		return
	}

	submitValue, _ := c.GetPostForm( "submit" )

	if submitValue == "Deny access" {
		// Looks like the consent request was denied by the user
		requestBody := new( hydra.RequestBody )
		requestBody.Error = "access_denied"
		requestBody.ErrorDescription = "The resource owner denied the request"
		rejectConsentResponse, err := hydra.RejectConsentRequest( http.DefaultClient, challenge, requestBody )

		if err != nil {
			// something is wrong
			return
		}

		// All we need to do now is to redirect the browser back to hydra!
		c.Redirect(http.StatusTemporaryRedirect, rejectConsentResponse.RedirectTo )
	} else {
		grantScope, _ := c.GetPostFormArray("grant_scope" )
		rememberValue, _ := c.GetPostForm("remember" )
		remember := rememberValue == "true"

		// Seems like the user authenticated! Let's tell hydra...
		getConsentResponse, err := hydra.GetConsentRequest( http.DefaultClient, challenge )

		if err != nil {
			// something is wrong
			return
		}

		// You can apply logic here, for example grant another scope, or do whatever...
		// ...

		// Now it's time to grant the consent request. You could also deny the request if something went terribly wrong
		requestBody := new( hydra.RequestBody )

		// We can grant all scopes that have been requested - hydra already checked for us that no additional scopes
		// are requested accidentally.
		requestBody.GrantScope = grantScope

		// ORY Hydra checks if requested audiences are allowed by the client, so we can simply echo this.
		requestBody.GrantAccessTokenAudience = getConsentResponse.RequestedAccessTokenAudience

		// This tells hydra to remember this consent request and allow the same client to request the same
		// scopes from the same user, without showing the UI, in the future.
		requestBody.Remember = remember

		// When this "remember" session expires, in seconds. Set this to 0 so it will never expire.
		requestBody.RememberFor = 3600

		// This data will be available when introspecting the token. Try to avoid sensitive information here,
		// unless you limit who can introspect tokens.
		// access_token: { foo: 'bar' },

		// This data will be available in the ID token.
		// id_token: { baz: 'bar' },
		requestBody.Session = new(hydra.Session)

		acceptConsentResponse, err := hydra.AcceptConsentRequest( http.DefaultClient, challenge, requestBody )

		if err != nil {
			// something is wrong
			return
		}

		// All we need to do now is to redirect the user back to hydra!
		c.Redirect(http.StatusTemporaryRedirect, acceptConsentResponse.RedirectTo )

	}

}


