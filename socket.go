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

func (s *Socket) ListenToRedis() {
	rConn := redis.PubSubConn{Conn: RedisPool.Get()}
	defer rConn.Close()

	rConn.PSubscribe(s.redisPattern())

	var (
		message *Message
		err     error
	)

	for {
		switch event := rConn.Receive().(type) {
		case redis.PMessage:
			err = json.Unmarshal(event.Data, &message)

			if err != nil {
				s.logMsg("[SECURITY] Redis message isn't JSON: %s", event.Data)
				continue
			}

			switch message.Event {
			case "message":
				if event.Channel == s.redisChannel() {
					// Message was sent by this connection, ignore.
					continue
				}

				s.logMsg("Received message from redis on '%s'", message.Channel)
				websocket.JSON.Send(s.ws, &message)
			case "close":
				if event.Channel == s.redisChannel() {
					rConn.PUnsubscribe(s.redisPattern())
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
		var (
			data    []byte
			message *Message
		)

		err := websocket.Message.Receive(s.ws, &data); if err != nil {
			s.disconnect()
			break
		}

		err = json.Unmarshal(data, &message); if err != nil {
			s.logMsg("[SECURITY] Invalid client message: %s", data)
			s.disconnect()
			continue
		}

		s.logMsg("Received message from socket on '%s'", message.Channel)

		s.redisPub(data)
	}
}

func (s *Socket) disconnect() {
	message := &Message{Event: "close"}
	data, _ := json.Marshal(message)

	s.redisPub(data)

	s.logMsg("Disconnecting from client")
	close(s.done)
}

func (s *Socket) Wait() {
	<-s.done
	s.logMsg("Disconnected", s.ID)
}

// Internal: Actual redis Pub/Sub channel to which we will emit events.
func (s *Socket) redisChannel() string {
	return "philote:channel:" + s.ID
}

// Internal: Pattern to PSUBSCRIBE to in redis events.
func (s *Socket) redisPattern() string {
	return "philote:channel:*"
}

func (s *Socket) redisPub(data []byte) {
	conn := RedisPool.Get()
	conn.Do("PUBLISH", s.redisChannel(), data)
	conn.Close()
}

func (s *Socket) logMsg(message string, args ...interface{}) {
	log.Printf("[" + s.ID + "] " + message + "\n", args...)
}
