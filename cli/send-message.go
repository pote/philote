package main

import (
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	"github.com/pote/redisurl"
)

type Message struct {
	UUID      string `json:"id,omitempty"`
	IssuerID  string `json:"issuer,omitempty"`
	Channel   string `json:"channel,omitempty"`
	Data      string `json:"data,omitempty"`
	Event     string `json:"event,omitempty"`
}


func publishMessage(channel, data string) (listeners int, err error) {
	message := &Message{
		Channel: channel,
		Data: data,
		Event: "message",
		IssuerID: "philote-cli",
	}

	payload, err := json.Marshal(message); if err != nil {
		return
	}

	r, err := redisurl.Connect(); if err != nil {
		return
	}
	defer r.Close()

	listeners, err = redis.Int(r.Do("PUBLISH", "philote:channel:" + channel, string(payload)))
	if err != nil {
		return
	}

	return
}
