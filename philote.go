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
}

func NewPhilote(ak *AccessKey, ws *websocket.Conn) (*Philote) {
  p := &Philote{
    ws:    ws,
    ID: uuid.NewV4().String(),
    AccessKey: ak,
  }

  return p
}

func (p *Philote) Listen() {
  for {
    message := &Message{}
    err := p.ws.ReadJSON(&message); if err != nil {
      p.logMsg("Invalid client message data: %s", err.Error() )
      if err.Error() == "EOF" {
        p.Hive.Disconnect <- p
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
  p.logMsg("Closing Philote")
  p.ws.Close()
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
