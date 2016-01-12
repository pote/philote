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
	accessKey := &AccessKey{
		Read: []string{"test-channel"},
		Write: []string{},
		Token: uuid.New(),
	}

	data, _ := json.Marshal(accessKey)

	r := RedisPool.Get()
	r.Do("SET", "philote:token:" + accessKey.Token, string(data))

	// Test authorization against the real thing.
	go main(); time.Sleep(time.Second) // Give it a second, will you?

	wsConnectionURL := fmt.Sprintf("ws://localhost:%v/%v", os.Getenv("PORT"), accessKey.Token)
	ws, err := websocket.Dial(wsConnectionURL, "", "http://localhost")

	if err != nil {
		t.Error(err)
	}

	if !ws.IsClientConn() {
		t.Error("created connection should be considered a client one")
	}
  time.Sleep(time.Second) // Give it a second, will you?

	r.Do("DEL", "philote" + accessKey.Token)
	r.Close()
}
