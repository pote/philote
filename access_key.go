package main

import(
  "errors"

  "github.com/dgrijalva/jwt-go"
)

type AccessKey struct {
  Read        []string `json:"read"`
  Write       []string `json:"write"`
  API         bool     `json:"api"`

  jwt.StandardClaims
}

func NewAccessKey(auth string) (*AccessKey, error) {
  ak := AccessKey{}

  verifyFunc := func(t *jwt.Token) (interface{}, error) {
    return Config.jwtSecret, nil
  }

  token, err := jwt.ParseWithClaims(auth, &ak, verifyFunc); if err != nil {
    return &ak, err
  }

  if claims, ok := token.Claims.(*AccessKey); ok && token.Valid {
    ak.Read = claims.Read
    ak.Write = claims.Write
    ak.API = !!claims.API
    return &ak, nil
  } else {
    return &ak, errors.New("invalid token")
  }

  return &ak, nil
}

func (ak *AccessKey) CanWrite(channel string) bool {
  for _, c := range ak.Write {
    if c == channel {
      return true
    }
  }

  return false
}

func (ak *AccessKey) CanRead(channel string) bool {
  for _, c := range ak.Read {
    if c == channel {
      return true
    }
  }

  return false
}
