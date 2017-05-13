package main

import(
  "testing"

  "github.com/dgrijalva/jwt-go"
)


func TestNewAccessKey(t *testing.T) {
  token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
      "read": []string{"test-channel"},
      "write": []string{"test-channel"},
  })

  tokenString, err := token.SignedString(Config.jwtSecret); if err != nil {
    t.Fatal(err)
  }

  ak, err := NewAccessKey(tokenString); if err != nil {
    t.Fatal(err)
  }

  if len(ak.Read) < 1 || ak.Read[0] != "test-channel" {
    t.Error("Access Key does not have proper read permissions")
  }

  if len(ak.Write) < 1 || ak.Write[0] != "test-channel" {
    t.Error("Access Key does not have proper write permissions")
  }
}
