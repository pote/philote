package main

import (
	"github.com/garyburd/redigo/redis"
	"github.com/pote/redisurl"
	"golang.org/x/net/websocket"
	"log"
	"lua"
	"net/http"
	"os"
	"runtime"
	"strings"
)

var RedisPool *redis.Pool = SetupRedis()
var Lua *lua.Lua = lua.NewClient(RedisPool)

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
		websocket.JSON.Send(ws, "No token in incoming request, dropped")
		return
	}

	ak, err := LoadKey(segs[1]); if err != nil {
		log.Println(err.Error())
		websocket.JSON.Send(ws, err.Error())
		return
	}

	if ak.UsageIsLimited() {
	err = ak.ConsumeUsage(); if err != nil {
			log.Println(err.Error())
			websocket.JSON.Send(ws, err.Error())
			return
		}
	}

	socket := NewSocket(ak, ws)
	go socket.ListenToRedis()
	go socket.ListenToSocket()

	socket.Wait()
}
