package main

import(
  "errors"
  "github.com/garyburd/redigo/redis"
  //"golang.org/x/go-uuid/uuid"
  "code.google.com/p/go-uuid/uuid"
  "golang.org/x/net/websocket"
  "log"
  "strings"
)

func ServeWebSocket(ws *websocket.Conn) {
  if _, _, err := RoutingInfo(ws.Request().URL.Path); err != nil {
    log.Fatal(err)
    ws.Close()
  }
  identifier := uuid.New()

  go ReceiveMessages(identifier, ws)
  go DispatchMessages(identifier, ws)
  select{}
}

func DispatchMessages(identifier string, ws *websocket.Conn) {
  _, channel, _ := RoutingInfo(ws.Request().URL.Path)
  pubSub := redis.PubSubConn{Conn: RedisPool.Get()}
  pubSub.PSubscribe(channel + ":*")

  for {
    switch m := pubSub.Receive().(type) {
    case redis.PMessage:
      if m.Channel != channel + ":" + identifier {
        websocket.Message.Send(ws, string(m.Data))
      }
    }
  }
}

func ReceiveMessages(identifier string, ws *websocket.Conn) {
  _, channel, _ := RoutingInfo(ws.Request().URL.Path)

  for {
    var message string
    websocket.Message.Receive(ws, &message)

    c := RedisPool.Get()
    c.Do("PUBLISH", channel + ":" + identifier, message)
    c.Close()
  }
}

func RoutingInfo(path string) (org, channel string, err error) {
  sections := strings.Split(path, "/")
  if len(sections) < 3 {
    err = errors.New("Connection needs to request an org/channel pair in path")
    return
  }

  routingInfo := []string{
    sections[1],
    sections[2],
  }

  org         = sections[1]
  channel     = strings.Join(routingInfo, ":")

  return
}
