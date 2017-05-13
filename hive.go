package main

type hive struct {
  Philotes []*Philote
  NewPhilotes chan *Philote
}

func NewHive() (*hive) {
  h := &hive{Philotes: []*Philote{}}

  go h.RegisterNewPhilotes()

  return h
}

func (h *hive) RegisterNewPhilotes() {
  for {
    philote := <- h.NewPhilotes

    philote.Hive = h
    h.Philotes = append(h.Philotes, philote)
  }
}

