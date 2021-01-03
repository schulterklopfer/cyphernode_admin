package cyphernodeApi

import (
  "encoding/json"
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
  Chain                string  `json:"chain"`
  Blocks               int32   `json:"blocks"`
  Headers              int32   `json:"headers"`
  BestBlockHash        string  `json:"bestblockhash"`
  Difficulty           float64 `json:"difficulty"`
  MedianTime           int64   `json:"mediantime"`
  VerificationProgress float64 `json:"verificationprogress,omitempty"`
  InitialBlockDownload bool    `json:"initialblockdownload,omitempty"`
  Pruned               bool    `json:"pruned"`
  PruneHeight          int32   `json:"pruneheight,omitempty"`
  ChainWork            string  `json:"chainwork,omitempty"`
  SoftForks map[string]*UnifiedSoftFork `json:"softforks"`
}

func (cyphernodeApi *CyphernodeApi) BitcoinCore_getBlockchainInfo() (*GetBlockChainInfoResult, error) {
  bearerToken, err := cyphernodeApi.CyphernodeKeys.BearerFromKey("000")

  if err != nil {
    return nil, err
  }

  //resp, err := cyphernodeApi.RestyClient.R().EnableTrace().SetHeader("Authorization", bearerToken).Get( cyphernodeApi.Url( "/v0/getblockchaininfo" ) )
  resp, err := cyphernodeApi.RestyClient.R().SetHeader("Authorization", bearerToken).Get( cyphernodeApi.Url( "/v0/getblockchaininfo" ) )

  if err != nil {
    return nil, err
  }

  var getBlockChainInfoResult GetBlockChainInfoResult

  err = json.Unmarshal( resp.Body(), &getBlockChainInfoResult )

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
