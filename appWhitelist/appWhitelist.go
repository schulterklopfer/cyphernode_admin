package appWhitelist

import (
  "bufio"
  "github.com/schulterklopfer/cyphernode_admin/helpers"
  "github.com/schulterklopfer/cyphernode_admin/logwrapper"
  "os"
  "strings"
  "sync"
  "time"
)

var appWhitelist *AppWhitelist

type AppWhitelistEntry struct {
  ClientID string
  AppType string
}

type AppWhitelist struct {
  LastFileInfo os.FileInfo
  LastUpate time.Time
  Path string
  entries []*AppWhitelistEntry
  mutex sync.Mutex
}

func Get() *AppWhitelist {
  return appWhitelist
}

func Init( whitelistFilePath string ) {
  if appWhitelist != nil {
    return
  }
  appWhitelist = &AppWhitelist{ Path:whitelistFilePath, entries: make( []*AppWhitelistEntry, 0 ) }
  appWhitelist.load()
  helpers.SetInterval(appWhitelist.checkForWhitelistChange, 1000, false)
}

func (appWhitelist *AppWhitelist) checkForWhitelistChange() {
  fileInfo, err := os.Stat( appWhitelist.Path )
  if err != nil {
    return
  }
  if appWhitelist.LastFileInfo != nil && (
      appWhitelist.LastFileInfo.Size() != fileInfo.Size() ||
      appWhitelist.LastFileInfo.ModTime().Before( fileInfo.ModTime() ) ) {
    err := appWhitelist.load()
    if err != nil {
      logwrapper.Logger().Errorf( "Failed to load whitelist: %s", err.Error() )
    }
    appWhitelist.LastUpate = time.Now()
  }
  appWhitelist.LastFileInfo = fileInfo
}

func (appWhitelist *AppWhitelist) ContainsClientID( clientID string ) bool {
  appWhitelist.mutex.Lock()
  defer appWhitelist.mutex.Unlock()
  if len(appWhitelist.entries) == 0 {
    return false
  }

  if helpers.SliceIndex( len(appWhitelist.entries), func(i int) bool {
    return appWhitelist.entries[i].ClientID == clientID
  } ) >= 0 {
    return true
  }
  return false
}

func (appWhitelist *AppWhitelist) load() error {
  appWhitelist.mutex.Lock()
  defer appWhitelist.mutex.Unlock()
  file, err := os.Open( appWhitelist.Path )
  if err != nil {
    return err
  }
  defer file.Close()

  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    line := scanner.Text()
    if strings.HasPrefix( line, "#" ) {
      // comment... skip
      continue
    }
    arr := strings.Split(line,":")

    if len(arr) == 2 && helpers.SliceIndex( len(appWhitelist.entries), func(i int) bool {
      return appWhitelist.entries[i].ClientID == arr[0]
    } ) == -1 {
      appWhitelist.entries = append( appWhitelist.entries, &AppWhitelistEntry{ClientID:arr[0], AppType:arr[1]})
    }
  }

  if err := scanner.Err(); err != nil {
    return err
  }
  return nil
}