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

package dockerApi

import (
  "bufio"
  "context"
  "encoding/base64"
  "github.com/docker/docker/api/types"
  "github.com/docker/docker/client"
  "github.com/schulterklopfer/cyphernode_admin/globals"
  "github.com/schulterklopfer/cyphernode_admin/helpers"
  "github.com/schulterklopfer/cyphernode_admin/logwrapper"
  "golang.org/x/text/unicode/norm"
  "io"
  "regexp"
  "strings"
  "sync"
)

type DockerApiConfig struct {
  Host string
  Port int
}

type DockerLogWrapper struct {
  Reader io.ReadCloser
  Active bool
  LastError error
  ContainerId string
  Lines []string
  LastLineIndex int
  Stop chan bool

  LinesMutex sync.Mutex

}

type DockerLogLine struct {
  ContainerId string
  Line string
  Active bool
}

type DockerApi struct {
  DockerClient           *client.Client
  containerList          []types.Container
  containerById          map[string]int
  containerByBase64Image map[string]int
  containerLogsById      map[string]*DockerLogWrapper

  Lines chan *DockerLogLine

  indexMutex sync.Mutex
  statsMutex sync.Mutex

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
      Lines: make( chan *DockerLogLine ),
      containerLogsById: make( map[string]*DockerLogWrapper ),
    }

    err = instance.Update()
    if err != nil {
      initOnceErr = err
      return
    }
    instance.initDockerLogCollectors()
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

func ( dockerApi *DockerApi ) initDockerLogCollectors() {
  //logwrapper.Logger().Debug(  "initDockerLogCollectors" )

  //dockerApi.indexMutex.Lock()
  //defer dockerApi.indexMutex.Unlock()

  dockerApi.checkDockerLogCollectors()

  go func() {
    for {
      select {
      case dockerLogLine := <- dockerApi.Lines:
        wrapper := dockerApi.containerLogsById[dockerLogLine.ContainerId]

        if dockerLogLine.Active {
          //logwrapper.Logger().Debug( dockerLogLine.ContainerId, "found log line" )

          wrapper.LinesMutex.Lock()
          wrapper.Lines = append(wrapper.Lines, dockerLogLine.Line)
          wrapper.LastLineIndex++
          lineCount := len(wrapper.Lines)
          if lineCount > globals.DOCKER_LOGS_MAX_LINES {
            wrapper.Lines = wrapper.Lines[lineCount-globals.DOCKER_LOGS_MAX_LINES:lineCount]
          }
          wrapper.LinesMutex.Unlock()
        } else {
          // Scanner stopped working
          logwrapper.Logger().Debug( dockerLogLine.ContainerId, "stop packet received" )
          wrapper.Active = false
          wrapper.Reader.Close()
        }
      }
    }
  }()

}

func ( dockerApi *DockerApi ) initDockerLogCollector( containerId string, dockerLogWrapper *DockerLogWrapper )  {

  //logwrapper.Logger().Debug( dockerLogWrapper.ContainerId, "initDockerLogCollector" )

  //dockerLogWrapper.LinesMutex.Lock()
  //defer dockerLogWrapper.LinesMutex.Unlock()
  //logwrapper.Logger().Debug( dockerLogWrapper.ContainerId, "trying to read docker logs for "+containerId )

  reader, err := dockerApi.LogReaderForContainer( containerId )
  dockerLogWrapper.LastError = err

  if err == nil {
    dockerLogWrapper.Active = true
    dockerLogWrapper.Reader = reader
    dockerLogWrapper.Lines = make( []string, 0)
    dockerLogWrapper.LastLineIndex = -1
    //logwrapper.Logger().Debug( dockerLogWrapper.ContainerId, "updated dockerLogWrapper" )

    go dockerApi.readDockerLogs( containerId, dockerLogWrapper )
  }

}

func ( dockerApi *DockerApi ) readDockerLogs( containerId string, dockerLogWrapper *DockerLogWrapper ) {

  logwrapper.Logger().Debug( dockerLogWrapper.ContainerId, "readDockerLogs" )

  reader := dockerLogWrapper.Reader

  scanner := bufio.NewScanner(reader)
  scanner.Split( bufio.ScanLines )

  scan := true

  go func() {
    for scan {
      logwrapper.Logger().Debug( dockerLogWrapper.ContainerId, "Waiting for stop" )
      select {
        case stop := <- dockerLogWrapper.Stop:
          if stop {
            logwrapper.Logger().Debug( dockerLogWrapper.ContainerId, "Scanning stopped from outside" )
            scan = false
          }
      }
    }
  }()

  logwrapper.Logger().Debug( dockerLogWrapper.ContainerId, "Scanning logs" )

  for scanner.Scan() && scan {
    bytes := scanner.Bytes()
    if len(bytes) < 8 {
      continue
    }
    //logwrapper.Logger().Debug( dockerLogWrapper.ContainerId, "Read bytes" )

    // cut away first 8 bytes since it only contains the length of the log
    // we could implement a custom split function which uses that information
    // but not now ;-)
    dockerApi.Lines <- &DockerLogLine{
      ContainerId: containerId,
      Line: string(norm.NFC.Bytes(bytes[8:])),
      Active: true,
    }
  }

  logwrapper.Logger().Debug( dockerLogWrapper.ContainerId, "Finished scanning" )

  _ = dockerLogWrapper.Reader.Close()

  dockerApi.Lines <- &DockerLogLine{
    ContainerId: containerId,
    Active: false,
  }
}

