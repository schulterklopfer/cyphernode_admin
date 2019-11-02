package cnaOIDC

/**
  based on https://github.com/markbates/goth. Thanks a lot! :D
 **/

import (
  "encoding/base64"
  "errors"
  "fmt"
  "github.com/gorilla/sessions"
  "github.com/schulterklopfer/cyphernode_admin/cnaErrors"
  "github.com/schulterklopfer/cyphernode_admin/cnaSessionStore"
  "github.com/schulterklopfer/cyphernode_admin/globals"
  "github.com/schulterklopfer/cyphernode_admin/helpers"
  "net/http"
  "net/url"
)


// SessionName is the key used to access the session store.
const SessionName = globals.SESSION_COOKIE_NAME
const sessionDataKey = "oidc"

var SessionStore *cnaSessionStore.CNASessionStore
var flow *Flow

type InitParams struct {
  ClientID string
  ClientSecret string
  CallbackURL string
  OIDCDiscoveryURL string
  SessionsEndpoint string
  CookieSecret []byte
  CookieDomain string
}

func NewInitParams( clientID string, clientSecret string, callbackURL string, OIDCDiscoveryURL string, sessionsEndpoint string, sessionStoreCookieSecret []byte, cookieDomain string ) *InitParams {
  return &InitParams{
    ClientID:                 clientID,
    ClientSecret:             clientSecret,
    CallbackURL:              callbackURL,
    OIDCDiscoveryURL:         OIDCDiscoveryURL,
    SessionsEndpoint:         sessionsEndpoint,
    CookieSecret:             sessionStoreCookieSecret,
    CookieDomain:             cookieDomain,
  }
}

func Init( params *InitParams ) {
  flow, _ =  NewFlow(params.ClientID, params.ClientSecret, params.CallbackURL, params.OIDCDiscoveryURL )
  SessionStore = cnaSessionStore.NewCNASessionStore( params.SessionsEndpoint, params.CookieDomain, params.CookieSecret )
}

func BeginAuthHandler(res http.ResponseWriter, req *http.Request) {
  url, err := GetAuthURL(res, req)
  if err != nil {
    res.WriteHeader(http.StatusBadRequest)
    fmt.Fprintln(res, err)
    return
  }

  http.Redirect(res, req, url, http.StatusTemporaryRedirect)
}

// SetState sets the state string associated with the given request.
// If no state string is associated with the request, one will be generated.
// This state is sent to the flow and can be retrieved during the
// callback.
var SetState = func(req *http.Request) string {
  state := req.URL.Query().Get("state")
  if len(state) > 0 {
    return state
  }

  // If a state query param is not passed in, generate a random
  // base64-encoded nonce so that the state on the auth URL
  // is unguessable, preventing CSRF attacks, as described in
  //
  // https://auth0.com/docs/protocols/oauth2/oauth-state#keep-reading
  return createState()
}

func createState() string {
  return helpers.RandomString(64, base64.URLEncoding.EncodeToString )
}

// GetState gets the state returned by the flow during the callback.
// This is used to prevent CSRF attacks, see
// http://tools.ietf.org/html/rfc6749#section-10.12
var GetState = func(req *http.Request) string {
  return req.URL.Query().Get("state")
}

func GetAuthURL(res http.ResponseWriter, req *http.Request) (string, error) {
  if SessionStore == nil {
    return "", cnaErrors.ErrNoSessionStore
  }

  sess, err := flow.BeginAuth(SetState(req))
  if err != nil {
    return "", err
  }

  url, err := sess.GetAuthURL()
  if err != nil {
    return "", err
  }

  err = StoreInSession(sessionDataKey, sess.Marshal(), req, res)

  if err != nil {
    return "", err
  }

  return url, err
}

