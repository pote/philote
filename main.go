package main

import (
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/garyburd/redigo/redis"
	"github.com/pote/redisurl"
	"golang.org/x/net/websocket"

	"lua"
)

var RedisPool *redis.Pool = SetupRedis()
var Lua *lua.Lua = lua.NewClient(RedisPool)

func SetupRedis() *redis.Pool {
	var err error
	var maxConnections int
	mcs := os.Getenv("REDIS_MAX_CONNECTIONS")
	if mcs != "" {
	 maxConnections, err = strconv.Atoi(mcs); if err != nil {
		 log.Fatalf(err.Error())
	 }
	} else {
		maxConnections = 400
	}

	redisURL := os.Getenv("REDIS_URL"); if redisURL == "" {
		redisURL = "redis://localhost:6379"
	}

	pool, err := redisurl.NewPoolWithURL(redisURL, 3, maxConnections, "240s")
	if err != nil {
		panic(err)
	}

	return pool
}

func main() {
	port := os.Getenv("PORT"); if port == "" {
		port = "6380"
	}
	runtime.GOMAXPROCS(runtime.NumCPU())

	log.Printf("[Main] Initializing Philotic Network\n")
	log.Printf("[Main] Version: %v\n", VERSION)
	log.Printf("[Main] Port: %v\n", port)
	log.Printf("[Main] Cores: %v\n", runtime.NumCPU())

	done := make(chan bool)
	RunServer(done, port)
}

func RunServer(done chan bool, port string) {
	go func() {
		err := http.ListenAndServe(":" + port, websocket.Handler(ServeWebSocket)); if err != nil {
			log.Println(err)
		}
	}()

	<- done
	log.Println("[Main] Stop signal detected, shutting down.")
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
