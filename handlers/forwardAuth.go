package handlers

import (
  "fmt"
  "github.com/dgrijalva/jwt-go"
  "github.com/gin-gonic/gin"
  "github.com/schulterklopfer/cyphernode_admin/globals"
  "github.com/schulterklopfer/cyphernode_admin/helpers"
  "github.com/schulterklopfer/cyphernode_admin/models"
  "github.com/schulterklopfer/cyphernode_admin/queries"
  "net/http"
  "strings"
)

func tokenFromBearerAuthHeader(authHeader string) string {
  if authHeader == "" {
    return ""
  }

  parts := strings.Split(authHeader, "Bearer")
  if len(parts) != 2 {
    return ""
  }

  token := strings.TrimSpace(parts[1])
  if len(token) < 1 {
    return ""
  }

  return token
}

func ForwardAuthAuth(c *gin.Context) {

  //secret := []byte("my_secret_key")

  // x-forwarded-prefix header idetentifies the app we want to auth against
  mountPoint := c.Request.Header.Get("x-forwarded-prefix")[1:]

  app, err := queries.GetAppByMountPoint( mountPoint )

  if err != nil {
    c.Header("X-Status-Reason", err.Error() )
    c.Status(http.StatusInternalServerError)
    return
  }

  uriInAp := c.Request.Header.Get("x-forwarded-uri")
  method := c.Request.Header.Get("x-forwarded-method")
  // check for public access
  accessGranted := false
  for _, accessPolicy := range app.AccessPolicies {
    accessGranted = accessPolicy.Check( method, uriInAp, nil )
    if accessGranted {
      break
    }
  }

  if accessGranted {
    c.Status(http.StatusOK)
    return
  }

  // access not granted. See if we have a valid token
  // and check access again


  tokenString := tokenFromBearerAuthHeader( c.Request.Header.Get("authorization") )

  // Parse takes the token string and a function for looking up the key. The latter is especially
  // useful if you use multiple keys for your application.  The standard is to use 'kid' in the
  // head of the token to identify which key to use, but the parsed token (head and claims) is provided
  // to the callback, providing flexibility.
  token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
    // Don't forget to validate the alg is what you expect:
    if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
      return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
    }
    return []byte(helpers.GetenvOrDefault(globals.CNA_COOKIE_SECRET_ENV_KEY)), nil
  })

  if err != nil {
    c.Header("X-Status-Reason", err.Error() )
    c.Status(http.StatusUnauthorized)
    return
  }

  if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

    subject, exists := claims["sub"]

    if !exists {
      c.Header("X-Status-Reason", "no subject claims" )
      c.Status(http.StatusUnauthorized)
      return
    }

    parts := strings.Split( token.Raw, "." )
    c.Header("X-Auth-User-Claims", parts[1] )

    userId := uint(subject.(float64))
    var user models.UserModel
    err := queries.Get( &user, userId,true )

    if err != nil || user.ID == 0 {
      c.Status(http.StatusUnauthorized)
      return
    }

    var roleNames []string;
    for _, role := range user.Roles {
      if role.AppId != app.ID {
        continue
      }
      roleNames = append( roleNames, role.Name )
    }

    accessGranted := false
    for _, accessPolicy := range app.AccessPolicies {
      accessGranted = accessPolicy.Check( method, uriInAp, roleNames )
      if accessGranted {
        break
      }
    }

    if accessGranted {
      c.Status(http.StatusOK)
      return
    }

  }

  c.Status(http.StatusUnauthorized)

}
