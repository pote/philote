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

func TestPhilotesExchangingMessages(t *testing.T) {
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
  conn1, _, err := websocket.DefaultDialer.Dial(u.String(), header); if err != nil {
    t.Error(err)
  }
  conn2, _, err := websocket.DefaultDialer.Dial(u.String(), header); if err != nil {
    t.Error(err)
  }

  if len(h.Philotes) != 2 {
    t.Error("Both philotes should be connected and registered")
  }
  originalMessage := &Message{Event: "message", Data: "yo!", Channel: "test-channel"}

  go func() { time.Sleep(time.Second); conn1.WriteJSON(originalMessage) }()

  receivedMessage := &Message{}
  err = conn2.ReadJSON(receivedMessage); if err != nil {
    t.Error(err)
  }

  if receivedMessage.Data != "yo!" {
    t.Error("incorrect message data")
  }
  if receivedMessage.Event != "message" {
    t.Error("incorrect message data")
  }
}

