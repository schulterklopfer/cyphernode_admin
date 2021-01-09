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
  "github.com/schulterklopfer/cyphernode_admin/cyphernodeState"
  "github.com/schulterklopfer/cyphernode_admin/globals"
  "net/http"
  "strconv"
)

func GetLatestBlocks(c *gin.Context) {

  count := globals.LATEST_BLOCK_COUNT

  if c, exists := c.GetQuery("c"); exists {
    i, err := strconv.Atoi(c)
    if err == nil {
      count = i
    }
  }

  if count < 1 {
    count = 1
  }

  if count > globals.LATEST_BLOCK_COUNT {
    count = globals.LATEST_BLOCK_COUNT
  }

  if len(cyphernodeState.Instance().LatestBlocks) < globals.LATEST_BLOCK_COUNT {
    c.Status(http.StatusNotFound)
    return
  }

  c.JSON(http.StatusOK, cyphernodeState.Instance().LatestBlocks[0:count] )
}
