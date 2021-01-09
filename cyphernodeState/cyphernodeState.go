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
  "github.com/schulterklopfer/cyphernode_admin/helpers"
  "github.com/schulterklopfer/cyphernode_admin/logwrapper"
  "sync"
  "time"
)

type CyphernodeState struct {
  LastUpdate                time.Time                              `json:"lastUpdate"`
  BlockchainInfo            *cyphernodeApi.GetBlockChainInfoResult `json:"blockchainInfo"`
  CyphernodeInfo            *cyphernodeInfo.CyphernodeInfo         `json:"cyphernodeInfo"`
  LatestBlocks              []*cyphernodeApi.GetBlockInfoResult    `json:"latestBlocks"`
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
      LatestBlocks: []*cyphernodeApi.GetBlockInfoResult{},
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