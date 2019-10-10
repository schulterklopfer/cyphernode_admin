package handlers

import (
	"github.com/gin-gonic/gin"
	hydraAdmin "github.com/ory/hydra/sdk/go/hydra/client/admin"
	"github.com/ory/hydra/sdk/go/hydra/models"
	hydraModels "github.com/ory/hydra/sdk/go/hydra/models"
	"github.com/schulterklopfer/cyphernode_admin/globals"
	"github.com/schulterklopfer/cyphernode_admin/helpers"
	"github.com/schulterklopfer/cyphernode_admin/hydraAPI"
	"net/http"
)

func GetHydraConsent( c *gin.Context ) {
	challenge, _ := c.GetQuery( "consent_challenge" )

	if challenge == "" {
		c.Header("X-Status-Reason","no challenge" )
		c.Status(http.StatusBadRequest )
		return
	}

	consentRequestParams := hydraAdmin.NewGetConsentRequestParams()
	consentRequestParams.ConsentChallenge = challenge
	getConsentResponse, err := hydraAPI.GetBackendClient().Admin.GetConsentRequest(consentRequestParams)

	if err != nil {
		c.Header("X-Status-Reason", err.Error() )
		c.Status(http.StatusBadRequest )
		return
	}

	if getConsentResponse.GetPayload().Skip {
		// You can apply logic here, for example grant another scope, or do whatever...
		// ...

		acceptConsentRequestParams := hydraAdmin.NewAcceptConsentRequestParams()
		var handledConsentRequest hydraModels.HandledConsentRequest
		acceptConsentRequestParams.ConsentChallenge = challenge
		acceptConsentRequestParams.Body = &handledConsentRequest

		// Now it's time to grant the consent request. You could also deny the request if something went terribly wrong

		// We can grant all scopes that have been requested - hydra already checked for us that no additional scopes
		// are requested accidentally.
		acceptConsentRequestParams.Body.GrantedScope = getConsentResponse.GetPayload().RequestedScope

		// Grant the roles scope to every client, so they can get the users roles over the backchannel
		// ... is roles scope already in requestBody.GrantScope
		// the client app will be able to request /api/v0/roles. This endpoint will only return
		// data if a valid appHash is provided in the querystring. Returned information
		// is strictly limited to the app requesting the roles. You only receive your roles in the
		// requesting app

		if helpers.SliceIndex( len(acceptConsentRequestParams.Body.GrantedScope), func(i int) bool {
			return acceptConsentRequestParams.Body.GrantedScope[i] == globals.HYDRA_ROLES_SCOPE
		} ) == -1 {
			acceptConsentRequestParams.Body.GrantedScope = append(acceptConsentRequestParams.Body.GrantedScope, globals.HYDRA_ROLES_SCOPE )
		}


		// ORY Hydra checks if requested audiences are allowed by the client, so we can simply echo this.
		acceptConsentRequestParams.Body.GrantedAudience = getConsentResponse.GetPayload().RequestedAudience

		// This data will be available when introspecting the token. Try to avoid sensitive information here,
		// unless you limit who can introspect tokens.
		// access_token: { foo: 'bar' },

		// This data will be available in the ID token.
		// id_token: { baz: 'bar' },

		session := new(models.ConsentRequestSessionData)
		acceptConsentRequestParams.Body.Session = session

		acceptConsentResponse, err := hydraAPI.GetBackendClient().Admin.AcceptConsentRequest(acceptConsentRequestParams)

		if err != nil {
			c.Header("X-Status-Reason", err.Error() )
			c.Status(http.StatusBadRequest )
			return
		}

		// All we need to do now is to redirect the user back to hydra!
		c.Redirect(http.StatusTemporaryRedirect, acceptConsentResponse.GetPayload().RedirectTo )
	} else {
		// If consent can't be skipped we MUST show the consent UI.
		// TODO: render consent UI here

		// If consent can't be skipped we MUST show the consent UI.
		c.HTML(http.StatusOK, "hydra/consent", gin.H{
			"title": "consent",
			"csrfToken": "",
			"challenge": challenge,
			"requestedScope": getConsentResponse.GetPayload().RequestedScope,
			"user": getConsentResponse.GetPayload().Subject,
			"client": getConsentResponse.GetPayload().Client,
		})

	}
}

