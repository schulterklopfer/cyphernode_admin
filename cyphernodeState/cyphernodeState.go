package cyphernodeState

import (
  "github.com/SatoshiPortal/cam/cyphernodeInfo"
  "github.com/schulterklopfer/cyphernode_admin/cyphernodeApi"
  "github.com/schulterklopfer/cyphernode_admin/globals"
  "github.com/schulterklopfer/cyphernode_admin/helpers"
  "github.com/schulterklopfer/cyphernode_admin/logwrapper"
  "sync"
  "time"
)

type CyphernodeState struct {
  LastUpdate                time.Time                              `json:"lastUpdate"`
  BlockchainInfo            *cyphernodeApi.GetBlockChainInfoResult `json:"blockchainInfo"`
  CyphernodeInfo            *cyphernodeInfo.CyphernodeInfo         `json:"cyphernodeInfo"`
  updateBlockchainInfoMutex sync.Mutex                             `json:"-"`
}

var instance *CyphernodeState
var once sync.Once

func initOnce( info *cyphernodeInfo.CyphernodeInfo ) error {
  var initOnceErr error
  once.Do(func() {
    instance = &CyphernodeState{
      BlockchainInfo: nil,
      CyphernodeInfo: info,
    }
    err := instance.updateBlockchainInfo()
    if err != nil {
      initOnceErr = err
      return
    }
    helpers.SetInterval( func() {
      err := instance.updateBlockchainInfo()
      if err != nil {
        logwrapper.Logger().Error( err.Error() )
      }
    }, globals.BLOCKCHAIN_INFO_UPDATE_INTERVAL, false )
  })
  return initOnceErr
}

func Init( info *cyphernodeInfo.CyphernodeInfo ) error {
  if instance == nil {
    err := initOnce( info )
    if err != nil {
      return err
    }
  }
  return nil
}

func Instance() *CyphernodeState {
  return instance
}

func ( cyphernodeState *CyphernodeState ) updateBlockchainInfo() error {
  cyphernodeState.updateBlockchainInfoMutex.Lock()
  defer cyphernodeState.updateBlockchainInfoMutex.Unlock()
  blockchainInfo, err := cyphernodeApi.Instance().BitcoinCore_getBlockchainInfo()
  if err != nil {
    return err
  }
  cyphernodeState.LastUpdate = time.Now()
  cyphernodeState.BlockchainInfo = blockchainInfo
  return nil
}