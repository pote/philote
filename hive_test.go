package main

import(
  "net/http"
  "net/http/httptest"
  "net/url"
  "testing"
  "time"

  "github.com/dgrijalva/jwt-go"
  "github.com/gorilla/websocket"
)

func TestHiveSuccessfulPhiloteRegistration(t *testing.T) {
  h := NewHive()
  if len(h.Philotes) != 0 {
    t.Error("new Hive shouldn't have registered philotes")
  }

  token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
      "read": []string{"test-channel"},
      "write": []string{"test-channel"},
  })

  tokenString, err := token.SignedString(Config.jwtSecret); if err != nil {
    t.Fatal(err)
  }

  server := httptest.NewServer(http.HandlerFunc(h.ServeNewConnection))
  header := map[string][]string{
    "Authorization": []string{"Bearer " + tokenString},
  }
  u, _ := url.Parse(server.URL)
  u.Scheme = "ws"
  _, _, err = websocket.DefaultDialer.Dial(u.String(), header); if err != nil {
    t.Error(err)
  }

  if len(h.Philotes) != 1 {
    t.Error("philote should  be registered on successful auth")
  }
}

func TestHiveIncorrectAuth(t *testing.T) {
  h := NewHive()
  if len(h.Philotes) != 0 {
    t.Error("new Hive shouldn't have registered philotes")
  }

  server := httptest.NewServer(http.HandlerFunc(h.ServeNewConnection))
  header := map[string][]string{
    "Authorization": []string{"Bearer " + "foo"},
  }
  u, _ := url.Parse(server.URL)
  u.Scheme = "ws"
  _, _, err := websocket.DefaultDialer.Dial(u.String(), header); if err == nil {
    t.Error("The Dial action should fail when there is no auth token")
  }

  if len(h.Philotes) != 0 {
    t.Error("philote should not be registered when missing auth")
  }
}

func TestHivePhiloteRegistrationWithNoAuth(t *testing.T) {
  h := NewHive()
  if len(h.Philotes) != 0 {
    t.Error("new Hive shouldn't have registered philotes")
  }

  server := httptest.NewServer(http.HandlerFunc(h.ServeNewConnection))
  u, _ := url.Parse(server.URL)
  u.Scheme = "ws"
  _, _, err := websocket.DefaultDialer.Dial(u.String(), nil); if err == nil {
    t.Error("The Dial action should fail when there is no auth token")
  }

  if len(h.Philotes) != 0 {
    t.Error("philote should not be registered when missing auth")
  }
}

func TestHiveDeregisterPhilote(t *testing.T) {
  h := NewHive()
  token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
      "read": []string{"test-channel"},
      "write": []string{"test-channel"},
  })
  tokenString, err := token.SignedString(Config.jwtSecret); if err != nil {
    t.Fatal(err)
  }
  server := httptest.NewServer(http.HandlerFunc(h.ServeNewConnection))
  header := map[string][]string{
    "Authorization": []string{"Bearer " + tokenString},
  }
  u, _ := url.Parse(server.URL)
  u.Scheme = "ws"
  conn, _, err := websocket.DefaultDialer.Dial(u.String(), header); if err != nil {
    t.Error(err)
  }

  if len(h.Philotes) != 1 {
    t.Error("philote should  be registered on successful auth")
  }

  conn.Close()
  time.Sleep(time.Second)

  if len(h.Philotes) != 0 {
    t.Error("Disconnected Philotes should be automatically deregistered")
  }
}
