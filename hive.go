package main

import(
  "errors"
  "fmt"

  "github.com/dgrijalva/jwt-go"
)

type hive struct {
  Philotes []*Philote
  NewPhilotes chan *Philote
  jwtSecret  []byte
}

func NewHive() (*hive) {
  h := &hive{Philotes: []*Philote{}}

  go h.RegisterNewPhilotes()

  return h
}

func (h *hive) RegisterNewPhilotes() {
  for {
    philote := <- h.NewPhilotes

    philote.Hive = h
    h.Philotes = append(h.Philotes, philote)
  }
}

func (h *hive) NewAccessKey(auth string) (*AccessKey, error) {
  ak := AccessKey{}

  verifyFunc := func(t *jwt.Token) (interface{}, error) {
    return h.jwtSecret, nil
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
