package cnaSessionStore

import (
  "bytes"
  "encoding/base32"
  "encoding/gob"
  "encoding/json"
  "errors"
  "fmt"
  "github.com/gin-contrib/sessions"
  "github.com/gorilla/securecookie"
  gsessions "github.com/gorilla/sessions"
  "github.com/schulterklopfer/cyphernode_admin/cnaErrors"
  "github.com/schulterklopfer/cyphernode_admin/transforms"
  "io/ioutil"
  "net/http"
  "strings"
)

var sessionExpire = 86400 * 30

type CNASessionStore struct {
  sessions.Store
  Codecs []securecookie.Codec
  options *gsessions.Options
  DefaultMaxAge int
  serializer SessionSerializer
  url string
}

func NewCNASessionStore( url string, domain string, keyPairs ...[]byte ) *CNASessionStore {

  if !strings.HasSuffix( url,"/" ) {
    url = url+"/"
  }

  store := &CNASessionStore{
    Codecs: securecookie.CodecsFromPairs(keyPairs...),
    options: &gsessions.Options{
      Domain: domain,
      Path:   "/",
      MaxAge: sessionExpire,
    },
    DefaultMaxAge: 60 * 20, // 20 minutes seems like a reasonable default
    serializer: JSONSerializer{},
    url: url,
  }
  return store
}

func (cnaSessionStore *CNASessionStore) SetMaxAge(v int) {
  var c *securecookie.SecureCookie
  var ok bool
  cnaSessionStore.options.MaxAge = v
  for i := range cnaSessionStore.Codecs {
    if c, ok = cnaSessionStore.Codecs[i].(*securecookie.SecureCookie); ok {
      c.MaxAge(v)
    } else {
      fmt.Printf("Can't change MaxAge on codec %v\n", cnaSessionStore.Codecs[i])
    }
  }
}

func (cnaSessionStore *CNASessionStore) Get(r *http.Request, name string) (*gsessions.Session, error) {
  return gsessions.GetRegistry(r).Get(cnaSessionStore, name)
}

// New should create and return a new session.
//
// Note that New should never return a nil session, even in the case of
// an error if using the Registry infrastructure to cache the session.
func (cnaSessionStore *CNASessionStore) New(r *http.Request, name string) (*gsessions.Session, error) {
  var (
    err error
    ok  bool
  )
  session := gsessions.NewSession(cnaSessionStore, name)
  // make a copy
  session.Options = cnaSessionStore.options
  session.IsNew = true
  if c, errCookie := r.Cookie(name); errCookie == nil {
    err = securecookie.DecodeMulti(name, c.Value, &session.ID, cnaSessionStore.Codecs...)
    if err == nil {
      ok, err = cnaSessionStore.load(session)
      session.IsNew = !(err == nil && ok) // not new if no error and data available
    }
  }
  return session, err
}

// Save should persist session to the underlying store implementation.
// To delete a session set session.Options.MaxAge = -1 and call Save
func (cnaSessionStore *CNASessionStore) Save(r *http.Request, w http.ResponseWriter, session *gsessions.Session) error {
  // Marked for deletion.
  if session.Options.MaxAge <= 0 {
    if err := cnaSessionStore.delete(session); err != nil {
      return err
    }
    http.SetCookie(w, gsessions.NewCookie(session.Name(), "", session.Options))
  } else {
    // Build an alphanumeric key for the store.BeA
    if session.ID == "" {
      session.ID = strings.TrimRight(base32.StdEncoding.EncodeToString(securecookie.GenerateRandomKey(32)), "=")
    }
    if err := cnaSessionStore.save(session); err != nil {
      return err
    }
    encoded, err := securecookie.EncodeMulti(session.Name(), session.ID, cnaSessionStore.Codecs...)
    if err != nil {
      return err
    }
    http.SetCookie(w, gsessions.NewCookie(session.Name(), encoded, session.Options))
  }
  return nil
}

func (cnaSessionStore *CNASessionStore) load(session *gsessions.Session) (bool, error) {

  res, err := http.DefaultClient.Get( cnaSessionStore.url+session.ID )

  if err != nil {
    return false, err
  }

  if res.StatusCode == http.StatusOK {
    // existing session
    body, err := ioutil.ReadAll(res.Body)
    if err != nil {
      return false, err
    }
    var transformedSession transforms.SessionV0
    err = json.Unmarshal(body, &transformedSession)
    if err != nil {
      return false, err
    }
    return true, cnaSessionStore.serializer.Deserialize([]byte(transformedSession.Values), session)
  } else {
    return false, cnaErrors.ErrNoSuchSession
  }
}

