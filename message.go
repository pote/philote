package main

type Message struct {
	IssuerID  string `json:"issuer,omitempty"`
	Channel   string `json:"channel,omitempty"`
	Data      string `json:"data,omitempty"`
	Event     string `json:"event,omitempty"`
}
