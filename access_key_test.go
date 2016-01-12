package main

import (
	"code.google.com/p/go-uuid/uuid"
	"encoding/json"
	"testing"
)

func TestLoadKey(t *testing.T) {
	accessKey := &AccessKey{
		Read: []string{"test-channel"},
		Write: []string{},
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

	if len(loadedKey.Write) != 0 {
		t.Error("invalid number of permissions retrieved")
		return
	}

	if len(loadedKey.Read) != 1 {
		t.Error("unexpected number of readable channels")
		return
		
	}

	if loadedKey.Read[0] != "test-channel" {
		t.Error("unexpected readable channels")
		return
		
	}

	if loadedKey.CanWrite("test-channel") {
		t.Error("key should not be able to write to channel")
		return
	}
}
