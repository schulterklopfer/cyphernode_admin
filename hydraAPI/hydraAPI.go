package hydraAPI

import (
	"github.com/ory/hydra/sdk/go/hydra/client"
	"github.com/ory/x/urlx"
	"github.com/schulterklopfer/cyphernode_admin/globals"
	"github.com/schulterklopfer/cyphernode_admin/logwrapper"
	"os"
)

var backendClient *client.OryHydra

func GetBackendClient() *client.OryHydra {
	return backendClient
}

func Init() {
	if backendClient != nil {
		return
	}
	var err error
	logwrapper.Logger().Info( "Initialising hydra backend API")
	backendClient = client.NewHTTPClientWithConfig(nil, &client.TransportConfig{Schemes: []string{"http"}, Host: urlx.ParseOrPanic(os.Getenv( globals.HYDRA_ADMIN_URL_ENV_KEY )).Host})
	if err != nil {
		logwrapper.Logger().Panic("failed to connect to database" )
	}
}
