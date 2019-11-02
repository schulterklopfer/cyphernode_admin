package cnaOIDC

/**
 based on https://github.com/markbates/goth. Thanks a lot! :D
**/

import (
  "bytes"
  "encoding/base64"
  "encoding/json"
  "errors"
  "fmt"
  "golang.org/x/oauth2"
  "io/ioutil"
  "net/http"
  "strings"
  "time"
)

const (
  // Standard Claims http://openid.net/specs/openid-connect-core-1_0.html#StandardClaims
  // fixed, cannot be changed
  subjectClaim  = "sub"
  expiryClaim   = "exp"
  audienceClaim = "aud"
  issuerClaim   = "iss"

  PreferredUsernameClaim = "preferred_username"
  EmailClaim             = "email"
  NameClaim              = "name"
  NicknameClaim          = "nickname"
  PictureClaim           = "picture"
  GivenNameClaim         = "given_name"
  FamilyNameClaim        = "family_name"
  AddressClaim           = "address"

  // Unused but available to set in Flow claims
  MiddleNameClaim          = "middle_name"
  ProfileClaim             = "profile"
  WebsiteClaim             = "website"
  EmailVerifiedClaim       = "email_verified"
  GenderClaim              = "gender"
  BirthdateClaim           = "birthdate"
  ZoneinfoClaim            = "zoneinfo"
  LocaleClaim              = "locale"
  PhoneNumberClaim         = "phone_number"
  PhoneNumberVerifiedClaim = "phone_number_verified"
  UpdatedAtClaim           = "updated_at"

  clockSkew = 10 * time.Second
)

type Flow struct {
  ClientKey    string
  Secret       string
  CallbackURL  string
  config       *oauth2.Config
  openIDConfig *OpenIDConfig

  UserIdClaims    []string
  NameClaims      []string
  NickNameClaims  []string
  EmailClaims     []string
  AvatarURLClaims []string
  FirstNameClaims []string
  LastNameClaims  []string
  LocationClaims  []string

  SkipUserInfoRequest bool
}

type OpenIDConfig struct {
  AuthEndpoint     string `json:"authorization_endpoint"`
  TokenEndpoint    string `json:"token_endpoint"`
  LogoutEndpoint   string `json:"end_session_endpoint"`
  UserInfoEndpoint string `json:"userinfo_endpoint"`
  Issuer           string `json:"issuer"`
}

func NewFlow(clientKey, secret, callbackURL, openIDAutoDiscoveryURL string, scopes ...string) (*Flow, error) {
  flow := &Flow{
    ClientKey:   clientKey,
    Secret:      secret,
    CallbackURL: callbackURL,

    UserIdClaims:    []string{subjectClaim},
    NameClaims:      []string{NameClaim},
    NickNameClaims:  []string{NicknameClaim, PreferredUsernameClaim},
    EmailClaims:     []string{EmailClaim},
    AvatarURLClaims: []string{PictureClaim},
    FirstNameClaims: []string{GivenNameClaim},
    LastNameClaims:  []string{FamilyNameClaim},
    LocationClaims:  []string{AddressClaim},

  }

  openIDConfig, err := getOpenIDConfig(flow, openIDAutoDiscoveryURL)
  if err != nil {
    return nil, err
  }
  flow.openIDConfig = openIDConfig

  flow.config = newConfig(flow, scopes, openIDConfig)
  return flow, nil
}

func (flow *Flow) Client() *http.Client {
  return http.DefaultClient
}

// Debug is a no-op for the openidConnect package.
func (flow *Flow) Debug(debug bool) {}

// BeginAuth asks the OpenID Connect flow for an authentication end-point.
func (flow *Flow) BeginAuth(state string) (*Session, error) {
  url := flow.config.AuthCodeURL(state)
  session := &Session{
    AuthURL: url,
  }
  return session, nil
}

