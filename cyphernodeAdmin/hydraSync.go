package cyphernodeAdmin

import (
	"github.com/schulterklopfer/cyphernode_admin/dataSource"
	"github.com/schulterklopfer/cyphernode_admin/hydra"
	"github.com/schulterklopfer/cyphernode_admin/models"
	"github.com/schulterklopfer/cyphernode_admin/queries"
	"net/http"
)

func (cyphernodeAdmin *CyphernodeAdmin) checkHydraClients() {
	cyphernodeAdmin.addNewHydraClients()
	cyphernodeAdmin.syncHydraClients()
}

func (cyphernodeAdmin *CyphernodeAdmin) syncHydraClients() {
	oauthClients, err := hydra.GetClients( http.DefaultClient )

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
			_ = hydra.DeleteClient( http.DefaultClient, oauthClient.ClientID )
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
			var oauthClient hydra.Client
			oauthClient.ClientSecret = hydraClient.Secret

			err = hydra.CreateClient( http.DefaultClient, &oauthClient )

			if err == nil {
				hydraClient.ClientID = oauthClient.ClientID
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
		var oauthClient hydra.Client

		oauthClient.ClientSecret = apps[i].Hash

		err := hydra.CreateClient( http.DefaultClient, &oauthClient )
		if err == nil {
			hydraClient.AppID = apps[i].ID
			hydraClient.ClientID = oauthClient.ClientID
			hydraClient.Secret = oauthClient.ClientSecret
			tx := db.Begin()
			tx.Save( &hydraClient )
			apps[i].HydraClientID = hydraClient.ID
			tx.Save( &apps[i] )
			tx.Commit()
		}

	}
}
