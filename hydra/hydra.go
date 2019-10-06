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

type Client struct {
	ClientID      string   `json:"client_id,omitempty"`
	ClientName    string   `json:"client_name,omitempty"`
	ClientSecret  string   `json:"client_secret,omitempty"`
	ClientUri     string   `json:"client_uri,omitempty"`
	RedirectUris  []string `json:"redirect_uris,omitempty"`
	GrantTypes    []string `json:"grant_types,omitempty"`
	ResponseTypes []string `json:"response_types,omitempty"`
	Scope         string   `json:"scope,omitempty"`
	Audience      []string `json:"audience,omitempty"`
	Owner         string   `json:"owner,omitempty"`
	PolicyUri     string   `json:"policy_uri,omitempty"`
	AllowedCorsOrigins []string `json:"allowed_cors_origins,omitempty"`
	TosUri string `json:"tos_uri,omitempty"`
	LogoUri string `json:"logo_uri,omitempty"`
	Contacts []string `json:"contacts,omitempty"`
	ClientSecretExpiresAt int `json:"client_secret_expires_at,omitempty"`
	SubjectType string `json:"subject_type,omitempty"`
	TokenEndpointAuthMethod string `json:"token_endpoint_auth_method,omitempty"`
	UserinfoSignedResponseAlg string `json:"userinfo_signed_response_alg,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

type Response struct {
	Skip bool `json:"skip,omitempty"`
	RequestedScope []string `json:"requested_scope,omitempty"`
	RequestedAccessTokenAudience []string `json:"requested_access_token_audience,omitempty"`
	RedirectTo string `json:"redirect_to,omitempty"`
	Subject string `json:"subject,omitempty"`
	Client Client `json:"client,omitempty"`
	Error string `json:"error,omitempty"`
	ErrorDescription string `json:"error_description,omitempty"`
}

type Session struct {
}

type RequestBody struct {
	Subject string `json:"subject,omitempty"`
	GrantScope []string `json:"grant_scope,omitempty"`
	GrantAccessTokenAudience []string `json:"grant_access_token_audience,omitempty"`
	Session *Session `json:"session,omitempty"`
	Error string `json:"error,omitempty"`
	ErrorDescription string `json:"error_description,omitempty"`
	Remember bool `json:"remember,omitempty"`
	RememberFor uint `json:"remember_for,omitempty"`
	Acr uint `json:"acr,omitempty"`
}

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

func httpRequest( client *http.Client, method string, url *url.URL, body []byte ) ([]byte, error) {

	req, err := http.NewRequest( method, url.String(), bytes.NewReader(body) )

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

	if res.StatusCode < 200 || res.StatusCode > 302 {
		return resBody, cnaErrors.ErrUnexpectedHydraResponse
	}

	return resBody, nil
}

func flowRequest( client *http.Client, method string, flow string, action string, challenge string, body *RequestBody ) (*Response, error) {

	flowAndAction := flow

	if action != "" {
		flowAndAction += "/"+action
	}

	u, err := absoluteURL( globals.HYDRA_FLOW_BASE_PATH + flowAndAction)
	if err != nil {
		return nil, err
	}

	u.RawQuery = flow +"_challenge="+challenge
	outBytes, err := json.Marshal( body )

	if err != nil {
		return nil, err
	}

	responseBytes, err := httpRequest( client, method, u, outBytes )

	if err != nil {
		return nil, err
	}

	var response Response
	err = json.Unmarshal(responseBytes, &response)

	if err != nil {
		return nil, err
	}

	return &response, nil

}

func GetLoginRequest( client *http.Client, challenge string ) (*Response, error) {
	return flowRequest( client, "GET", "login", "", challenge, nil )
}

func AcceptLoginRequest( client *http.Client, challenge string, body *RequestBody ) (*Response, error) {
	return flowRequest( client, "PUT", "login", "accept", challenge, body )
}

func RejectLoginRequest( client *http.Client, challenge string, body *RequestBody ) (*Response, error) {
	return flowRequest( client, "PUT", "login", "reject", challenge, body )
}

func GetConsentRequest( client *http.Client, challenge string ) (*Response, error) {
	return flowRequest( client, "GET", "consent", "", challenge, nil )
}

func AcceptConsentRequest( client *http.Client, challenge string, body *RequestBody ) (*Response, error) {
	return flowRequest( client, "PUT", "consent", "accept", challenge, body )
}

func RejectConsentRequest( client *http.Client, challenge string, body *RequestBody ) (*Response, error) {
	return flowRequest( client, "PUT", "consent", "reject", challenge, body )
}

func GetLogoutRequest( client *http.Client, challenge string ) (*Response, error) {
	return flowRequest( client, "GET", "logout", "", challenge, nil )
}

func AcceptLogoutRequest( client *http.Client, challenge string, body *RequestBody ) (*Response, error) {
	return flowRequest( client, "PUT", "logout", "accept", challenge, body )
}

func RejectLogoutRequest( client *http.Client, challenge string, body *RequestBody ) (*Response, error) {
	return flowRequest( client, "PUT", "logout", "reject", challenge, body )
}

func GetClients( client *http.Client )  ([]Client, error) {

	u, err := absoluteURL("/clients" )

	if err != nil {
		return nil, err
	}

	responseBytes, err := httpRequest( client, "GET", u, nil )

	if err != nil {
		return nil, err
	}

	var clients []Client
	err = json.Unmarshal(responseBytes, &clients)

	return clients, nil
}

func CreateClient( client *http.Client, oauthClient *Client )  error {
	u, err := absoluteURL("/clients" )

	if err != nil {
		return err
	}

	body, err := json.Marshal( &oauthClient )

	responseBytes, err := httpRequest( client, "POST", u, body )

	if err != nil {
		return err
	}

	err = json.Unmarshal(responseBytes, oauthClient)

	if err != nil {
		return err
	}

	return nil
}

func DeleteClient( client *http.Client, clientID string ) error {
	u, err := absoluteURL("/clients/"+clientID )

	if err != nil {
		return err
	}

	_, err = httpRequest( client, "DELETE", u, nil )

	if err != nil {
		return err
	}

	return nil
}
