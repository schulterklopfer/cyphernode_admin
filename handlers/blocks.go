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
