package hydra

import (
	"encoding/json"
	"github.com/schulterklopfer/cyphernode_admin/cnaErrors"
	"github.com/schulterklopfer/cyphernode_admin/globals"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

func absoluteURL( path string ) (*url.URL, error) {
	hydraAdminUrlString := os.Getenv( globals.HYDRA_ADMIN_URL_ENV_KEY )

	if hydraAdminUrlString == "" {
		return nil, cnaErrors.ErrNoHydraAdminUrl
	}

	hydraAdminUrl, err := url.Parse(hydraAdminUrlString)
	if err != nil {
		return nil, err
	}

	relativeUrl, err := url.Parse( path )
	if err != nil {
		return nil, err
	}

	u := hydraAdminUrl.ResolveReference(relativeUrl)
	return u, nil
}

func requestFromBackend( client *http.Client, method string, flowAndAction string, challenge string, body io.Reader ) (*map[string]interface{}, error) {

	u, err := absoluteURL( globals.HYDRA_FLOW_BASE_PATH + flowAndAction)
	if err != nil {
		return nil, err
	}

	u.RawQuery = flowAndAction +"_challenge="+challenge
	urlString := u.String()
	req, err := http.NewRequest( method, urlString, body )

	if err != nil {
		return nil, err
	}

	res, err := client.Do( req )

	if err != nil {
		return nil, err
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if err != nil {
		return nil, err
	}

	var resMap map[string]interface{}
	err = json.Unmarshal(resBody, &resMap)

	if err != nil {
		return nil, err
	}

	if res.StatusCode < 200 || res.StatusCode > 302 {
		return &resMap, cnaErrors.ErrUnexpectedHydraResponse
	}

	return &resMap, nil
}

func GetLoginRequest( client *http.Client, challenge string ) (*map[string]interface{}, error) {
	return requestFromBackend( client, "GET", "login", challenge, nil )
}

func AcceptLoginRequest( client *http.Client, challenge string, body io.Reader ) (*map[string]interface{}, error) {
	return requestFromBackend( client, "PUT", "login/accept", challenge, body )
}

func RejectLoginRequest( client *http.Client, challenge string, body io.Reader ) (*map[string]interface{}, error) {
	return requestFromBackend( client, "PUT", "login/reject", challenge, body )
}

func GetConsentRequest( client *http.Client, challenge string ) (*map[string]interface{}, error) {
	return requestFromBackend( client, "GET", "consent", challenge, nil )
}

func AcceptConsentRequest( client *http.Client, challenge string, body io.Reader ) (*map[string]interface{}, error) {
	return requestFromBackend( client, "PUT", "consent/accept", challenge, body )
}

func RejectConsentRequest( client *http.Client, challenge string, body io.Reader ) (*map[string]interface{}, error) {
	return requestFromBackend( client, "PUT", "consent/reject", challenge, body )
}

func GetLogoutRequest( client *http.Client, challenge string ) (*map[string]interface{}, error) {
	return requestFromBackend( client, "GET", "logout", challenge, nil )
}

func AcceptLogoutRequest( client *http.Client, challenge string, body io.Reader ) (*map[string]interface{}, error) {
	return requestFromBackend( client, "PUT", "logout/accept", challenge, body )
}

func RejectLogoutRequest( client *http.Client, challenge string, body io.Reader ) (*map[string]interface{}, error) {
	return requestFromBackend( client, "PUT", "logout/reject", challenge, body )
}
