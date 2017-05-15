package main

import(
)

type hive struct {
  Philotes    map[string]*Philote
  Connect     chan *Philote
  Disconnect  chan *Philote
}

func NewHive() (*hive) {
  h := &hive{Philotes: map[string]*Philote{}}

  go h.MaintainPhiloteIndex()

  return h
}

func (h *hive) MaintainPhiloteIndex() {
  for {
    select {
    case p := <- h.Connect:
      p.Hive = h
      h.Philotes[p.ID] = p
      go p.Listen()
    case p := <- h.Disconnect:
      delete(h.Philotes, p.ID)
      p.disconnect()
    }
  }
}

