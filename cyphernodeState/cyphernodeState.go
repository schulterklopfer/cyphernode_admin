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

package cyphernodeState

import (
  "github.com/SatoshiPortal/cam/cyphernodeInfo"
  "github.com/schulterklopfer/cyphernode_admin/cyphernodeApi"
  "github.com/schulterklopfer/cyphernode_admin/globals"
  "github.com/schulterklopfer/cyphernode_fauth/helpers"
  "github.com/schulterklopfer/cyphernode_fauth/logwrapper"
  "sync"
  "time"
)

type CyphernodeState struct {
  LastUpdate                   time.Time                                 `json:"lastUpdate"`
  BlockchainInfo               *cyphernodeApi.GetBlockChainInfoResult    `json:"blockchainInfo"`
  LightningNodeInfo            *cyphernodeApi.GetLightningNodeInfoResult `json:"lightningNodeInfo"`
  CyphernodeInfo               *cyphernodeInfo.CyphernodeInfo            `json:"cyphernodeInfo"`
  LatestBlocks                 []*cyphernodeApi.GetBlockInfoResult       `json:"latestBlocks"`
  updateBlockchainInfoMutex    sync.Mutex                                `json:"-"`
  updateLightningNodeInfoMutex sync.Mutex                                `json:"-"`

}

var instance *CyphernodeState
var once sync.Once

func initOnce( info *cyphernodeInfo.CyphernodeInfo ) error {
  var initOnceErr error
  once.Do(func() {
    instance = &CyphernodeState{
      BlockchainInfo: nil,
      CyphernodeInfo: info,
      LatestBlocks: []*cyphernodeApi.GetBlockInfoResult{},
    }

    upbci := func() {
      trying := true

      for trying {
        err := instance.updateBlockchainInfo()
        if err == nil {
          trying = false
          continue
        }
        logwrapper.Logger().Error( err.Error() )
        time.Sleep( 1*time.Second )
      }
    }

    upbci()
    helpers.SetInterval( upbci, globals.BLOCKCHAIN_INFO_UPDATE_INTERVAL, false )

    featureIndex := helpers.SliceIndex(len(info.OptionalFeatures), func(i int) bool {
      return info.OptionalFeatures[i].Label == "lightning" && info.OptionalFeatures[i].Active
    })

    if featureIndex > -1 {
      uplni := func() {
        trying := true

        for trying {
          err := instance.updateLightningNodeInfo()
          if err == nil {
            trying = false
            info.OptionalFeatures[featureIndex].Extra["pubkey"] = instance.LightningNodeInfo.Id
            continue
          }
          logwrapper.Logger().Error( err.Error() )
          time.Sleep( 1*time.Second )
        }
      }

      uplni()
      helpers.SetInterval( uplni, globals.BLOCKCHAIN_INFO_UPDATE_INTERVAL, false )
    }
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

  if err != nil {
    return err
  }

  var blocks []*cyphernodeApi.GetBlockInfoResult
  hash := blockchainInfo.BestBlockHash

  for i:=0; i<globals.LATEST_BLOCK_COUNT; i++ {
    block, err := cyphernodeApi.Instance().BitcoinCore_getBlock( hash )
    if err != nil {
      break
    }
    block.Tx = []string{}
    blocks = append(blocks,block)
    hash = block.PreviousHash
  }

  cyphernodeState.LatestBlocks = blocks
  return nil
}

func ( cyphernodeState *CyphernodeState ) updateLightningNodeInfo() error {
  cyphernodeState.updateLightningNodeInfoMutex.Lock()
  defer cyphernodeState.updateLightningNodeInfoMutex.Unlock()
  lightningNodeInfo, err := cyphernodeApi.Instance().Lightning_getNodeInfo()
  if err != nil {
    return err
  }
  cyphernodeState.LastUpdate = time.Now()
  cyphernodeState.LightningNodeInfo = lightningNodeInfo
  return nil
}