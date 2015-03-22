package routes

import(
  "io"
  "net/http"
)

type Root struct { }

func (r *Root) Match(req *http.Request) bool {
  return true
}

func (r *Root) Perform(w http.ResponseWriter, req *http.Request) {
  io.WriteString(w, "Welcome to the Philotic Network.\n")
}

