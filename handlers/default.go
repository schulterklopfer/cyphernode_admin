package handlers

import (
  "github.com/dgrijalva/jwt-go"
  "github.com/gin-gonic/gin"
  "github.com/schulterklopfer/cyphernode_admin/authentication"
  "github.com/schulterklopfer/cyphernode_admin/globals"
  "github.com/schulterklopfer/cyphernode_admin/helpers"
  "github.com/schulterklopfer/cyphernode_admin/queries"
  "net/http"
  "time"
)

func DefaultRoot( c *gin.Context ) {
  //if user, err := cnaOIDC.GetUser( c.Writer, c.Request ); err == nil {
  //  c.Set("user", user )
  //  c.Redirect( http.StatusTemporaryRedirect, globals.BASE_ENDPOINT_PRIVATE+globals.PRIVATE_ENDPOINTS_HOME)
  //} else {
    c.Redirect( http.StatusTemporaryRedirect, globals.BASE_ENDPOINT_PUBLIC+globals.PUBLIC_ENDPOINTS_LOGIN)
  //}
}

func DefaultLogin( c *gin.Context ) {
  // create token and session
  var input map[string]string

  err := c.Bind( &input )
  if err != nil {
    c.Header("X-Status-Reason", err.Error() )
    c.Status(http.StatusBadRequest )
    return
  }

  username, exists := input["username"]

  if !exists {
    c.Header("X-Status-Reason","username missing" )
    c.Status(http.StatusBadRequest )
    return
  }

  password, exists := input["password"]

  if !exists {
    c.Header("X-Status-Reason", "password missing" )
    c.Status(http.StatusBadRequest )
    return
  }

  user, err := authentication.CheckUserPassword( username, password )

  if err != nil {
    c.Header("X-Status-Reason", err.Error() )
    c.Status(http.StatusUnauthorized)
    return
  }

  roles, err := queries.GetRolesOfUser( user )

  if err != nil {
    c.Header("X-Status-Reason", err.Error() )
    c.Status(http.StatusBadRequest )
    return
  }

  var roleStrings []string

  for i:=0; i<len(roles); i++ {
    roleStrings = append(roleStrings, roles[i].Name )
  }

  nowUnix := time.Now().Unix()

  token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
    "sub": user.ID,
    "roles": roleStrings,
    "exp": nowUnix + 1*24*60*60, // expires in one day
    "auth_time": nowUnix,
  })

  tokenString, err := token.SignedString([]byte(helpers.GetenvOrDefault(globals.CNA_COOKIE_SECRET_ENV_KEY)))

  if err != nil {
    c.Header("X-Status-Reason", err.Error() )
    c.Status(http.StatusBadRequest )
    return
  }

  result := map[string]interface{} {
    "token": tokenString,
  }

  c.JSON(200, result )

}

func DefaultLogout( c *gin.Context ) {
  if _, exists := c.Get("user"); exists {
    //err := cnaOIDC.Logout( c.Writer, c.Request, helpers.AbsoluteURL(globals.URLS_BYEBYE) )
    //if err != nil {
    //  c.Header("X-Status-Reason", err.Error() )
    //  c.Status(http.StatusBadRequest )
    //  return
    //}
    //c.Redirect( http.StatusTemporaryRedirect, globals.BASE_ENDPOINT_PUBLIC+globals.PUBLIC_ENDPOINTS_BYEBYE)
    //return
  }
  c.Status( http.StatusNotFound )
}

func DefaultCallback( c *gin.Context ) {

  //gothUser, err := cnaOIDC.CompleteUserAuth( c.Writer, c.Request )

  //if err != nil {
  //  c.Header("X-Status-Reason", err.Error() )
  //  c.Status(http.StatusBadRequest )
  //  return
  //}

  //c.Set("user", gothUser )
  //c.Redirect( http.StatusTemporaryRedirect, globals.BASE_ENDPOINT_PRIVATE+globals.PRIVATE_ENDPOINTS_HOME)

}

func DefaultHome( c *gin.Context ) {
  //user, err := cnaOIDC.GetUser(c.Writer, c.Request)
  //if err == nil {
  //  c.JSON( http.StatusOK, &user )
  //}
}

func DefaultByeBye( c *gin.Context ) {
//  c.JSON( http.StatusOK, map[string]string{ "message": "bye bye!" })
}


/* https://godoc.org/golang.org/x/oauth2
ctx := context.Background()
conf := &oauth2.Config{
    Hash:     "YOUR_CLIENT_ID",
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
