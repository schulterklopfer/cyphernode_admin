package dockerApi

import (
  "context"
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
  containerById          map[string]*types.Container
  containerByHashedImage map[string]*types.Container
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

func ( dockerApi *DockerApi ) ContainerByHashedImage( image string ) *types.Container {
  return dockerApi.containerByHashedImage[image]
}

func ( dockerApi *DockerApi ) ContainerById( id string ) *types.Container {
  return dockerApi.containerById[id]
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

  dockerApi.containerById = make( map[string]*types.Container )
  dockerApi.containerByHashedImage = make( map[string]*types.Container )

  for _, container := range dockerApi.containerList {
    dockerApi.containerById[container.ID] = &container
    image := strings.Split(container.Image,"@")
    hashBytes := helpers.TrimmedRipemd160Hash( []byte(image[0]) )
    dockerApi.containerByHashedImage[string(hashBytes)] = &container
  }
}