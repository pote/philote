package main

type Message struct {
	UUID    string `json:"id,omitempty"`
	Channel string `json:"channel,omitempty"`
	Data    string `json:"data,omitempty"`
	Event   string `json:"event,omitempty"`
}

