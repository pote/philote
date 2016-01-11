package main

import (
	"code.google.com/p/go-uuid/uuid"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/pote/redisurl"
	"os"
	"strings"
)

type Socket struct {
	Channels map[string]string `json:"channels"`
}

func main() {
	token := flag.String("token", uuid.New(), "authorization token")
	channels := flag.String("channels", "test-channel", "comma-separated list of channels")
	flag.Parse()

	socket := &Socket{Channels: map[string]string{}}
	for _, channel := range strings.Split(*channels, ",") {
		socket.Channels[channel] = "read,write"
	}
	data, err := json.Marshal(socket)

	r, err := redisurl.Connect(); if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer r.Close()

	_, err = r.Do("SET", "philote:token:" + *token, string(data)); if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(*token)
}
