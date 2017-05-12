package main

import (
  "log"
  "net/http"
  "runtime"
  "strings"

  "github.com/ianschenck/envflag"
  "github.com/gorilla/websocket"
)

var Hive *hive = NewHive()

var Upgrader = websocket.Upgrader{
  ReadBufferSize:  1024,
  WriteBufferSize: 1024,
}

func main() {
  //jwtToken := envflag.String("JWT_SECRET", "", "Secret JWT token that validates access keys.")
  port := envflag.String("PORT", "6380", "Port in which to serve Philote websocket connections")

  envflag.Parse()

  log.Printf("[Main] Initializing Philotic Network\n")
  log.Printf("[Main] Version: %v\n", VERSION)
  log.Printf("[Main] Port: %v\n", *port)
  log.Printf("[Main] Cores: %v\n", runtime.NumCPU())

  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    results := strings.Split(r.Header.Get("Authorization"), "Bearer "); if len(results) < 2 {
      err :=  errors.New("No access token found")
      log.Println(err)
      w.Write([]byte(err.Error()))
      return
    }

    accessKey, err := NewAccessKey(results[1]); if err != nil {
      log.Println(err)
      w.Write([]byte(err.Error()))
      return
    }

    connection, err := Upgrader.Upgrade(w, r, nil); if err != nil {
      log.Println(err)
      w.Write([]byte(err.Error()))
      return
    }

    philote := NewPhilote(accessKey, connection)
    Hive.NewPhilotes <- philote
    go philote.ListenToSocket()
    philote.Wait()
  })

  err := http.ListenAndServe(":" + *port, nil)
  if err != nil {
    log.Fatal("ListenAndServe: ", err)
  }
}
