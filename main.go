package main

import (
	"code.google.com/p/go-uuid/uuid"
	"encoding/json"
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

	done := make(chan bool)

	go ReceiveMessages(connectionId, ws, done)
	go DispatchMessages(token.Hub, connectionId, ws)

	<-done
	LogMsg("Disconnected", connectionId)
}

type Message struct {
	UUID    string `json:"id"`
	Channel string `json:"channel"`
	Data    string `json:"data"`
	Event   string `json:"event,omitempty"`
}

func DispatchMessages(hub, identifier string, ws *websocket.Conn) {
	pubSub := redis.PubSubConn{Conn: RedisPool.Get()}
	defer pubSub.Close()

	pubSub.PSubscribe(hub + ":*")

	var message *Message
	var err error

	for {
		switch event := pubSub.Receive().(type) {
		case redis.PMessage:
			err = json.Unmarshal(event.Data, &message)

			if err != nil {
				continue
			}

			switch message.Event {
			case "message":
				if event.Channel == hub+":"+identifier {
					continue
				}

				LogMsg("Received message from redis on '%s'", identifier, hub)
				websocket.JSON.Send(ws, &message)
			case "close":
				if event.Channel == hub+":"+identifier {
					pubSub.PUnsubscribe(hub + ":*")
					break
				}
			default:
				log.Println(message)
			}
		}
	}
}

func ReceiveMessages(identifier string, ws *websocket.Conn, done chan bool) {
	token, err := TokenFromConn(ws)

	if err != nil {
		log.Fatal(err)
		ws.Close()
	}

	for {
		var data []byte
		var message *Message

		err = websocket.Message.Receive(ws, &data)

		if err != nil {
			message = &Message{Event: "close"}
			data, err = json.Marshal(message)

			c := RedisPool.Get()
			c.Do("PUBLISH", token.Hub+":"+identifier, data)
			c.Close()

			LogMsg("Received client disconnection", identifier)

			done <- true

			return
		}

		err = json.Unmarshal(data, &message)

		message.Event = "message"

		if err != nil {
			LogMsg("Client message with wrong format", identifier)
			continue
		}

		LogMsg("Received message from socket on '%s'", identifier, message.Channel)

		if token.CanAccess(message.Channel) {
			c := RedisPool.Get()
			c.Do("PUBLISH", token.Hub+":"+identifier, data)
			c.Close()
		}
	}
}

func LogMsg(message, connection string, args ...interface{}) {
	args = append([]interface{}{connection}, args...)
	log.Printf("[%s] "+message+"\n", args...)
}
