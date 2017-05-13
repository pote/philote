package main

import(
  "github.com/ianschenck/envflag"
)

type config struct {
  jwtSecret []byte
  port      string
  version   string
}

func LoadConfig() (*config) {
  c := &config{}

  jwtSecret := envflag.String("JWT_SECRET", "", "JWT secret used to validate access keys.")
  port := envflag.String("PORT", "6380", "Port in which to serve Philote websocket connections")
  envflag.Parse()

  c.jwtSecret = []byte(*jwtSecret)
  c.port = *port

  return c
}
