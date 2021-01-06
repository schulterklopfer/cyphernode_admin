package handlers

import (
  "github.com/gin-gonic/gin"
  "github.com/schulterklopfer/cyphernode_admin/dockerApi"
  "net/http"
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
