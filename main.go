package main

import(
  "github.com/garyburd/redigo/redis"
  "github.com/pote/chronicler"
  "github.com/pote/redisurl"
  "golang.org/x/net/websocket"
  "net/http"
  "log"
  "routes"
  "runtime"
)

var RedisPool *redis.Pool = SetupRedis()

func SetupRedis() *redis.Pool{
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
  go http.ListenAndServe(":8181", nil)

  web := chronicler.NewStory()
  web.Register(&routes.Root{})

  web.Serve(":8080")
}

