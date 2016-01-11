package main

import (
	"code.google.com/p/go-uuid/uuid"
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	"golang.org/x/net/websocket"
	"log"
)

type Socket struct {
	Token    string          `json:"-"`
	ID       string          `json:"-"`
	Channels []string        `json:"channels"`
	ws       *websocket.Conn `json:"-"`
	done     chan bool       `json:"-"`
}

func LoadSocket(token string, ws *websocket.Conn) (*Socket, error) {
	r := RedisPool.Get()
	rawSocket, err := redis.String(r.Do("GET", "philote:token:" + token))
	r.Close()
	if err != nil {
		return &Socket{}, err
	}
	
	if rawSocket  == "" {
		return &Socket{}, InvalidSocketTokenError{"unknown token"}
	}

	socket := &Socket{
		ws:    ws,
		done: make(chan bool),
		ID: uuid.New(),
	}

	err = json.Unmarshal([]byte(rawSocket), &socket); if err != nil {
		return socket, InvalidSocketTokenError{"invalid token data: " + err.Error()}
	}

	return socket, nil
}

func (s *Socket) redisChannels() []interface{} {
	channels := make([]interface{}, len(s.Channels))

	for index, channel := range s.Channels {
		channels[index] = "philote:channel:" + channel
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

		s.logMsg("Received message from socket on '%s'", message.IssuerID)

		s.publish(message)
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

// Internal: Actual redis Pub/Sub channel to which we will emit events.
func (s *Socket) redisChannel() string {
	return "philote:channel:" + s.ID
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
