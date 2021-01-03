package cyphernodeApi

import (
  "crypto/tls"
  "crypto/x509"
  "github.com/go-resty/resty/v2"
  "github.com/schulterklopfer/cyphernode_admin/cyphernodeKeys"
  "github.com/schulterklopfer/cyphernode_admin/globals"
  "github.com/schulterklopfer/cyphernode_admin/helpers"
  "io/ioutil"
  "net/http"
  "os"
  "path"
  "strconv"
  "sync"
)

type CyphernodeApiConfig struct {
  Version string
  Host string
  Port int
}

type CyphernodeApi struct {
  Config         *CyphernodeApiConfig
  CyphernodeKeys *cyphernodeKeys.CyphernodeKeys
  RestyClient    *resty.Client
}


var instance *CyphernodeApi
var once sync.Once

func initOnce( config *CyphernodeApiConfig ) error {
  var initOnceErr error
  once.Do(func() {
    keysFile, err := os.Open( helpers.GetenvOrDefault( globals.KEYS_FILE_ENV_KEY ) )

    if err != nil {
      initOnceErr = err
      return
    }

    cnk, err := cyphernodeKeys.NewCyphernodeKeysFromFile(keysFile)

    if err != nil {
      initOnceErr = err
      return
    }

    caCert, err := ioutil.ReadFile(helpers.GetenvOrDefault( globals.CERT_FILE_ENV_KEY ))
    if err != nil {
      initOnceErr = err
      return
    }

    caCertPool := x509.NewCertPool()
    caCertPool.AppendCertsFromPEM(caCert)

    httpClient := &http.Client{
      Transport: &http.Transport{
        TLSClientConfig: &tls.Config{
          RootCAs: caCertPool,
        },
      },
    }

    restyClient := resty.NewWithClient( httpClient )

    instance = &CyphernodeApi{
      config,
      cnk,
      restyClient,
    }
  })
  return initOnceErr
}

func Init( config *CyphernodeApiConfig ) error {
  if instance == nil {
    err := initOnce( config )
    if err != nil {
      return err
    }
  }
  return nil
}

func Instance() *CyphernodeApi {
  return instance
}

func (cyphernodeApi *CyphernodeApi) Url( rel string ) string {
  return "https://"+path.Join(cyphernodeApi.Config.Host+":"+strconv.Itoa(cyphernodeApi.Config.Port), rel )
}
