package main

import(
  log "github.com/sirupsen/logrus"
  "github.com/ianschenck/envflag"
)

type config struct {
  jwtSecret []byte
  port      string
  version   string
  log       log.Level
}


func LoadConfig() (*config) {
  c := &config{}

  jwtSecret := envflag.String("JWT_SECRET", "", "JWT secret used to validate access keys.")
  port := envflag.String("PORT", "6380", "Port in which to serve Philote websocket connections")
  logLevel := envflag.String("LOG", "info", "Log level, accepts: 'debug', 'info', 'warning', 'error', 'fatal', 'panic'")
  envflag.Parse()

  c.jwtSecret = []byte(*jwtSecret)
  c.port = *port

  var err error
  c.log, err = log.ParseLevel(*logLevel); if err != nil {
    log.WithFields(log.Fields{"error": err}).Panic("Unparsable log level")
  }
  log.SetLevel(c.log)

  return c
}
