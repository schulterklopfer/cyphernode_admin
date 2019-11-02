package cnaOIDC

/**
 based on https://github.com/markbates/goth. Thanks a lot! :D
**/

import (
  "encoding/json"
  "errors"
  "golang.org/x/oauth2"
  "net/url"
  "strings"
  "time"
)

// Session stores data during the auth process with the OpenID Connect flow.
type Session struct {
  AuthURL      string
  AccessToken  string
  RefreshToken string
  ExpiresAt    time.Time
  IDToken      string
}

// GetAuthURL will return the URL set by calling the `BeginAuth` function on the OpenID Connect flow.
func (s Session) GetAuthURL() (string, error) {
  if s.AuthURL == "" {
    return "", errors.New("an AuthURL has not be set")
  }
  return s.AuthURL, nil
}

// Authorize the session with the OpenID Connect flow and return the access token to be stored for future use.
func (s *Session) Authorize(provider *Flow, params url.Values) (string, error) {
  p := provider
  token, err := p.config.Exchange(oauth2.NoContext, params.Get("code"))
  if err != nil {
    return "", err
  }

  if !token.Valid() {
    return "", errors.New("Invalid token received from flow")
  }

  s.AccessToken = token.AccessToken
  s.RefreshToken = token.RefreshToken
  s.ExpiresAt = token.Expiry
  s.IDToken = token.Extra("id_token").(string)
  return token.AccessToken, err
}

// Marshal the session into a string
func (s Session) Marshal() string {
  b, _ := json.Marshal(s)
  return string(b)
}

func (s Session) String() string {
  return s.Marshal()
}

// UnmarshalSession will unmarshal a JSON string into a session.
func (flow *Flow) UnmarshalSession(data string) (*Session, error) {
  sess := &Session{}
  err := json.NewDecoder(strings.NewReader(data)).Decode(sess)
  return sess, err
}