/*
 * MIT License
 *
 * Copyright (c) 2021 schulterklopfer/__escapee__
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILIT * Y, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package handlers

import (
  "github.com/dgrijalva/jwt-go"
  "github.com/gin-gonic/gin"
  "github.com/schulterklopfer/cyphernode_admin/authentication"
  "github.com/schulterklopfer/cyphernode_admin/globals"
  "github.com/schulterklopfer/cyphernode_fauth/helpers"
  "net/http"
  "time"
)

func Login( c *gin.Context ) {
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

  nowUnix := time.Now().Unix()

  token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
    "id": user.ID,
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
