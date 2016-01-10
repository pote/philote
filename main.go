package main

import (
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
	if len(segs) < 2  {
		log.Println("No token in incoming request, dropped")
		ws.Write([]byte("No token in incoming request, dropped"))
		return
	}

	socket, err := LoadSocket(segs[1], ws)
	if err != nil {
		log.Println(err.Error())
		ws.Write([]byte(err.Error()))
		return
	}

	go socket.ListenToRedis()
	go socket.ListenToSocket()

	socket.Wait()
}
