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
		Channels: map[string]string{
			"test-channel": "read,write",
		},
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

	if loadedSocket.Channels["test-channel"] != "read,write" {
		t.Errorf("unexpected channel name: %+v", loadedSocket.Channels)
		return
	}
}
