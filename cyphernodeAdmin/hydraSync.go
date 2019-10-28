package cyphernodeAdmin

import (
  hydraAdmin "github.com/ory/hydra/sdk/go/hydra/client/admin"
  hydraModels "github.com/ory/hydra/sdk/go/hydra/models"
  "github.com/schulterklopfer/cyphernode_admin/dataSource"
  "github.com/schulterklopfer/cyphernode_admin/hydraAPI"
  "github.com/schulterklopfer/cyphernode_admin/models"
  "github.com/schulterklopfer/cyphernode_admin/queries"
)

func (cyphernodeAdmin *CyphernodeAdmin) checkHydraClients() {
  cyphernodeAdmin.addNewHydraClients()
  cyphernodeAdmin.syncHydraClients()
}

func (cyphernodeAdmin *CyphernodeAdmin) syncHydraClients() {
  params := hydraAdmin.NewListOAuth2ClientsParams()
  result, err := hydraAPI.GetBackendClient().Admin.ListOAuth2Clients(params)

  oauthClients := result.GetPayload()

  if err != nil {
    return
  }

  hydraClients := make( []models.HydraClientModel, 0 )

  err = queries.Find( &hydraClients, []interface{}{}, "", -1,0,true )

  if err != nil {
    return
  }

  // go through all oauth clients and check if corresponding hydra client exists
  // if not, delete oauth client

  for i:=0; i<len(oauthClients); i++ {
    oauthClient := oauthClients[i]
    found := false
    for j:=0; j<len(hydraClients); j++ {
      hydraClient := hydraClients[j]
      if hydraClient.ClientID == oauthClient.ClientID {
        found = true
        break
      }
    }
    if !found {
      // delete from hydra
      // TODO: HANDLE ERROR ?
      params := hydraAdmin.NewDeleteOAuth2ClientParams()
      params.ID = oauthClient.ClientID
      _, _ = hydraAPI.GetBackendClient().Admin.DeleteOAuth2Client(params)
    }
  }

  // go through all hydra clients and check if corresponding oauth client exists
  // if not, create oauth client
  db := dataSource.GetDB()

  for i:=0; i<len(hydraClients); i++ {
    hydraClient := hydraClients[i]
    found := false
    for j:=0; j<len(oauthClients); j++ {
      oauthClient := oauthClients[j]
      if hydraClient.ClientID == oauthClient.ClientID {
        found = true
        break
      }
    }
    if !found {
      // create client
      // TODO: HANDLE ERROR ?
      params := hydraAdmin.NewCreateOAuth2ClientParams()
      var oauthClient hydraModels.Client
      oauthClient.Secret = hydraClient.Secret
      oauthClient.ClientID = hydraClient.ClientID
      oauthClient.RedirectURIs = []string{ hydraClient.CallbackURL }
      oauthClient.PostLogoutRedirectURIs = []string{ hydraClient.PostLogoutCallbackURL }

      params.Body = &oauthClient
      result, err := hydraAPI.GetBackendClient().Admin.CreateOAuth2Client( params )

      if err == nil {
        hydraClient.ClientID = result.Payload.ClientID
        hydraClient.Synced = true
        db.Save( hydraClient )
      }
    }
  }
}

func (cyphernodeAdmin *CyphernodeAdmin) addNewHydraClients() {
  apps := make( []*models.AppModel, 0 )

  // select all apps with no hydra oauth client
  err := queries.Find( &apps,  []interface{}{"hydra_client_id = 0"}, "", -1,0,true)

  if err != nil {
    return
  }

  db := dataSource.GetDB()
  for i:=0; i<len(apps); i++ {
    var hydraClient models.HydraClientModel
    var oauthClient hydraModels.Client
    params := hydraAdmin.NewCreateOAuth2ClientParams()
    oauthClient.Secret = apps[i].ClientSecret
    oauthClient.ClientID = apps[i].ClientID
    oauthClient.RedirectURIs = []string{ apps[i].CallbackURL }
    oauthClient.PostLogoutRedirectURIs = []string{ apps[i].PostLogoutCallbackURL }

    params.Body = &oauthClient
    _, err := hydraAPI.GetBackendClient().Admin.CreateOAuth2Client( params )

    if err == nil {
      hydraClient.AppID = apps[i].ID
      hydraClient.ClientID = oauthClient.ClientID
      hydraClient.Secret = apps[i].ClientSecret
      hydraClient.CallbackURL = apps[i].CallbackURL
      hydraClient.PostLogoutCallbackURL = apps[i].PostLogoutCallbackURL
      hydraClient.Synced = true
      tx := db.Begin()
      tx.Save( &hydraClient )
      apps[i].HydraClientID = hydraClient.ID
      tx.Save( &apps[i] )
      tx.Commit()
    }

  }
}
