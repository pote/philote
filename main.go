package main

import (
	"code.google.com/p/go-uuid/uuid"
	"github.com/garyburd/redigo/redis"
	"github.com/pote/redisurl"
	"golang.org/x/net/websocket"
	"log"
	"net/http"
	"os"
	"runtime"
)

var RedisPool *redis.Pool = SetupRedis()

func SetupRedis() *redis.Pool {
	pool, err := redisurl.NewPool(3, 400, "240s")
	if err != nil {
		panic(err)
	}

	return pool
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.Printf("[Main] Initializing Philotic Network on %v core(s)\n", runtime.NumCPU())

	http.Handle("/", websocket.Handler(ServeWebSocket))
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}

func ServeWebSocket(ws *websocket.Conn) {
	token, err := TokenFromConn(ws)

	if err != nil {
		log.Fatal(err)
		ws.Close()
	}

	connectionId := uuid.New()

	LogMsg("Connected and listening", connectionId)

	go ReceiveMessages(connectionId, ws)

	for _, channel := range token.Channels {
		go DispatchMessages(channel, connectionId, ws)
	}
	select {}
}

type Message struct {
	UUID    string `json:"id"`
	Channel string `json:"channel"`
	Data    string `json:"data"`
}

func DispatchMessages(channel, identifier string, ws *websocket.Conn) {
	pubSub := redis.PubSubConn{Conn: RedisPool.Get()}
	pubSub.PSubscribe(channel + ":*")

	for {
		switch event := pubSub.Receive().(type) {
		case redis.PMessage:
			if event.Channel != channel+":"+identifier {
				LogMsg("Received message from redis on '%s'", identifier, channel)
				websocket.JSON.Send(ws, &Message{
					UUID:    uuid.New(),
					Channel: channel,
					Data:    string(event.Data),
				})
			}
		}
	}
}

func ReceiveMessages(identifier string, ws *websocket.Conn) {
	token, err := TokenFromConn(ws)

	if err != nil {
		log.Fatal(err)
		ws.Close()
	}

	for {
		var message *Message
		websocket.JSON.Receive(ws, &message)

		LogMsg("Received message from socket on '%s'", identifier, message.Channel)

		if token.CanAccess(message.Channel) {
			c := RedisPool.Get()
			c.Do("PUBLISH", message.Channel+":"+identifier, message.Data)
			c.Close()
		}
	}
}

func LogMsg(message, connection string, args ...interface{}) {
	args = append([]interface{}{connection}, args...)
	log.Printf("[%s] "+message+"\n", args...)
}
