package main

import (
  log "github.com/sirupsen/logrus"

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
  log.WithFields(log.Fields{"philote": p.ID}).Debug("Listening to Philote")
  for {
    message := &Message{}
    err := p.ws.ReadJSON(&message); if err != nil {
      log.WithFields(log.Fields{
        "philote": p.ID,
        "channel": message.Channel,
        "error": err.Error()}).Warn("Invalid client message data")
      if err.Error() == "EOF" {
        p.Hive.Disconnect <- p
        break
      } else {
        continue
      }
    }

    log.WithFields(log.Fields{"philote": p.ID, "channel": message.Channel}).Debug("Received message from socet")

    if p.AccessKey.CanWrite(message.Channel) {
      go p.publish(message)
    } else {
      log.WithFields(log.Fields{
        "philote": p.ID,
        "channel": message.Channel}).Info("Message dropped due to insufficient write permissions")
    }
  }
}

func (p *Philote) disconnect() {
  p.publish(&Message{Event: "close"})
  log.WithFields(log.Fields{"philote": p.ID}).Debug("Closing Philote")
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
