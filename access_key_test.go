package main

import (
	"testing"
)

func TestLoadKey(t *testing.T) {
	ak, _ := newAccessKey()
	defer ak.Delete()

	loadedKey, err := LoadKey(ak.Token)

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