// FetchUser will use the the id_token and access requested information about the user.
func (flow *Flow) FetchUser(session *Session) (User, error) {
  sess := session

  expiresAt := sess.ExpiresAt

  if sess.IDToken == "" {
    return User{}, fmt.Errorf("cannot get user information without id_token")
  }

  // decode returned id token to get expiry
  claims, err := decodeJWT(sess.IDToken)

  if err != nil {
    return User{}, fmt.Errorf("oauth2: error decoding JWT token: %v", err)
  }

  expiry, err := flow.validateClaims(claims)
  if err != nil {
    return User{}, fmt.Errorf("oauth2: error validating JWT token: %v", err)
  }

  if expiry.Before(expiresAt) {
    expiresAt = expiry
  }

  if err := flow.getUserInfo(sess.AccessToken, claims); err != nil {
    return User{}, err
  }

  user := User{
    AccessToken:  sess.AccessToken,
    RefreshToken: sess.RefreshToken,
    ExpiresAt:    expiresAt,
    RawData:      claims,
  }

  flow.userFromClaims(claims, &user)
  return user, err
}

//RefreshTokenAvailable refresh token is provided by auth flow or not
func (flow *Flow) RefreshTokenAvailable() bool {
  return true
}

//RefreshToken get new access token based on the refresh token
func (flow *Flow) RefreshToken(refreshToken string) (*oauth2.Token, error) {
  token := &oauth2.Token{RefreshToken: refreshToken}
  ts := flow.config.TokenSource(oauth2.NoContext, token)
  newToken, err := ts.Token()
  if err != nil {
    return nil, err
  }
  return newToken, err
}

// validate according to standard, returns expiry
// http://openid.net/specs/openid-connect-core-1_0.html#IDTokenValidation
func (flow *Flow) validateClaims(claims map[string]interface{}) (time.Time, error) {
  audience := getClaimValue(claims, []string{audienceClaim})
  if audience != flow.ClientKey {
    found := false
    audiences := getClaimValues(claims, []string{audienceClaim})
    for _, aud := range audiences {
      if aud == flow.ClientKey {
        found = true
        break
      }
    }
    if !found {
      return time.Time{}, errors.New("audience in token does not match client key")
    }
  }

  issuer := getClaimValue(claims, []string{issuerClaim})
  if issuer != flow.openIDConfig.Issuer {
    return time.Time{}, errors.New("issuer in token does not match issuer in OpenIDConfig discovery")
  }

  // expiry is required for JWT, not for UserInfoResponse
  // is actually a int64, so force it in to that type
  expiryClaim := int64(claims[expiryClaim].(float64))
  expiry := time.Unix(expiryClaim, 0)
  if expiry.Add(clockSkew).Before(time.Now()) {
    return time.Time{}, errors.New("user info JWT token is expired")
  }
  return expiry, nil
}

func (flow *Flow) userFromClaims(claims map[string]interface{}, user *User) {
  // required
  user.UserID = getClaimValue(claims, flow.UserIdClaims)

  user.Name = getClaimValue(claims, flow.NameClaims)
  user.NickName = getClaimValue(claims, flow.NickNameClaims)
  user.Email = getClaimValue(claims, flow.EmailClaims)
  user.AvatarURL = getClaimValue(claims, flow.AvatarURLClaims)
  user.FirstName = getClaimValue(claims, flow.FirstNameClaims)
  user.LastName = getClaimValue(claims, flow.LastNameClaims)
  user.Location = getClaimValue(claims, flow.LocationClaims)
}

func (flow *Flow) getUserInfo(accessToken string, claims map[string]interface{}) error {
  // skip if there is no UserInfoEndpoint or is explicitly disabled
  if flow.openIDConfig.UserInfoEndpoint == "" || flow.SkipUserInfoRequest {
    return nil
  }

  userInfoClaims, err := flow.fetchUserInfo(flow.openIDConfig.UserInfoEndpoint, accessToken)
  if err != nil {
    return err
  }

  // The sub (subject) Claim MUST always be returned in the UserInfo Response.
  // http://openid.net/specs/openid-connect-core-1_0.html#UserInfoResponse
  userInfoSubject := getClaimValue(userInfoClaims, []string{subjectClaim})
  if userInfoSubject == "" {
    return fmt.Errorf("userinfo response did not contain a 'sub' claim: %#v", userInfoClaims)
  }

  // The sub Claim in the UserInfo Response MUST be verified to exactly match the sub Claim in the ID Token;
  // if they do not match, the UserInfo Response values MUST NOT be used.
  // http://openid.net/specs/openid-connect-core-1_0.html#UserInfoResponse
  subject := getClaimValue(claims, []string{subjectClaim})
  if userInfoSubject != subject {
    return fmt.Errorf("userinfo 'sub' claim (%s) did not match id_token 'sub' claim (%s)", userInfoSubject, subject)
  }

  // Merge in userinfo claims in case id_token claims contained some that userinfo did not
  for k, v := range userInfoClaims {
    claims[k] = v
  }

  return nil
}

