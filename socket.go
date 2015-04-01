package main

import (
	"code.google.com/p/go-uuid/uuid"
	"encoding/json"
	"errors"
	"github.com/garyburd/redigo/redis"
	"golang.org/x/net/websocket"
	"log"
)

type Socket struct {
	ID    string
	Token *AccessToken
	done  chan bool
	ws    *websocket.Conn
}

type Message struct {
	UUID    string `json:"id,omitempty"`
	Channel string `json:"channel,omitempty"`
	Data    string `json:"data,omitempty"`
	Event   string `json:"event,omitempty"`
}

func NewSocket(ws *websocket.Conn) (socket *Socket, err error) {
	tokenString := ws.Request().FormValue("token")

	if tokenString == "" {
		err = errors.New("Need a `token` query parameter when connecting")
		logMsg("[FATAL] Connection didn't send token param", "new")
		ws.Close()
		return
	}

	var token *AccessToken
	token, err = ParseAccessToken(tokenString)

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
			s.redisPub(data)
		}
	}
}

func (s *Socket) disconnect() {
	message := &Message{Event: "close"}
	data, _ := json.Marshal(message)

	s.redisPub(data)

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

func (s *Socket) redisPub(data interface{}) {
	conn := RedisPool.Get()
	conn.Do("PUBLISH", s.redisChannel(), data)
	conn.Close()
}

func logMsg(message, connection string, args ...interface{}) {
	args = append([]interface{}{connection}, args...)
	log.Printf("[%s] "+message+"\n", args...)
}
