package dockerApi

import (
  "context"
  "encoding/base64"
  "github.com/docker/docker/api/types"
  "github.com/docker/docker/client"
  "github.com/schulterklopfer/cyphernode_admin/globals"
  "github.com/schulterklopfer/cyphernode_admin/helpers"
  "github.com/schulterklopfer/cyphernode_admin/logwrapper"
  "regexp"
  "strings"
  "sync"
)

type DockerApiConfig struct {
  Host string
  Port int
}

type DockerApi struct {
  DockerClient           *client.Client
  containerList          []types.Container
  containerById          map[string]int
  containerByBase64Image map[string]int
}


var instance *DockerApi
var once sync.Once

func initOnce() error {
  var initOnceErr error
  once.Do(func() {
    cli, err := client.NewClientWithOpts(client.FromEnv)
    if err != nil {
      initOnceErr = err
      return
    }

    instance = &DockerApi{
      DockerClient: cli,
    }

    err = instance.Update()
    if err != nil {
      initOnceErr = err
      return
    }
    helpers.SetInterval( func() {
      err := instance.Update()
      if err != nil {
        logwrapper.Logger().Error( err.Error() )
      }
    }, globals.DOCKER_API_UPDATE_INTERVAL, false )
  })

  return initOnceErr
}

func Init() error {
  if instance == nil {
    err := initOnce()
    if err != nil {
      return err
    }
  }
  return nil
}

func Instance() *DockerApi {
  return instance
}

func ( dockerApi *DockerApi ) ContainerList() []types.Container {
  return dockerApi.containerList
}

func ( dockerApi *DockerApi ) ContainerByBase64Image( image string ) *types.Container {
  if index, exists := dockerApi.containerByBase64Image[image]; exists {
    return &dockerApi.containerList[index]
  }
  return nil
}

func ( dockerApi *DockerApi ) ContainerById( id string ) *types.Container {
  if index, exists := dockerApi.containerById[id]; exists {
    return &dockerApi.containerList[index]
  }
  return nil
}

func ( dockerApi *DockerApi ) ContainerByNameRegexp( pattern string ) *types.Container {
  for _, container := range dockerApi.containerList {
    for _, name := range container.Names {
      if matched, err := regexp.MatchString( pattern, name ); matched && err==nil {
        return &container
      }
    }
  }
  return nil
}

func ( dockerApi *DockerApi ) Update() error {
  containerList, err := dockerApi.DockerClient.ContainerList(context.Background(), types.ContainerListOptions{})
  if err != nil {
    return err
  }
  dockerApi.containerList = containerList
  dockerApi.updateContainerIndexes()
  return nil
}

func ( dockerApi *DockerApi ) updateContainerIndexes() {

  dockerApi.containerById = make( map[string]int )
  dockerApi.containerByBase64Image = make( map[string]int )

  for index, container := range dockerApi.containerList {
    dockerApi.containerById[container.ID] = index
    image := strings.Split(container.Image,"@")
    base64Image := strings.Trim(base64.StdEncoding.EncodeToString([]byte(image[0])),"=")
    dockerApi.containerByBase64Image[base64Image] = index
  }
}