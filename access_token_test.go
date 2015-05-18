package main

import "testing"

func TestParseAccessToken(t *testing.T) {
	redis := RedisPool.Get()
	defer redis.Close()

	redis.Do("SET", "hubs:deadbeef", "deadbeefsecret")

	var tokenString = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJodWIiOiJkZWFkYmVlZiIsImNoYW5uZWxzIjpbImNoYXQiLCJ1cGRhdGVzIl19.vlqHRxfxSidyH9_oW-rMl_LvLR8UqhK5uGc5KRjTxl0"

	token, err := ParseAccessToken(tokenString)

	redis.Do("DEL", "hubs:deadbeef")

	if err != nil {
		t.Error(err)
		return
	}

	if token.Hub != "deadbeef" {
		t.Error("Expected hub to be deadbeef, but got:", token.Hub)
		return
	}

	if token.Channels[0] != "chat" || token.Channels[1] != "updates" {
		t.Error("Expected channels to be [chat, updates], but was:", token.Channels)
		return
	}

	return
}
