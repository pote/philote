package main

import "testing"

func TestRoutingInfo(t *testing.T) {
	var token = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJodWIiOiJkZWFkYmVlZiIsImNoYW5uZWxzIjpbImNoYXQiLCJ1cGRhdGVzIl19.vlqHRxfxSidyH9_oW-rMl_LvLR8UqhK5uGc5KRjTxl0"

	h, c, e := RoutingInfo(token)

	if e != nil {
		t.Error(e)
		return
	}

	if h != "deadbeef" {
		t.Error("Expected hub to be 'deadbeef', was", h)
		return
	}

	// FIXME: c should actually be ['chat', 'updates'], but we're ignoring
	// that it's an array for now.
	if c != "chat" {
		t.Error("Expected channel to be 'chat', but was:", c)
		return
	}
}
