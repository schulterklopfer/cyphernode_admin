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

package cyphernodeApi

import (
  "encoding/json"
  "errors"
  "strconv"
)

type Bip9SoftForkDescription struct {
  Status     string `json:"status"`
  Bit        uint8  `json:"bit"`
  StartTime1 int64  `json:"startTime"`
  StartTime2 int64  `json:"start_time"`
  Timeout    int64  `json:"timeout"`
  Since      int32  `json:"since"`
}

type UnifiedSoftFork struct {
  Type                    string                   `json:"type"`
  BIP9SoftForkDescription *Bip9SoftForkDescription `json:"bip9"`
  Height                  int32                    `json:"height"`
  Active                  bool                     `json:"active"`
}

// valid from bitcoin core > 0.19
type GetBlockChainInfoResult struct {
  Chain                string                      `json:"chain"`
  Blocks               int32                       `json:"blocks"`
  Headers              int32                       `json:"headers"`
  BestBlockHash        string                      `json:"bestblockhash"`
  Difficulty           float64                     `json:"difficulty"`
  MedianTime           int64                       `json:"mediantime"`
  VerificationProgress float64                     `json:"verificationprogress,omitempty"`
  InitialBlockDownload bool                        `json:"initialblockdownload,omitempty"`
  Pruned               bool                        `json:"pruned"`
  PruneHeight          int32                       `json:"pruneheight,omitempty"`
  ChainWork            string                      `json:"chainwork,omitempty"`
  SoftForks            map[string]*UnifiedSoftFork `json:"softforks"`
}

type GetBlockInfoResult struct {
  Hash          string   `json:"hash"`
  Confirmations int64    `json:"confirmations"`
  StrippedSize  int32    `json:"strippedsize"`
  Size          int32    `json:"size"`
  Weight        int32    `json:"weight"`
  Height        int64    `json:"height"`
  Version       int32    `json:"version"`
  VersionHex    string   `json:"versionHex"`
  MerkleRoot    string   `json:"merkleroot"`
  Tx            []string `json:"tx"`
  Time          int64    `json:"time"`
  Mediantime    int64    `json:"mediantime"`
  Nonce         uint32   `json:"nonce"`
  Bits          string   `json:"bits"`
  Difficulty    float64  `json:"difficulty"`
  Chainwork     string   `json:"chainwork"`
  TxCount       int32    `json:"nTx"`
  PreviousHash  string   `json:"previousblockhash"`
}

func (cyphernodeApi *CyphernodeApi) BitcoinCore_getBlockchainInfo() (*GetBlockChainInfoResult, error) {
  bearerToken, err := cyphernodeApi.CyphernodeKeys.BearerFromKey("000")

  if err != nil {
    return nil, err
  }

  //resp, err := cyphernodeApi.RestyClient.R().EnableTrace().SetHeader("Authorization", bearerToken).Get( cyphernodeApi.Url( "/v0/getblockchaininfo" ) )
  resp, err := cyphernodeApi.RestyClient.R().SetHeader("Authorization", bearerToken).Get(cyphernodeApi.Url("/v0/getblockchaininfo"))

  if err != nil {
    return nil, err
  }

  if resp.StatusCode() != 200 {
    println(string(resp.Body()))
    return nil, errors.New("Invalid response code: " + strconv.Itoa(resp.StatusCode()))
  }

  var getBlockChainInfoResult GetBlockChainInfoResult

  err = json.Unmarshal(resp.Body(), &getBlockChainInfoResult)

  if err != nil {
    return nil, err
  }

  /*  // Explore response object
      fmt.Println("Response Info:")
      fmt.Println("  Error      :", err)
      fmt.Println("  Status Code:", resp.StatusCode())
      fmt.Println("  Status     :", resp.Status())
      fmt.Println("  Proto      :", resp.Proto())
      fmt.Println("  Time       :", resp.Time())
      fmt.Println("  Received At:", resp.ReceivedAt())
      fmt.Println("  Body       :\n", resp)
      fmt.Println()

      // Explore trace info
      fmt.Println("Request Trace Info:")
      ti := resp.Request.TraceInfo()
      fmt.Println("  DNSLookup     :", ti.DNSLookup)
      fmt.Println("  ConnTime      :", ti.ConnTime)
      fmt.Println("  TCPConnTime   :", ti.TCPConnTime)
      fmt.Println("  TLSHandshake  :", ti.TLSHandshake)
      fmt.Println("  ServerTime    :", ti.ServerTime)
      fmt.Println("  ResponseTime  :", ti.ResponseTime)
      fmt.Println("  TotalTime     :", ti.TotalTime)
      fmt.Println("  IsConnReused  :", ti.IsConnReused)
      fmt.Println("  IsConnWasIdle :", ti.IsConnWasIdle)
      fmt.Println("  ConnIdleTime  :", ti.ConnIdleTime)
  */
  return &getBlockChainInfoResult, nil
}

type bitcoinCore_getBlockResult struct {
  Result *GetBlockInfoResult `json:"result"`
}

func (cyphernodeApi *CyphernodeApi) BitcoinCore_getBlock(hash string) (*GetBlockInfoResult, error) {
  bearerToken, err := cyphernodeApi.CyphernodeKeys.BearerFromKey("002")

  if err != nil {
    return nil, err
  }

  //resp, err := cyphernodeApi.RestyClient.R().EnableTrace().SetHeader("Authorization", bearerToken).Get( cyphernodeApi.Url( "/v0/getblockchaininfo" ) )
  resp, err := cyphernodeApi.RestyClient.R().SetHeader("Authorization", bearerToken).Get(cyphernodeApi.Url("/v0/getblockinfo/" + hash))

  if err != nil {
    return nil, err
  }

  if resp.StatusCode() != 200 {
    println(string(resp.Body()))
    return nil, errors.New("Invalid response code: " + strconv.Itoa(resp.StatusCode()))
  }

  var result bitcoinCore_getBlockResult
  body := resp.Body()
  err = json.Unmarshal(body, &result)

  if err != nil {
    return nil, err
  }

  return result.Result, nil
}
