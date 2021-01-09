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

  LinesMutex sync.Mutex

}

type DockerLogLine struct {
  ContainerId string
  Line string
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
  dockerApi.indexMutex.Lock()
  defer dockerApi.indexMutex.Unlock()
  dockerApi.containerLogsById = make( map[string]*DockerLogWrapper )

  for _, container := range dockerApi.containerList {
    dockerLogWrapper := DockerLogWrapper{
      Active: false,
    }
    dockerApi.initDockerLogCollector( container.ID, &dockerLogWrapper)
    dockerApi.containerLogsById[container.ID] = &dockerLogWrapper
  }

  go func() {
    for {
      select {
      case dockerLogLine := <- dockerApi.Lines:
        wrapper := dockerApi.containerLogsById[dockerLogLine.ContainerId]
        wrapper.LinesMutex.Lock()
        wrapper.Lines = append(wrapper.Lines, dockerLogLine.Line)
        wrapper.LastLineIndex++
        lineCount := len(wrapper.Lines)
        if lineCount > globals.DOCKER_LOGS_MAX_LINES {
          wrapper.Lines = wrapper.Lines[lineCount-globals.DOCKER_LOGS_MAX_LINES:]
        }
        wrapper.LinesMutex.Unlock()
      }
    }
  }()

}

func ( dockerApi *DockerApi ) initDockerLogCollector( containerId string, dockerLogWrapper *DockerLogWrapper )  {
  reader, err := dockerApi.LogReaderForContainer( containerId )
  dockerLogWrapper.LastError = err

  if err == nil {
    dockerLogWrapper.Reader = reader
    dockerLogWrapper.Lines = make( []string, 0)
    dockerLogWrapper.Active = true
    dockerLogWrapper.LastLineIndex = -1
    go dockerApi.readDockerLogs( containerId, dockerLogWrapper.Reader )
  }

}

func ( dockerApi *DockerApi ) readDockerLogs( containerId string, reader io.ReadCloser ) {
  scanner := bufio.NewScanner(reader)
  scanner.Split( bufio.ScanLines )

  for scanner.Scan() {
    bytes := scanner.Bytes()
    if len(bytes) < 8 {
      continue
    }
    // cut away first 8 bytes since it only contains the length of the log
    // we could implement a custom split function which uses that information
    // but not now ;-)
    dockerApi.Lines <- &DockerLogLine{
      ContainerId: containerId,
      Line: string(norm.NFC.Bytes(bytes[8:])),
    }
  }
}

func ( dockerApi *DockerApi ) checkDockerLogCollectors() {
  for containerId, dockerLogWrapper := range dockerApi.containerLogsById {
    if !dockerLogWrapper.Active {
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
