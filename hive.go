package main

import(
  "net/http"
  "strings"

  log "github.com/sirupsen/logrus"
)

type hive struct {
  Philotes    map[string]*Philote
  Connect     chan *Philote
  Disconnect  chan *Philote
}

func NewHive() (*hive) {
  h := &hive{
    Philotes:   map[string]*Philote{},
    Connect:    make(chan *Philote),
    Disconnect: make(chan *Philote),
  }

  go h.MaintainPhiloteIndex()

  return h
}

func (h *hive) MaintainPhiloteIndex() {
  log.Debug("Starting bookeeper")

  for {
    select {
    case p := <- h.Connect:
      if len(h.Philotes) >= Config.maxConnections {
        log.WithFields(log.Fields{"philote": p.ID}).Warn("MAX_CONNECTIONS limit reached, dropping new connection")
        p.disconnect()
      }

      log.WithFields(log.Fields{"philote": p.ID}).Debug("Registering Philote")
      p.Hive = h
      h.Philotes[p.ID] = p
      go p.Listen()
    case p := <- h.Disconnect:
      log.WithFields(log.Fields{"philote": p.ID}).Debug("Disconnecting Philote")
      delete(h.Philotes, p.ID)
      p.disconnect()
    }
  }
}

func (h *hive) ServeNewConnection(w http.ResponseWriter, r *http.Request) {
  auth := strings.TrimSpace(strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer"))
  if auth == "" {
    r.ParseForm()
    auth = r.Form.Get("auth")
    log.WithFields(log.Fields{"auth": auth}).Debug("Empty Authorization header, trying querystring #auth param")
  }

  accessKey, err := NewAccessKey(auth); if err != nil {
    log.WithFields(log.Fields{"error": err.Error(), "auth": auth }).Warn("Can't create Access key")
    w.Write([]byte(err.Error()))
    return
  }

  connection, err := Config.Upgrader.Upgrade(w, r, nil); if err != nil {
    log.WithFields(log.Fields{"error": err.Error()}).Warn("Can't upgrade connection")
    w.Write([]byte(err.Error()))
    return
  }

  philote := NewPhilote(accessKey, connection)
  h.Connect <- philote
}
