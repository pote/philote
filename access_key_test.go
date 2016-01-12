package main

import (
	"code.google.com/p/go-uuid/uuid"
	"encoding/json"
	"testing"
)

func TestLoadKey(t *testing.T) {
	accessKey := &AccessKey{
		Channels: map[string]string{
			"test-channel": "read,write",
		},
		Token: uuid.New(),
	}

	data, _ := json.Marshal(accessKey)

	r := RedisPool.Get()
	r.Do("SET", "philote:token:" + accessKey.Token, string(data))

	loadedKey, err := LoadKey(accessKey.Token)

	if err != nil {
		t.Error(err)
		return
	}

	if len(loadedKey.Channels) != 1 {
		t.Error("invalid number of channels retrieved")
		return
	}

	if loadedKey.Channels["test-channel"] != "read,write" {
		t.Errorf("unexpected channel name: %+v", loadedKey.Channels)
		return
	}
}
