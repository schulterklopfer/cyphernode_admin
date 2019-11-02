package appWhitelist

import (
  "bufio"
  "github.com/schulterklopfer/cyphernode_admin/helpers"
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
    appWhitelist.load()
    appWhitelist.LastUpate = time.Now()
  }
  appWhitelist.LastFileInfo = fileInfo
}

func (appWhitelist *AppWhitelist) ContainsClientID( clientID string ) bool {
  if len(appWhitelist.entries) == 0 {
    return false
  }

  for i:=0; i<len(appWhitelist.entries); i++ {
    if appWhitelist.entries[i].ClientID == clientID {
      return true
    }
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
    if len(arr) == 2 && !appWhitelist.ContainsClientID( arr[0] ) {
      appWhitelist.entries = append( appWhitelist.entries, &AppWhitelistEntry{ClientID:arr[0], AppType:arr[1]})
    }
  }

  if err := scanner.Err(); err != nil {
    return err
  }
  return nil
}