func (cnaSessionStore *CNASessionStore) save(session *gsessions.Session) error {

  b, err := cnaSessionStore.serializer.Serialize(session)

  if err != nil {
    return err
  }

  res, err := http.DefaultClient.Get( cnaSessionStore.url+session.ID )

  if err != nil {
    return err
  }

  var transformedSession transforms.SessionV0

  if res.StatusCode == http.StatusOK {
    // existing session
    body, err := ioutil.ReadAll(res.Body)
    if err != nil {
      return err
    }
    err = json.Unmarshal( body, &transformedSession )
    if err != nil {
      return err
    }
    transformedSession.Values=string(b)

    body, err = json.Marshal( transformedSession )
    if err != nil {
      return err
    }
    // patch
    req, err := http.NewRequest("PATCH", cnaSessionStore.url+session.ID, bytes.NewBuffer(body) )
    req.Header.Set("Content-Type", "application/json")
    if err != nil {
      return err
    }

    res, err := http.DefaultClient.Do(req)
    if err != nil {
      return err
    }

    if res.StatusCode != http.StatusOK {
      return errors.New( "session save failed [patch]" )
    }

  } else {
    // new session
    transformedSession.SessionID = session.ID
    transformedSession.Values=string(b)

    body, err := json.Marshal( transformedSession )
    if err != nil {
      return err
    }

    // create
    res, err := http.DefaultClient.Post( cnaSessionStore.url, "application/json", bytes.NewBuffer(body) )

    if err != nil {
      return err
    }

    if res.StatusCode != http.StatusCreated {
      return errors.New( "session save failed [create]" )
    }

  }

  return nil
}

func (cnaSessionStore *CNASessionStore) delete(session *gsessions.Session) error {

  if session.ID == "" {
    return nil
  }

  req, err := http.NewRequest("DELETE", cnaSessionStore.url+session.ID,nil)

  if err != nil {
    return err
  }

  res, err := http.DefaultClient.Do(req)

  if err != nil {
    return err
  }

  if res.StatusCode != http.StatusNoContent {
    return errors.New( "session delete failed" )
  }
  return nil
}

func (cnaSessionStore *CNASessionStore) Options(options sessions.Options) {
  cnaSessionStore.options = &gsessions.Options{
    Path:     options.Path,
    Domain:   options.Domain,
    MaxAge:   options.MaxAge,
    Secure:   options.Secure,
    HttpOnly: options.HttpOnly,
  }
}

/* session serializer/deserializer */

// SessionSerializer provides an interface hook for alternative serializers
type SessionSerializer interface {
  Deserialize(d []byte, ss *gsessions.Session) error
  Serialize(ss *gsessions.Session) ([]byte, error)
}

// JSONSerializer encode the session map to JSON.
type JSONSerializer struct{}

// Serialize to JSON. Will err if there are unmarshalable key values
func (s JSONSerializer) Serialize(ss *gsessions.Session) ([]byte, error) {
  m := make(map[string]interface{}, len(ss.Values))
  for k, v := range ss.Values {
    ks, ok := k.(string)
    if !ok {
      err := fmt.Errorf("Non-string key value, cannot serialize session to JSON: %v", k)
      fmt.Printf("redistore.JSONSerializer.serialize() Error: %v", err)
      return nil, err
    }
    m[ks] = v
  }

  return json.Marshal(m)
}

// Deserialize back to map[string]interface{}
func (s JSONSerializer) Deserialize(d []byte, ss *gsessions.Session) error {
  m := make(map[string]interface{})
  err := json.Unmarshal(d, &m)
  if err != nil {
    fmt.Printf("redistore.JSONSerializer.deserialize() Error: %v", err)
    return err
  }
  for k, v := range m {
    ss.Values[k] = v
  }
  return nil
}

// GobSerializer uses gob package to encode the session map
type GobSerializer struct{}

// Serialize using gob
func (s GobSerializer) Serialize(ss *gsessions.Session) ([]byte, error) {
  buf := new(bytes.Buffer)
  enc := gob.NewEncoder(buf)
  err := enc.Encode(ss.Values)
  if err == nil {
    return buf.Bytes(), nil
  }
  return nil, err
}

// Deserialize back to map[interface{}]interface{}
func (s GobSerializer) Deserialize(d []byte, ss *gsessions.Session) error {
  dec := gob.NewDecoder(bytes.NewBuffer(d))
  return dec.Decode(&ss.Values)
}