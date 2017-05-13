package main

import(
  "github.com/dgrijalva/jwt-go"
)

type AccessKey struct {
  Read        []string `json:"read"`
  Write       []string `json:"write"`

  jwt.StandardClaims
}


func (ak *AccessKey) CanWrite(channel string) bool {
  for _, c := range ak.Write {
    if c == channel {
      return true
    }
  }

  return false
}
