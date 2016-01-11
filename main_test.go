package main

import (
	"code.google.com/p/go-uuid/uuid"
	"encoding/json"
	"fmt"
	"golang.org/x/net/websocket"
	"os"
	"testing"
	"time"
)

func TestBasicAuthorization (t *testing.T) {
	socket := &Socket{
		Channels: map[string]string{
			"test-channel": "read,write",
		},
	}

	token := uuid.New()
	data, _ := json.Marshal(socket)

	r := RedisPool.Get()
	r.Do("SET", "philote:token:" + token, string(data))

	// Test authorization against the real thing.
	go main(); time.Sleep(time.Second) // Give it a second, will you?

	wsConnectionURL := fmt.Sprintf("ws://localhost:%v/%v", os.Getenv("PORT"), token)
	ws, err := websocket.Dial(wsConnectionURL, "", "http://localhost")
	if err != nil {
		t.Error(err)
	}

	if !ws.IsClientConn() {
		t.Error("created connection should be considered a client one")
	}
  time.Sleep(5 * time.Second) // Give it a second, will you?

	r.Do("DEL", "philote" + token)
	r.Close()
}
