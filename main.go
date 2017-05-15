package main

import (
  "net/http"
  "runtime"

  log "github.com/sirupsen/logrus"
  "github.com/gorilla/websocket"
)

var Config *config = LoadConfig()
// TODO: These buffer sizes should be configurable.
var Upgrader = websocket.Upgrader{
  ReadBufferSize:  1024,
  WriteBufferSize: 1024,
}

func main() {
  log.WithFields(log.Fields{
    "version": VERSION,
    "port": Config.port,
    "cores": runtime.NumCPU()}).Info("Initializing Philotic Network")

  h := NewHive()
  http.HandleFunc("/", h.ServeNewConnection)

  err := http.ListenAndServe(":" + Config.port, nil); if err != nil {
    log.Fatal("ListenAndServe: ", err)
  }
}
