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

	if len(c) != 2 || c[0] != "chat" || c[1] != "updates" {
		t.Error("Expected channels to be ['chat', 'updates'], but was:", c)
		return
	}
}
