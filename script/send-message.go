package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/pote/redisurl"
	"os"
)

type Message struct {
	UUID      string `json:"id,omitempty"`
	IssuerID  string `json:"issuer,omitempty"`
	Channel   string `json:"channel,omitempty"`
	Data      string `json:"data,omitempty"`
	Event     string `json:"event,omitempty"`
}


func main() {
	data := flag.String("message", "A message from the Philotic Network!", "message #data payload")
	channel := flag.String("channel", "test-channel", "the channel the message will be broadcasted to")
	flag.Parse()


	message := &Message{
		Channel: *channel,
		Data: *data,
		Event: "message",
		IssuerID: "script/send-message.go",
	}

	payload, err := json.Marshal(message); if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	r, err := redisurl.Connect(); if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer r.Close()

	clients, err := redis.Int64(r.Do("PUBLISH", "philote:channel:" + *channel, string(payload)))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Messaged received by %v clients\n", clients)
}
