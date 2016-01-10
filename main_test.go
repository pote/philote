package main

import (
	"code.google.com/p/go-uuid/uuid"
	"encoding/json"
	"fmt"
	"golang.org/x/net/websocket"
	"os"
	"testing"
)

func TestBasicAuthorization (t *testing.T) {
	socket := &Socket{
		Channels: []string{"test"},
	}

	token := uuid.New()
	data, _ := json.Marshal(socket)

	r := RedisPool.Get()
	defer r.Close()

	r.Do("SET", "philote:" + token, string(data))
	defer r.Do("DEL", "philote" + token)

	// Test authorization against the real thing.
	go main()

	wsConnectionURL := fmt.Sprintf("ws://localhost:%v/%v", os.Getenv("PORT"), token)
	_, err := websocket.Dial(wsConnectionURL, "", "http://localhost")
	if err != nil {
		t.Error(err)
	}

}
