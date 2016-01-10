package main

import (
	"code.google.com/p/go-uuid/uuid"
	"encoding/json"
	"golang.org/x/net/websocket"
	"testing"
)

func TestLoadSocket(t *testing.T) {
	r := RedisPool.Get()
	defer r.Close()

	socket := &Socket{
		Channels: []string{"test-channel"},
	}

	token := uuid.New()
	data, _ := json.Marshal(socket)
	r.Do("SET", "philote:token:" + token, string(data))

	loadedSocket, err := LoadSocket(token, &websocket.Conn{})

	if err != nil {
		t.Error(err)
		return
	}

	if len(loadedSocket.Channels) != 1 {
		t.Error("invalid number of channels retrieved")
		return
	}

	if loadedSocket.Channels[0] != "test-channel" {
		t.Error("unexpected channel name: " + loadedSocket.Channels[0])
		return
	}
}
