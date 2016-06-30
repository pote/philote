package main

import (
	"fmt"
	"os"
	"testing"
	"time"

	"golang.org/x/net/websocket"
)

func createMatchingAccessKeys() (*AccessKey, *AccessKey) {
	ak1 := &AccessKey{
		Read: []string{"test-channel"},
		Write: []string{"test-channel"},
		Token: "blah",
	}

	ak1.Save()

	ak2 := &AccessKey{
		Read: []string{"test-channel"},
		Write: []string{"test-channel"},
		Token: "bleh",
	}

	ak2.Save()

	return ak1, ak2
}

func TestBasicPubSub (t *testing.T) {
	ak1, ak2 := createMatchingAccessKeys()
	go main(); time.Sleep(time.Second) // Give it a second, will you?

	url1 := fmt.Sprintf("ws://localhost:%v/%v", os.Getenv("PORT"), ak1.Token)
	url2 := fmt.Sprintf("ws://localhost:%v/%v", os.Getenv("PORT"), ak2.Token)

	ws1, err := websocket.Dial(url1, "", "http://localhost"); if err != nil {
		t.Fatal(err)
	}

	ws2, err := websocket.Dial(url2, "", "http://localhost"); if err != nil {
		t.Fatal(err)
	}

	channel := make(chan Message, 1)

	go func() {
		receivedMsg := Message{}
		err := websocket.JSON.Receive(ws2, &receivedMsg); if err != nil {
			t.Fatal(err)
		}

		channel <- receivedMsg
	}()


	originalMessage := &Message{
		Channel: "test-channel",
		Data: "hey there!",
		Event: "message",
	}

	go func() {
		time.Sleep(time.Second * 1)
		websocket.JSON.Send(ws1, &originalMessage)
	}()

	select {
	case receivedMessage := <- channel:
		if originalMessage.Channel != receivedMessage.Channel {
			t.Error("Received message channel is incorrect")
		}

		if originalMessage.Data != receivedMessage.Data {
			t.Error("Received message Data is incorrect")
		}

		if originalMessage.Event != receivedMessage.Event {
			t.Error("Received message Event is incorrect")
		}

	case <- time.After(time.Second * 3):
		t.Error("timout reached, no message received")
	}
}
