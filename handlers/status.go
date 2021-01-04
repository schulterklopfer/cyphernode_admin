package handlers

import (
  "github.com/gin-gonic/gin"
  "github.com/schulterklopfer/cyphernode_admin/cyphernodeState"
  "net/http"
)

func GetStatus(c *gin.Context) {
  c.JSON(http.StatusOK, cyphernodeState.Instance())
}
