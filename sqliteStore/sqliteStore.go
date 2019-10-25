package sqliteStore

import (
  "encoding/base32"
  "encoding/json"
  "fmt"
  "github.com/gin-contrib/sessions"
  "github.com/gorilla/securecookie"
  gsessions "github.com/gorilla/sessions"
  "github.com/schulterklopfer/cyphernode_admin/cnaErrors"
  "github.com/schulterklopfer/cyphernode_admin/models"
  "github.com/schulterklopfer/cyphernode_admin/queries"
  "net/http"
  "strings"
)

var sessionExpire = 86400 * 30

type SqliteStore struct {
  sessions.Store
  Codecs []securecookie.Codec
  options *gsessions.Options
  DefaultMaxAge int
  serializer SessionSerializer
}

type Session struct {
  sessions.Session
}

type Options struct {
  gsessions.Options
}

func NewSqliteStore( keyPairs ...[]byte ) *SqliteStore {
  store := &SqliteStore{
    Codecs: securecookie.CodecsFromPairs(keyPairs...),
    options: &gsessions.Options{
      Path:   "/",
      MaxAge: sessionExpire,
    },
    DefaultMaxAge: 60 * 20, // 20 minutes seems like a reasonable default
    serializer: JSONSerializer{},
  }
  return store
}

func (sqliteStore *SqliteStore) SetMaxAge(v int) {
  var c *securecookie.SecureCookie
  var ok bool
  sqliteStore.options.MaxAge = v
  for i := range sqliteStore.Codecs {
    if c, ok = sqliteStore.Codecs[i].(*securecookie.SecureCookie); ok {
      c.MaxAge(v)
    } else {
      fmt.Printf("Can't change MaxAge on codec %v\n", sqliteStore.Codecs[i])
    }
  }
}

func (sqliteStore *SqliteStore) Get(r *http.Request, name string) (*gsessions.Session, error) {
  return gsessions.GetRegistry(r).Get(sqliteStore, name)
}

// New should create and return a new session.
//
// Note that New should never return a nil session, even in the case of
// an error if using the Registry infrastructure to cache the session.
func (sqliteStore *SqliteStore) New(r *http.Request, name string) (*gsessions.Session, error) {
  var (
    err error
    ok  bool
  )
  session := gsessions.NewSession(sqliteStore, name)
  // make a copy
  session.Options = sqliteStore.options
  session.IsNew = true
  if c, errCookie := r.Cookie(name); errCookie == nil {
    err = securecookie.DecodeMulti(name, c.Value, &session.ID, sqliteStore.Codecs...)
    if err == nil {
      ok, err = sqliteStore.load(session)
      session.IsNew = !(err == nil && ok) // not new if no error and data available
    }
  }
  return session, err
}

// Save should persist session to the underlying store implementation.
// To delete a session set session.Options.MaxAge = -1 and call Save
func (sqliteStore *SqliteStore) Save(r *http.Request, w http.ResponseWriter, session *gsessions.Session) error {
  // Marked for deletion.
  if session.Options.MaxAge <= 0 {
    if err := sqliteStore.delete(session); err != nil {
      return err
    }
    http.SetCookie(w, gsessions.NewCookie(session.Name(), "", session.Options))
  } else {
    // Build an alphanumeric key for the store.BeA
    if session.ID == "" {
      session.ID = strings.TrimRight(base32.StdEncoding.EncodeToString(securecookie.GenerateRandomKey(32)), "=")
    }
    if err := sqliteStore.save(session); err != nil {
      return err
    }
    encoded, err := securecookie.EncodeMulti(session.Name(), session.ID, sqliteStore.Codecs...)
    if err != nil {
      return err
    }
    http.SetCookie(w, gsessions.NewCookie(session.Name(), encoded, session.Options))
  }
  return nil
}

func (sqliteStore *SqliteStore) load(session *gsessions.Session) (bool, error) {

  // TODO: handle case:
  // * cookie is set, but session is not in db
  var sessionModels []*models.SessionModel
  err := queries.Find( &sessionModels,  []interface{}{"session_id = ?", session.ID}, "", 1,0,false)

  if err != nil {
    return false, err
  }

  if len(sessionModels) == 0 {
    return false, cnaErrors.ErrNoSuchSession
  }

  return true, sqliteStore.serializer.Deserialize([]byte(sessionModels[0].Values), session)
}

func (sqliteStore *SqliteStore) save(session *gsessions.Session) error {
  // TODO save session in sqlite table

  b, err := sqliteStore.serializer.Serialize(session)

  if err != nil {
    return err
  }

  var sessionModels []*models.SessionModel
  err = queries.Find( &sessionModels,  []interface{}{"session_id = ?", session.ID}, "", 1,0,false)

  if err != nil {
    return err
  }

  var sessionModel *models.SessionModel

  if len(sessionModels) == 0 {
    sessionModel = new(models.SessionModel)
  } else {
    sessionModel = sessionModels[0]
  }

  sessionModel.Values=string(b)
  sessionModel.SessionID =session.ID

  err = queries.Update( &sessionModel )

  if err != nil {
    return err
  }

  return nil
}

func (sqliteStore *SqliteStore) delete(session *gsessions.Session) error {
  return queries.DeleteSession( session.ID )
}

func (sqliteStore *SqliteStore) Options(options sessions.Options) {
  sqliteStore.options = &gsessions.Options{
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
type GOBSerializer struct{}

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
