package main

import (
  "net/http"
  "runtime"
  "strings"

  log "github.com/sirupsen/logrus"
  "github.com/gorilla/websocket"
)

var Hive *hive = NewHive()
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

  http.HandleFunc("/", ServeNewConnection)

  err := http.ListenAndServe(":" + Config.port, nil); if err != nil {
    log.Fatal("ListenAndServe: ", err)
  }
}

func ServeNewConnection(w http.ResponseWriter, r *http.Request) {
  auth := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
  accessKey, err := NewAccessKey(auth); if err != nil {
    log.WithFields(log.Fields{"error": err.Error()}).Warn("Can't create Access key")
    w.Write([]byte(err.Error()))
    return
  }

  connection, err := Upgrader.Upgrade(w, r, nil); if err != nil {
    log.WithFields(log.Fields{"error": err.Error()}).Warn("Can't upgrade connection")
    w.Write([]byte(err.Error()))
    return
  }

  philote := NewPhilote(accessKey, connection)
  Hive.Connect <- philote
}
