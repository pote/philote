package main

import(
  "fmt"
  "errors"

  "github.com/dgrijalva/jwt-go"
)

type AccessKey struct {
  Read        []string `json:"read"`
  Write       []string `json:"write"`

  jwt.StandardClaims
}

func NewAccessKey(auth string) (*AccessKey, error) {
  ak := AccessKey{}

  verifyFunc := func(t *jwt.Token) (interface{}, error) {
    return Config.jwtSecret, nil
  }

  token, err := jwt.ParseWithClaims(auth, &ak, verifyFunc)

  return &ak, err

  if claims, ok := token.Claims.(*AccessKey); ok && token.Valid {
    fmt.Printf("%v %v", claims.Read, claims.Write)
  } else {
    fmt.Println(err)
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
