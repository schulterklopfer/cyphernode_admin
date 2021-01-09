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
  "github.com/gorilla/websocket"
  "github.com/schulterklopfer/cyphernode_admin/dockerApi"
  "log"
  "net/http"
  "time"
)

func GetContainerByBas64Image(c *gin.Context) {
  // param 0 is first param in url pattern
  hashedImage := c.Params[0].Value

  container := dockerApi.Instance().ContainerByBase64Image(hashedImage)

  if container == nil {
    c.Status(http.StatusNotFound )
    return
  }

  c.JSON(http.StatusOK, container )
}


func GetContainerByName(c *gin.Context) {
  // param 0 is first param in url pattern
  name := c.Params[0].Value

  container := dockerApi.Instance().ContainerByNameRegexp( name )

  if container == nil {
    c.Status(http.StatusNotFound )
    return
  }

  c.JSON(http.StatusOK, container )
}

var upGrader = websocket.Upgrader{
  CheckOrigin: func(r *http.Request) bool {
    return true
  },
}

//webSocket returns json format
func WSLogsByContainerId(c *gin.Context) {
  if len(c.Params) < 1 {
    c.Status( http.StatusNotFound )
    return
  }

  containerId := c.Params[0].Value

  if !dockerApi.Instance().ContainerExists( containerId ) {
    c.Status( http.StatusNotFound )
    return
  }

  //Upgrade get request to webSocket protocol
  ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
  if err != nil {
    log.Println("error get connection")
    c.Status( http.StatusNotFound )
    return
  }
  defer ws.Close()

  sending := true

  go func() {
    lastSeenLineIndex := -1
    for sending {
      lines, lastLineIndex := dockerApi.Instance().LogLinesById( containerId )

      if lastLineIndex != lastSeenLineIndex {
        lineDiff := lastLineIndex - lastSeenLineIndex
        lineCount := len(lines)

        // lastLineIndex may be larger than lineCount cause we only store DOCKER_LOGS_MAX_LINES lines
        if lineDiff > lineCount {
          lineDiff = lineCount
        }

        lastSeenLineIndex = lastLineIndex
        err := ws.WriteJSON( lines[ lineCount-lineDiff: ] )

        if err != nil {
          sending = false
        }
      }

      time.Sleep( 250 * time.Millisecond )
    }
  }()

  // listening for close message here
  for sending {
    _, _, err := ws.ReadMessage()
    if err != nil {
      break
    }
  }
  sending = false
}