var CompleteUserAuth = func(res http.ResponseWriter, req *http.Request) (User, error) {
  if SessionStore == nil {
    return User{}, cnaErrors.ErrNoSessionStore
  }

  value, err := GetFromSession(sessionDataKey, req)
  if err != nil {
    return User{}, err
  }

  sess, err := flow.UnmarshalSession(value)
  if err != nil {
    return User{}, err
  }

  err = validateState(req.URL.Query().Get("state"), sess)
  if err != nil {
    return User{}, err
  }

  user, err := flow.FetchUser(sess)
  if err == nil {
    // user can be found with existing session data
    return user, err
  }

  // get new token and retry fetch
  _, err = sess.Authorize(flow, req.URL.Query())
  if err != nil {
    return User{}, err
  }

  err = StoreInSession(sessionDataKey, sess.Marshal(), req, res)

  if err != nil {
    return User{}, err
  }

  gu, err := flow.FetchUser(sess)
  return gu, err
}

// validateState ensures that the state token param from the original
// AuthURL matches the one included in the current (callback) request.
func validateState(state string, sess *Session) error {
  rawAuthURL, err := sess.GetAuthURL()
  if err != nil {
    return err
  }

  authURL, err := url.Parse(rawAuthURL)
  if err != nil {
    return err
  }

  originalState := authURL.Query().Get("state")
  if originalState != "" && (originalState != state ) {
    return errors.New("state token mismatch")
  }
  return nil
}

// Logout invalidates a user session.
func Logout( res http.ResponseWriter, req *http.Request, postLogoutRedirectURL string ) error {

  if SessionStore == nil {
    return cnaErrors.ErrNoSessionStore
  }

  if SessionStore == nil {
    return cnaErrors.ErrNoSessionStore
  }

  session, err := SessionStore.Get(req, SessionName)
  if err != nil {
    return err
  }
  session.Options.MaxAge = -1
  session.Values = make(map[interface{}]interface{})
  err = session.Save(req, res)
  if err != nil {
    return errors.New("Could not delete user session ")
  }
  return nil
}

// StoreInSession stores a specified key/value pair in the session.
func StoreInSession(key string, value string, req *http.Request, res http.ResponseWriter) error {
  session, _ := SessionStore.New(req, SessionName)

  if err := updateSessionValue(session, key, value); err != nil {
    return err
  }

  return session.Save(req, res)
}

// GetFromSession retrieves a previously-stored value from the session.
// If no value has previously been stored at the specified key, it will return an error.
func GetFromSession(key string, req *http.Request) (string, error) {
  session, _ := SessionStore.Get(req, SessionName)
  value, err := getSessionValue(session, key)
  if err != nil {
    return "", errors.New("could not find a matching session for this request")
  }

  return value, nil
}

var GetUser = func(res http.ResponseWriter, req *http.Request) (User, error) {
  if SessionStore == nil {
    return User{}, cnaErrors.ErrNoSessionStore
  }

  value, err := GetFromSession(sessionDataKey, req)
  if err != nil {
    return User{}, err
  }

  sess, err := flow.UnmarshalSession(value)
  if err != nil {
    return User{}, err
  }

  user, err := flow.FetchUser(sess)
  if err == nil {
    // user can be found with existing session data
    return user, err
  }

  // get new token and retry fetch
  _, err = sess.Authorize(flow, req.URL.Query())
  if err != nil {
    return User{}, err
  }

  err = StoreInSession(sessionDataKey, sess.Marshal(), req, res)

  if err != nil {
    return User{}, err
  }

  gu, err := flow.FetchUser(sess)
  return gu, err
}

func getSessionValue(session *sessions.Session, key string) (string, error) {
  value := session.Values[key]
  if value == nil {
    return "", fmt.Errorf("could not find a matching session for this request")
  }

  /*
    rdata := strings.NewReader(value.(string))
    r, err := gzip.NewReader(rdata)
    if err != nil {
      return "", err
    }
    s, err := ioutil.ReadAll(r)
    if err != nil {
      return "", err
    }

    return string(s), nil
  */
  return value.(string), nil
}

func updateSessionValue(session *sessions.Session, key, value string) error {
  /*
    var b bytes.Buffer
    gz := gzip.NewWriter(&b)
    // TOOD: debug gzip compression... Probably is has sth to do with session serialisation. Maybe
    // try sth other than json marshalling
    if _, err := gz.Write([]byte(value)); err != nil {
      return err
    }
    if err := gz.Flush(); err != nil {
      return err
    }
    if err := gz.Close(); err != nil {
      return err
    }

    session.Values[key] = b.String()
  */
  session.Values[key] = value
  return nil
}