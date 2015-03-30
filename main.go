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
	if _, _, err := RoutingInfo(getRequestToken(ws.Request())); err != nil {
		log.Fatal(err)
		ws.Close()
	}
	identifier := uuid.New()

	go ReceiveMessages(identifier, ws)
	go DispatchMessages(identifier, ws)
	select {}
}

func DispatchMessages(identifier string, ws *websocket.Conn) {
	_, channel, _ := RoutingInfo(getRequestToken(ws.Request()))
	pubSub := redis.PubSubConn{Conn: RedisPool.Get()}
	pubSub.PSubscribe(channel + ":*")

	for {
		switch m := pubSub.Receive().(type) {
		case redis.PMessage:
			if m.Channel != channel+":"+identifier {
				websocket.Message.Send(ws, string(m.Data))
			}
		}
	}
}

func ReceiveMessages(identifier string, ws *websocket.Conn) {
	_, channel, _ := RoutingInfo(getRequestToken(ws.Request()))

	for {
		var message string
		websocket.Message.Receive(ws, &message)

		c := RedisPool.Get()
		c.Do("PUBLISH", channel+":"+identifier, message)
		c.Close()
	}
}

func getRequestToken(req *http.Request) (token string) {
	if tokens, ok := req.Form["token"]; ok {
		return tokens[len(tokens)-1]
	} else {
		return ""
	}
}

// FIXME: This should return multiple channels, so that we start listening on
// all of them.
func RoutingInfo(at string) (hub, channel string, err error) {
	token, err := ParseAccessToken(at)

	if err != nil {
		return
	}

	hub = token.Hub
	channel = token.Channels[0]

	return
}
