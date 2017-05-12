package main

import (
  "errors"
  "fmt"
  "github.com/dgrijalva/jwt-go"
)

type AccessKey struct {
  Read        []string `json:"read"`
  Write       []string `json:"write"`

  jwt.StandardClaims
}

func NewAccessKey(auth string) (*AccessKey, error) {
  ak := AccessKey{}
  _, err := jwt.ParseWithClaims(auth, &ak, func(t *jwt.Token) (interface{}, error) {
    return []byte{}, nil
  })

  return &ak, err
}

func LoadKey(rawToken string) (*AccessKey, error) {
  ak := &AccessKey{}

  keyFunc := func(token *jwt.Token) (interface{}, error) {
    return []byte("AllYourBase"), nil
  }

  token, err := jwt.ParseWithClaims(rawToken, &AccessKey{}, keyFunc); if err != nil {
    return ak, err
  }

  if claims, ok := token.Claims.(*AccessKey); ok && token.Valid {
    fmt.Printf("%v %v", claims.Read, claims.Write)
  } else {
    fmt.Println(err)
    return ak, errors.New("invalid token")
  }

	return ak, nil
}

func (ak *AccessKey) CanWrite(channel string) bool {
	for _, c := range ak.Write {
		if c == channel {
			return true
		}
	}
	
	return false
}
