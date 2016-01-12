package main

import (
	"code.google.com/p/go-uuid/uuid"
	"fmt"
	"golang.org/x/net/websocket"
	"os"
	"testing"
	"time"
)

func newAccessKey() (*AccessKey, error) {
	ak := &AccessKey{
		Read: []string{"test-channel"},
		Write: []string{},
		Token: uuid.New(),
	}

	err := ak.Save()
	return ak, err
}

func TestBasicAuthorization (t *testing.T) {
	ak, _ := newAccessKey()
	defer ak.Delete()

	// Test authorization against the real thing.
	go main(); time.Sleep(time.Second) // Give it a second, will you?

	wsConnectionURL := fmt.Sprintf("ws://localhost:%v/%v", os.Getenv("PORT"), ak.Token)
	ws, err := websocket.Dial(wsConnectionURL, "", "http://localhost")

	if err != nil {
		t.Error(err)
	}

  time.Sleep(time.Second) // Give it a second, will you?
	if !ws.IsClientConn() {
		t.Error("created connection should be considered a client one")
	}
}