func ( dockerApi *DockerApi ) checkDockerLogCollectors() {
  //logwrapper.Logger().Debug( "checkDockerLogCollectors" )

  // remove logs from log db
  for containerId, dockerLogWrapper := range dockerApi.containerLogsById {
    // check if container in log db is in list of containers
    if helpers.SliceIndex(len(dockerApi.containerList), func( index int ) bool {
      return dockerApi.containerList[index].ID == containerId
    }) == -1 {
      //not found ... remove from logs db
      // send the stop signal
      dockerLogWrapper.Stop <- true
      delete( dockerApi.containerLogsById, containerId )
    }
  }

  // check containerList against containerLogsById to see what we need to add or remove
  for _, container := range dockerApi.containerList {
    if _, exists := dockerApi.containerLogsById[container.ID]; !exists {
      dockerApi.containerLogsById[container.ID] = &DockerLogWrapper{
        Active: false,
      }
    }
  }

  // finally init inactive log collectors
  for containerId, dockerLogWrapper := range dockerApi.containerLogsById {
    if !dockerLogWrapper.Active {
      if dockerLogWrapper.Reader != nil {
        _ = dockerLogWrapper.Reader.Close()
      }
      dockerApi.initDockerLogCollector( containerId, dockerLogWrapper )
    }
  }
}

func Instance() *DockerApi {
  return instance
}

func ( dockerApi *DockerApi ) ContainerList() []types.Container {
  dockerApi.indexMutex.Lock()
  defer dockerApi.indexMutex.Unlock()
  return dockerApi.containerList
}

func ( dockerApi *DockerApi ) ContainerExists( containerId string ) bool {
  dockerApi.indexMutex.Lock()
  defer dockerApi.indexMutex.Unlock()
  _, exists := dockerApi.containerById[containerId]
  return exists
}

func ( dockerApi *DockerApi ) ContainerByBase64Image( image string ) *types.Container {
  dockerApi.indexMutex.Lock()
  defer dockerApi.indexMutex.Unlock()
  if index, exists := dockerApi.containerByBase64Image[image]; exists {
    return &dockerApi.containerList[index]
  }
  return nil
}

func ( dockerApi *DockerApi ) ContainerById( id string ) *types.Container {
  dockerApi.indexMutex.Lock()
  defer dockerApi.indexMutex.Unlock()
  if index, exists := dockerApi.containerById[id]; exists {
    return &dockerApi.containerList[index]
  }
  return nil
}

func ( dockerApi *DockerApi ) LogLinesById( id string ) ([]string, int) {
  if _, exists := dockerApi.containerLogsById[id]; exists {
    dockerApi.containerLogsById[id].LinesMutex.Lock()
    defer dockerApi.containerLogsById[id].LinesMutex.Unlock()
    return dockerApi.containerLogsById[id].Lines, dockerApi.containerLogsById[id].LastLineIndex
  }
  return nil, -1
}

func ( dockerApi *DockerApi ) ContainerByNameRegexp( pattern string ) *types.Container {
  dockerApi.indexMutex.Lock()
  defer dockerApi.indexMutex.Unlock()
  for _, container := range dockerApi.containerList {
    for _, name := range container.Names {
      if matched, err := regexp.MatchString( pattern, name ); matched && err==nil {
        return &container
      }
    }
  }
  return nil
}

// TODO read into local cache and provide this local cache to websocket later
func ( dockerApi *DockerApi ) LogReaderForContainer( containerId string ) (io.ReadCloser, error) {
  //since := strconv.FormatInt(time.Now().Add(time.Duration( - 60 * time.Minute )).Unix(),10)
  reader, err := dockerApi.DockerClient.ContainerLogs( context.Background(), containerId, types.ContainerLogsOptions{
    ShowStdout: true,
    ShowStderr: true,
    Timestamps: true,
    Follow:     true,
    Details:    true,
    //Since:      since,
  })

  if err != nil {
    return nil,err
  }

  return reader, nil
}

func ( dockerApi *DockerApi ) Update() error {
  //filters := filters.NewArgs()
  //filters.Add("label", "cyphernode=true")
  containerList, err := dockerApi.DockerClient.ContainerList(context.Background(), types.ContainerListOptions{
    //Filters: filters,
  })
  if err != nil {
    return err
  }
  dockerApi.containerList = containerList
  dockerApi.updateContainerIndexes()
  dockerApi.checkDockerLogCollectors()
  return nil
}

func ( dockerApi *DockerApi ) updateContainerIndexes() {
  dockerApi.indexMutex.Lock()
  defer dockerApi.indexMutex.Unlock()

  dockerApi.containerById = make( map[string]int )
  dockerApi.containerByBase64Image = make( map[string]int )

  for index, container := range dockerApi.containerList {
    dockerApi.containerById[container.ID] = index
    image := strings.Split(container.Image,"@")
    base64Image := strings.Trim(base64.StdEncoding.EncodeToString([]byte(image[0])),"=")
    dockerApi.containerByBase64Image[base64Image] = index
  }
}
