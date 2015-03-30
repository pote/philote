package main

import "testing"

func TestParseAccessToken(t *testing.T) {
	var tokenString = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJodWIiOiJkZWFkYmVlZiIsImNoYW5uZWxzIjpbImNoYXQiLCJ1cGRhdGVzIl19.vlqHRxfxSidyH9_oW-rMl_LvLR8UqhK5uGc5KRjTxl0"

	token, err := ParseAccessToken(tokenString)

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
