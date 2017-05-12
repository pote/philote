package main

import (
  "log"

  "github.com/gorilla/websocket"
  "github.com/satori/go.uuid"
)

type Philote struct {
  ID         string
  AccessKey  *AccessKey
  Hive       *hive
  ws         *websocket.Conn
  done       chan bool
}

func NewPhilote(ak *AccessKey, ws *websocket.Conn) (*Philote) {
  return &Philote{
    ws:    ws,
    done: make(chan bool),
    ID: uuid.NewV4().String(),
    AccessKey: ak,
  }
}

func (p *Philote) ListenToSocket() {
  for {
    message := &Message{}
    err := p.ws.ReadJSON(&message); if err != nil {
      p.logMsg("Invalid client message data: %s", err.Error() )
      if err.Error() == "EOF" {
        p.disconnect()
        break
      } else {
        continue
      }
    }

    p.logMsg("Received message from socket in channel " + message.Channel)

    if p.AccessKey.CanWrite(message.Channel) {
      go p.publish(message)
    } else {
      p.logMsg("Client does not have write permission for channel " + message.Channel + ", message dropped")
    }
  }
}

func (p *Philote) disconnect() {
  p.publish(&Message{Event: "close"})
  p.logMsg("Disconnecting from client")
  close(p.done)
}

func (p *Philote) Wait() {
  <-p.done
  p.logMsg("Disconnected")
}

func (p *Philote) publish(message *Message) {
  message.IssuerID = p.ID

  for _, philote := range p.Hive.Philotes {
    if p.ID == philote.ID {
      continue
    }

    for _, channel := range philote.AccessKey.Read {
      if message.Channel == channel {
        philote.ws.WriteJSON(message)
        break
      }
    }
  }
}

func (p *Philote) logMsg(message string, args ...interface{}) {
  log.Printf("[" + p.ID + "] " + message + "\n", args...)
}
