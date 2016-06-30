package main

import (
	"encoding/json"
	"log"

	"code.google.com/p/go-uuid/uuid"
	"github.com/garyburd/redigo/redis"
	"golang.org/x/net/websocket"
)

type Socket struct {
	ID         string
	AccessKey  *AccessKey
	ws         *websocket.Conn
	done       chan bool
}

func NewSocket(ak *AccessKey, ws *websocket.Conn) (*Socket) {
	socket := &Socket{
		ws:    ws,
		done: make(chan bool),
		ID: uuid.New(),
		AccessKey: ak,
	}

	return socket
}

func (s *Socket) redisChannels() []interface{} {
	channels := make([]interface{}, len(s.AccessKey.Read))

	i := 0
	for _, channel := range s.AccessKey.Read {
		channels[i] = "philote:channel:" + channel
		i = i + 1
	}

	return channels
}

func (s *Socket) ListenToRedis() {
	rConn := redis.PubSubConn{Conn: RedisPool.Get()}
	defer rConn.Close()

	rConn.Subscribe(s.redisChannels()...)

	var (
		message *Message
		err     error
	)

	for {
		switch event := rConn.Receive().(type) {
		case redis.Message:
			err = json.Unmarshal(event.Data, &message)

			if err != nil {
				s.logMsg("[SECURITY] Redis message isn't JSON: %s", event.Data)
				continue
			}

			switch message.Event {
			case "message":
				if message.IssuerID == s.ID {
					// Message was sent by this connection, ignore.
					continue
				}

				s.logMsg("Received message from redis on '%s'", message.Channel)
				websocket.JSON.Send(s.ws, &message)
			case "close":
				if message.IssuerID == s.ID {
					rConn.PUnsubscribe(s.redisChannels()...)
					break
				}
			}
		case error:
			rConn.Close()
			rConn = redis.PubSubConn{Conn: RedisPool.Get()}
		}
	}
}

func (s *Socket) ListenToSocket() {
	for {
		message := &Message{}
		err := websocket.JSON.Receive(s.ws, &message);
		if err != nil {
			s.logMsg("Invalid client message data: %s", err.Error() )
			if err.Error() == "EOF" {
				s.disconnect()
				break
			} else {
				continue
			}
		}

		s.logMsg("Received message from socket in channel " + message.Channel)

		if s.AccessKey.CanWrite(message.Channel) {
			s.publish(message)
		} else {
			s.logMsg("Client does not have write permission for channel " + message.Channel + ", message dropped")
		}
	}
}

func (s *Socket) disconnect() {
	message := &Message{Event: "close"}
	s.publish(message)
	s.logMsg("Disconnecting from client")
	close(s.done)
}

func (s *Socket) Wait() {
	<-s.done
	s.logMsg("Disconnected")
}

func (s *Socket) publish(message *Message) error {
	conn := RedisPool.Get()
	defer conn.Close()

	message.IssuerID = s.ID
	data, err := json.Marshal(message); if err != nil {
		return err
	}

	_, err = conn.Do("PUBLISH", "philote:channel:" + message.Channel, string(data))
	return err
}

func (s *Socket) logMsg(message string, args ...interface{}) {
	log.Printf("[" + s.ID + "] " + message + "\n", args...)
}
