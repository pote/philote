package main

import (
	"fmt"
	"os"
	"testing"
	"time"

	"golang.org/x/net/websocket"
)

func TestBasicPubSub (t *testing.T) {
	done := make(chan bool)
	go RunServer(done, os.Getenv("PORT"))
	defer func() { done <- true }()

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

func TestChannelIsolation (t *testing.T) {
	done := make(chan bool)
	go RunServer(done, os.Getenv("PORT"))
	defer func() { done <- true }()

	ak1 := &AccessKey{
		Read: []string{},
		Write: []string{"another-test-channel"},
		Token: "blah",
	}
	ak1.Save()

	ak2 := &AccessKey{
		Read: []string{"test-channel"},
		Write: []string{},
		Token: "bleh",
	}
	ak2.Save()

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
		Channel: "another-test-channel",
		Data: "hey there!",
		Event: "message",
	}

	go func() {
		time.Sleep(time.Second * 1)
		websocket.JSON.Send(ws1, &originalMessage)
	}()

	select {
	case <- channel:
		t.Error("this message should not have been received")
	case <- time.After(time.Second * 3):
		// NoOp, all is well.
	}
}

func TestWritePermissionsApplied (t *testing.T) {
	done := make(chan bool)
	go RunServer(done, os.Getenv("PORT"))
	defer func() { done <- true }()

	ak1 := &AccessKey{
		Read: []string{},
		Write: []string{},
		Token: "blah",
	}

	ak1.Save()

	ak2 := &AccessKey{
		Read: []string{"test-channel"},
		Write: []string{},
		Token: "bleh",
	}

	ak2.Save()

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
	case <- channel:
		t.Error("This message should not have been received")
	case <- time.After(time.Second * 3):
		// NoOp, all is good.
	}
}
