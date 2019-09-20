package hydra

import (
	"bytes"
	"encoding/json"
	"github.com/schulterklopfer/cyphernode_admin/cnaErrors"
	"github.com/schulterklopfer/cyphernode_admin/globals"
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

type Client struct {
	PolicyUri string `json:"policy_uri"`
	TosUri string `json:"tos_uri"`
	LogoUri string `json:"logo_uri"`
	ClientName string `json:"client_name"`
	ClientId string `json:"client_id"`
}

type Response struct {
	Skip bool `json:"skip"`
	RequestedScope []string `json:"requested_scope"`
	RequestedAccessTokenAudience []string `json:"requested_access_token_audience"`
	RedirectTo string `json:"redirect_to"`
	Subject string `json:"subject"`
	Client Client `json:"client"`
}

type Session struct {
}

type RequestBody struct {
	Subject string `json:"subject"`
	GrantScope []string `json:"grant_scope"`
	GrantAccessTokenAudience []string `json:"grant_access_token_audience"`
	Session *Session `json:"session"`
	Error string `json:"error"`
	ErrorDescription string `json:"error_description"`
	Remember bool `json:"remember"`
	RememberFor uint `json:"remember_for"`
	Acr uint `json:"acr"`
}

func requestFromBackend( client *http.Client, method string, flowAndAction string, challenge string, body *RequestBody ) (*Response, error) {

	u, err := absoluteURL( globals.HYDRA_FLOW_BASE_PATH + flowAndAction)
	if err != nil {
		return nil, err
	}

	u.RawQuery = flowAndAction +"_challenge="+challenge
	urlString := u.String()
	outBytes, err := json.Marshal( body )

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest( method, urlString, bytes.NewReader(outBytes) )

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

	var response Response
	err = json.Unmarshal(resBody, &response)

	if err != nil {
		return nil, err
	}

	if res.StatusCode < 200 || res.StatusCode > 302 {
		return &response, cnaErrors.ErrUnexpectedHydraResponse
	}

	return &response, nil
}

func GetLoginRequest( client *http.Client, challenge string ) (*Response, error) {
	return requestFromBackend( client, "GET", "login", challenge, nil )
}

func AcceptLoginRequest( client *http.Client, challenge string, body *RequestBody ) (*Response, error) {
	return requestFromBackend( client, "PUT", "login/accept", challenge, body )
}

func RejectLoginRequest( client *http.Client, challenge string, body *RequestBody ) (*Response, error) {
	return requestFromBackend( client, "PUT", "login/reject", challenge, body )
}

func GetConsentRequest( client *http.Client, challenge string ) (*Response, error) {
	return requestFromBackend( client, "GET", "consent", challenge, nil )
}

func AcceptConsentRequest( client *http.Client, challenge string, body *RequestBody ) (*Response, error) {
	return requestFromBackend( client, "PUT", "consent/accept", challenge, body )
}

func RejectConsentRequest( client *http.Client, challenge string, body *RequestBody ) (*Response, error) {
	return requestFromBackend( client, "PUT", "consent/reject", challenge, body )
}

func GetLogoutRequest( client *http.Client, challenge string ) (*Response, error) {
	return requestFromBackend( client, "GET", "logout", challenge, nil )
}

func AcceptLogoutRequest( client *http.Client, challenge string, body *RequestBody ) (*Response, error) {
	return requestFromBackend( client, "PUT", "logout/accept", challenge, body )
}

func RejectLogoutRequest( client *http.Client, challenge string, body *RequestBody ) (*Response, error) {
	return requestFromBackend( client, "PUT", "logout/reject", challenge, body )
}
