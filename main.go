package main

import (
  "net/http"
  "runtime"

  log "github.com/sirupsen/logrus"
  "github.com/gorilla/websocket"
)

var Config *config = LoadConfig()
var Upgrader = websocket.Upgrader{
  ReadBufferSize:  Config.readBufferSize,
  WriteBufferSize: Config.writeBufferSize,
  CheckOrigin: func(r *http.Request) bool {
    return true
  },
}

func main() {
  log.WithFields(log.Fields{
    "version": VERSION,
    "port": Config.port,
    "cores": runtime.NumCPU()}).Info("Initializing Philotic Network")

  log.WithFields(log.Fields{
    "read-buffer-size": Config.readBufferSize,
    "write-buffer-size": Config.writeBufferSize,
    "max-connections": Config.maxConnections}).Debug("Configuration options:")

  h := NewHive()
  http.HandleFunc("/", h.ServeNewConnection)
  http.HandleFunc("/api", h.ServeAPICall)

  err := http.ListenAndServe(":" + Config.port, nil); if err != nil {
    log.Fatal("ListenAndServe: ", err)
  }
}
