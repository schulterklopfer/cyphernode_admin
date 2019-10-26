package handlers

import (
  "github.com/gin-gonic/gin"
  "github.com/schulterklopfer/cyphernode_admin/globals"
  "github.com/schulterklopfer/cyphernode_admin/oidc"
  "net/http"
)

func DefaultRoot( c *gin.Context ) {
  if _, exists :=  c.Get( "user"); exists {
    c.Redirect( http.StatusTemporaryRedirect, globals.ROUTER_GROUPS_BASE_ENDPOINT_PRIVATE+globals.PRIVATE_ENDPOINTS_HOME)
  } else {
    c.Redirect( http.StatusTemporaryRedirect, globals.ROUTER_GROUPS_BASE_ENDPOINT_PUBLIC+globals.PUBLIC_ENDPOINTS_LOGIN)
  }
}

func DefaultLogin( c *gin.Context ) {
  //url := oauth2.Get().Config.AuthCodeURL(helpers.RandomString(16 ), xoauth2.AccessTypeOffline)
  //fmt.Printf("redirecting to: %v", url)
  //c.Redirect( http.StatusTemporaryRedirect, url )

  //if gothUser, err := oidc.CompleteUserAuth(c.Writer, c.Request); err == nil {
  //  println( &gothUser )
  //} else {
  oidc.BeginAuthHandler(c.Writer, c.Request)
  //}

}

func DefaultCallbackGet( c *gin.Context ) {

  gothUser, err := oidc.CompleteUserAuth( c.Writer, c.Request )

  if err != nil {
    c.Header("X-Status-Reason", err.Error() )
    c.Status(http.StatusBadRequest )
    return
  }

  c.Set("user", gothUser )
  c.Redirect( http.StatusTemporaryRedirect, globals.ROUTER_GROUPS_BASE_ENDPOINT_PRIVATE+globals.PRIVATE_ENDPOINTS_HOME)

}

func DefaultHome( c *gin.Context ) {
  if user, exists :=  c.Get( "user"); exists {
    c.JSON( http.StatusOK, &user )
  }
}

func DefaultLogout( c *gin.Context ) {
  print( "logout")
}


/* https://godoc.org/golang.org/x/oauth2
ctx := context.Background()
conf := &oauth2.Config{
    ClientID:     "YOUR_CLIENT_ID",
    ClientSecret: "YOUR_CLIENT_SECRET",
    Scopes:       []string{"SCOPE1", "SCOPE2"},
    Endpoint: oauth2.Endpoint{
        AuthURL:  "https://provider.com/o/oauth2/auth",
        TokenURL: "https://provider.com/o/oauth2/token",
    },
}

// Redirect user to consent page to ask for permission
// for the scopes specified above.
url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
fmt.Printf("Visit the URL for the auth dialog: %v", url)

// Use the authorization code that is pushed to the redirect
// URL. Exchange will do the handshake to retrieve the
// initial access token. The HTTP Client returned by
// conf.Client will refresh the token as necessary.
var code string
if _, err := fmt.Scan(&code); err != nil {
    log.Fatal(err)
}
tok, err := conf.Exchange(ctx, code)
if err != nil {
    log.Fatal(err)
}

client := conf.Client(ctx, tok)
client.Get("...")
*/
