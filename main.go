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
	"strings"
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
	http.ListenAndServe(":" + os.Getenv("PORT"), websocket.Handler(ServeWebSocket))
}

func ServeWebSocket(ws *websocket.Conn) {
	segs := strings.Split(ws.Request().URL.Path, "/")
	if len(segs) < 2 || segs[1] == "" {
		log.Println("No token in incoming request, dropped")
		ws.Write([]byte("No token in incoming request, dropped"))
		return
	}

	r := RedisPool.Get()
	rawSocket, err := redis.String(r.Do("GET", segs[1]))
	r.Close()
	if err != nil {
		log.Println(err.Error())
		ws.Write([]byte("No token in incoming request, dropped"))
		return
	}
	
	if rawSocket  == "" {
		log.Println("Unknown or expired token")
		ws.Write([]byte("Unknown or expired token"))
		return

	}

	socket := &Socket{
		ws:    ws,
		done: make(chan bool),
		ID: uuid.New(),
	}

	json.Unmarshal([]byte(rawSocket), &socket)

	go socket.ListenToRedis()
	go socket.ListenToSocket()

	socket.Wait()
}
