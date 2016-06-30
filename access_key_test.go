package main

import (
	"testing"

	"code.google.com/p/go-uuid/uuid"
)

func createTestAccessKey() (*AccessKey, error) {
	ak := &AccessKey{
		Read: []string{"test-channel"},
		Write: []string{},
		Token: uuid.New(),
	}

	err := ak.Save()
	return ak, err
}


func TestLoadKey(t *testing.T) {
	ak, _ := createTestAccessKey()
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

func TestConsumeUsage(t *testing.T) {
	ak, _ := createTestAccessKey()

	if ak.UsageIsLimited() {
		t.Error("AccessKey should be considered unlimited usage")
	}

	err := ak.ConsumeUsage()

	if err == nil {
		t.Error("should not be able to run #ConsumeUsage() in an unlimited AccessKey without errors")
	}


	// Trun key into limited access.
	ak.AllowedUses = 1
	ak.Save()

	if !ak.UsageIsLimited() {
		t.Error("AccessKey should be considered limited usage")
	}

	err = ak.ConsumeUsage()

	if err != nil {
		t.Error(err)
	}

	if ak.Uses != 1 {
		t.Error("AccessKey should track it's usage")
	}

	err = ak.ConsumeUsage()

	if err == nil {
		t.Error("Should not be able to consume usage on a depleted AccessKey")
	}
}
