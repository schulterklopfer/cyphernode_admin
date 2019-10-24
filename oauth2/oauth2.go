package oauth2

import (
  "golang.org/x/net/Context"
  "golang.org/x/oauth2"
  "net/http"
)

type OAUTH2 struct {
  Config    *oauth2.Config
  Context   context.Context
  Client    *http.Client
}

var instance *OAUTH2

func Init(config *oauth2.Config) {
  if instance != nil {
    return
  }
  instance = &OAUTH2{Config:config,Context:context.Background()}
}

func Get() *OAUTH2 {
  return instance
}