// fetch and decode JSON from the given UserInfo URL
func (flow *Flow) fetchUserInfo(url, accessToken string) (map[string]interface{}, error) {
  req, _ := http.NewRequest("GET", url, nil)
  req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

  resp, err := flow.Client().Do(req)
  if err != nil {
    return nil, err
  }
  defer resp.Body.Close()

  if resp.StatusCode != http.StatusOK {
    return nil, fmt.Errorf("Non-200 response from UserInfo: %d, WWW-Authenticate=%s", resp.StatusCode, resp.Header.Get("WWW-Authenticate"))
  }

  // The UserInfo Claims MUST be returned as the members of a JSON object
  // http://openid.net/specs/openid-connect-core-1_0.html#UserInfoResponse
  data, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return nil, err
  }

  return unMarshal(data)
}

func getOpenIDConfig(p *Flow, openIDAutoDiscoveryURL string) (*OpenIDConfig, error) {
  res, err := p.Client().Get(openIDAutoDiscoveryURL)
  if err != nil {
    return nil, err
  }
  defer res.Body.Close()

  body, err := ioutil.ReadAll(res.Body)
  if err != nil {
    return nil, err
  }

  openIDConfig := &OpenIDConfig{}
  err = json.Unmarshal(body, openIDConfig)
  if err != nil {
    return nil, err
  }

  return openIDConfig, nil
}

func newConfig(flow *Flow, scopes []string, openIDConfig *OpenIDConfig) *oauth2.Config {
  c := &oauth2.Config{
    ClientID:     flow.ClientKey,
    ClientSecret: flow.Secret,
    RedirectURL:  flow.CallbackURL,
    Endpoint: oauth2.Endpoint{
      AuthURL:  openIDConfig.AuthEndpoint,
      TokenURL: openIDConfig.TokenEndpoint,
    },
    Scopes: []string{},
  }

  if len(scopes) > 0 {
    foundOpenIDScope := false

    for _, scope := range scopes {
      if scope == "openid" {
        foundOpenIDScope = true
      }
      c.Scopes = append(c.Scopes, scope)
    }

    if !foundOpenIDScope {
      c.Scopes = append(c.Scopes, "openid")
    }
  } else {
    c.Scopes = []string{"openid"}
  }

  return c
}

func getClaimValue(data map[string]interface{}, claims []string) string {
  for _, claim := range claims {
    if value, ok := data[claim]; ok {
      if stringValue, ok := value.(string); ok && len(stringValue) > 0 {
        return stringValue
      }
    }
  }

  return ""
}

func getClaimValues(data map[string]interface{}, claims []string) []string {
  var result []string

  for _, claim := range claims {
    if value, ok := data[claim]; ok {
      if stringValues, ok := value.([]interface{}); ok {
        for _, stringValue := range stringValues {
          if s, ok := stringValue.(string); ok && len(s) > 0 {
            result = append(result, s)
          }
        }
      }
    }
  }

  return result
}

// decodeJWT decodes a JSON Web Token into a simple map
// http://openid.net/specs/draft-jones-json-web-token-07.html
func decodeJWT(jwt string) (map[string]interface{}, error) {
  jwtParts := strings.Split(jwt, ".")
  if len(jwtParts) != 3 {
    return nil, errors.New("jws: invalid token received, not all parts available")
  }

  // Re-pad, if needed
  encodedPayload := jwtParts[1]
  if l := len(encodedPayload) % 4; l != 0 {
    encodedPayload += strings.Repeat("=", 4-l)
  }

  decodedPayload, err := base64.StdEncoding.DecodeString(encodedPayload)
  if err != nil {
    return nil, err
  }

  return unMarshal(decodedPayload)
}

func unMarshal(payload []byte) (map[string]interface{}, error) {
  data := make(map[string]interface{})

  return data, json.NewDecoder(bytes.NewBuffer(payload)).Decode(&data)
}
