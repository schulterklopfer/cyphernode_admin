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
  "github.com/schulterklopfer/cyphernode_admin/cyphernodeKeys"
  "strconv"
)

type GetLightningNodeInfoResult struct {
  Id string `json:"id"`
  Alias string `json:"alias"`
  Color string `json:"color"`
  NumPeers int `json:"num_peers"`
  NumPendingChannels int `json:"num_pending_channels"`
  NumActiveChannels int `json:"num_active_channels"`
  NumInactiveChannels int `json:"num_inactive_channels"`
  Address []string `json:"address"`
  Version string `json:"version"`
  BlockHeight int `json:"blockheight"`
  Network string `json:"network"`
  MsatoshiFeesCollected int `json:"msatoshi_fees_collected"`
}

func (cyphernodeApi *CyphernodeApi) Lightning_getNodeInfo() (*GetLightningNodeInfoResult, error) {
  bearerToken, err := cyphernodeKeys.Instance().BearerFromKey("001")

  if err != nil {
    return nil, err
  }

  //resp, err := cyphernodeApi.RestyClient.R().EnableTrace().SetHeader("Authorization", bearerToken).Get( cyphernodeApi.Url( "/v0/getblockchaininfo" ) )
  resp, err := cyphernodeApi.RestyClient.R().SetHeader("Authorization", bearerToken).Get(cyphernodeApi.Url("/v0/ln_getinfo"))

  if err != nil {
    return nil, err
  }

  if resp.StatusCode() != 200 {
    println(string(resp.Body()))
    return nil, errors.New("Invalid response code: " + strconv.Itoa(resp.StatusCode()))
  }

  var result GetLightningNodeInfoResult
  body := resp.Body()
  err = json.Unmarshal(body, &result)

  if err != nil {
    return nil, err
  }

  return &result, nil
}

