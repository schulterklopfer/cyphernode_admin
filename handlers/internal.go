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
  "github.com/gin-gonic/gin"
  "github.com/schulterklopfer/cyphernode_fauth/models"
  "gopkg.in/validator.v2"
  "net/http"
)

func InternalRegisterApp( c *gin.Context ) {

  // get client id and check whitelist for id.
  // if client id is in whiteliste, register the app
  // if not, deny registration


  var app models.AppModel
  err := c.Bind( &app )


  if err != nil {
    c.Header("X-Status-Reason", err.Error() )
    c.Status(http.StatusBadRequest)
    return
  }

  err = validator.Validate(app)

  if err != nil {
    c.Header("X-Status-Reason", err.Error() )
    c.Status(http.StatusBadRequest)
    return
  }

  /*
  if appList.Get().ContainsClientID( app.Hash ) {

    var existingApps []*models.AppModel
    err := queries.Find( &existingApps,  []interface{}{"client_id = ?", app.Hash }, "", 1,0,false)
    if err != nil {
      c.Header("X-Status-Reason", err.Error() )
      c.Status(http.StatusBadRequest)
      return
    }

    if len(existingApps) > 0 {
      c.Header("X-Status-Reason", "app already registered" )
      c.Status(http.StatusBadRequest)
      return
    }

    err = queries.Create( &app )
    if err != nil {
      c.Header("X-Status-Reason", err.Error() )
      c.Status(http.StatusBadRequest)
      return
    }
    c.Status( http.StatusCreated )
  } else {
    c.Header("X-Status-Reason","client id not in whitelist" )
    c.Status(http.StatusBadRequest)
    return
  }

   */

}