func PostHydraConsent( c *gin.Context ) {
	challenge, _ := c.GetPostForm( "challenge" )

	if challenge == "" {
		c.Header("X-Status-Reason","no challenge" )
		c.Status(http.StatusBadRequest )
		return
	}

	submitValue, _ := c.GetPostForm( "submit" )

	if submitValue == "Deny access" {
		// Looks like the consent request was denied by the user

		rejectConsentRequestParams := hydraAdmin.NewRejectConsentRequestParams()
		var requestBody models.RequestDeniedError
		rejectConsentRequestParams.ConsentChallenge = challenge
		rejectConsentRequestParams.Body = &requestBody
		rejectConsentRequestParams.Body.Name = "access_denied"
		rejectConsentRequestParams.Body.Description = "The resource owner denied the request"

		rejectConsentResponse, err := hydraAPI.GetBackendClient().Admin.RejectConsentRequest(rejectConsentRequestParams)

		if err != nil {
			c.Header("X-Status-Reason", err.Error() )
			c.Status(http.StatusBadRequest )
			return
		}

		// All we need to do now is to redirect the browser back to hydra!
		c.Redirect(http.StatusTemporaryRedirect, rejectConsentResponse.GetPayload().RedirectTo )
	} else {
		grantScope, _ := c.GetPostFormArray("grant_scope" )
		rememberValue, _ := c.GetPostForm("remember" )
		remember := rememberValue == "1"

		// Seems like the user authenticated! Let's tell hydra...
		getConsentRequestParams := hydraAdmin.NewGetConsentRequestParams()
		getConsentRequestParams.ConsentChallenge = challenge
		getConsentResponse, err := hydraAPI.GetBackendClient().Admin.GetConsentRequest(getConsentRequestParams)

		if err != nil {
			c.Header("X-Status-Reason", err.Error() )
			c.Status(http.StatusBadRequest )
			return
		}

		// You can apply logic here, for example grant another scope, or do whatever...
		// ...

		// Now it's time to grant the consent request. You could also deny the request if something went terribly wrong

		acceptConsentRequestParams := hydraAdmin.NewAcceptConsentRequestParams()
		acceptConsentRequestParams.ConsentChallenge = challenge
		var handledConsentRequest hydraModels.HandledConsentRequest
		acceptConsentRequestParams.Body = &handledConsentRequest

		// We can grant all scopes that have been requested - hydra already checked for us that no additional scopes
		// are requested accidentally.
		acceptConsentRequestParams.Body.GrantedScope = grantScope

		// Grant the roles scope to every client, so they can get the users roles over the backchannel
		// ... is roles scope already in requestBody.GrantScope
		// the client app will be able to request /api/v0/roles. This endpoint will only return
		// data if a valid appHash is provided in the querystring. Returned information
		// is strictly limited to the app requesting the roles. You only receive your roles in the
		// requesting app

		if helpers.SliceIndex( len(acceptConsentRequestParams.Body.GrantedScope), func(i int) bool {
			return acceptConsentRequestParams.Body.GrantedScope[i] == globals.HYDRA_ROLES_SCOPE
		} ) == -1 {
			acceptConsentRequestParams.Body.GrantedScope = append(acceptConsentRequestParams.Body.GrantedScope, globals.HYDRA_ROLES_SCOPE )
		}

		// ORY Hydra checks if requested audiences are allowed by the client, so we can simply echo this.
		acceptConsentRequestParams.Body.GrantedAudience = getConsentResponse.GetPayload().RequestedAudience

		// This tells hydra to remember this consent request and allow the same client to request the same
		// scopes from the same user, without showing the UI, in the future.
		acceptConsentRequestParams.Body.Remember = remember

		// When this "remember" session expires, in seconds. Set this to 0 so it will never expire.
		acceptConsentRequestParams.Body.RememberFor = 3600

		// This data will be available when introspecting the token. Try to avoid sensitive information here,
		// unless you limit who can introspect tokens.
		// access_token: { foo: 'bar' },

		// This data will be available in the ID token.
		// id_token: { baz: 'bar' },
		//session := new(models.ConsentRequestSessionData)
		//acceptConsentRequestParams.Body.Session = session

		acceptConsentResponse, err := hydraAPI.GetBackendClient().Admin.AcceptConsentRequest(acceptConsentRequestParams)

		if err != nil {
			c.Header("X-Status-Reason", err.Error() )
			c.Status(http.StatusBadRequest )
			return
		}

		// All we need to do now is to redirect the user back to hydra!
		c.Redirect(http.StatusTemporaryRedirect, acceptConsentResponse.GetPayload().RedirectTo )

	}

}


