package main

import (
	"code.google.com/p/go-uuid/uuid"
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	"golang.org/x/net/websocket"
	"log"
)

type Socket struct {
	ID    string
	Token *AccessToken
	done  chan bool
	ws    *websocket.Conn
	redis redis.PubSubConn
}

type Message struct {
	UUID    string `json:"id,omitempty"`
	Channel string `json:"channel,omitempty"`
	Data    string `json:"data,omitempty"`
	Event   string `json:"event,omitempty"`
}

func NewSocket(ws *websocket.Conn) (socket *Socket, err error) {
	token, err := ParseAccessToken(ws.Request().FormValue("token"))

	if err != nil {
		logMsg("[FATAL] Can't parse connection token: %s", "new", err)
		ws.Close()
		return
	}

	socket = &Socket{
		ws:    ws,
		done:  make(chan bool),
		ID:    uuid.New(),
		Token: token,
	}

	logMsg("Connecting...", socket.ID)

	return
}

func (s *Socket) ListenToRedis() {
	s.redis = redis.PubSubConn{Conn: RedisPool.Get()}
	defer s.redis.Close()

	s.redis.PSubscribe(s.redisPattern())

	var (
		message *Message
		err     error
	)

	for {
		switch event := s.redis.Receive().(type) {
		case redis.PMessage:
			err = json.Unmarshal(event.Data, &message)

			if err != nil {
				logMsg("[SECURITY] Redis message isn't JSON: %s", s.ID, event.Data)
				continue
			}

			switch message.Event {
			case "message":
				if event.Channel == s.redisChannel() {
					// Message was sent by this connection, ignore.
					continue
				}

				logMsg("Received message from redis on '%s'", s.ID, message.Channel)
				websocket.JSON.Send(s.ws, &message)
			case "close":
				if event.Channel == s.redisChannel() {
					s.redis.PUnsubscribe(s.redisPattern())
					break
				}
			}
		}
	}
}

func (s *Socket) ListenToSocket() {
	for {
		var (
			data    []byte
			message *Message
		)

		err := websocket.Message.Receive(s.ws, &data)

		if err != nil {
			s.disconnect()
			break
		}

		err = json.Unmarshal(data, &message)

		if err != nil {
			logMsg("[SECURITY] Invalid client message: %s", s.ID, data)
			s.disconnect()
			continue
		}

		logMsg("Received message from socket on '%s'", s.ID, message.Channel)

		if s.Token.CanAccess(message.Channel) {
			c := RedisPool.Get()
			c.Do("PUBLISH", s.redisChannel(), data)
			c.Close()
		}
	}
}

func (s *Socket) disconnect() {
	message := &Message{Event: "close"}
	data, _ := json.Marshal(message)

	c := RedisPool.Get()
	c.Do("PUBLISH", s.redisChannel(), data)
	c.Close()

	logMsg("Disconnecting from client", s.ID)
	close(s.done)
}

func (s *Socket) Wait() {
	<-s.done
	logMsg("Disconnected", s.ID)
}

// Internal: Actual redis Pub/Sub channel to which we will emit events.
func (s *Socket) redisChannel() string {
	return s.Token.Hub + ":" + s.ID
}

// Internal: Pattern to PSUBSCRIBE to in redis events.
func (s *Socket) redisPattern() string {
	return s.Token.Hub + ":*"
}

func logMsg(message, connection string, args ...interface{}) {
	args = append([]interface{}{connection}, args...)
	log.Printf("[%s] "+message+"\n", args...)
}
