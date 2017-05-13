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

  tokenString, err := token.SignedString(Hive.jwtSecret); if err != nil {
    t.Fatal(err)
  }

  ak, err := Hive.NewAccessKey(tokenString); if err != nil {
    t.Fatal(err)
  }

  if ak.Read[0] != "test-channel" {
    t.Error("Access Key does not have proper read permissions")
  }

  if ak.Write[0] != "test-channel" {
    t.Error("Access Key does not have proper write permissions")
  }
}
