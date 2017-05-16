package main

import(
  log "github.com/sirupsen/logrus"
  "github.com/ianschenck/envflag"
)

type config struct {
  jwtSecret       []byte
  port            string
  version         string
  readBufferSize  int
  writeBufferSize int
  maxConnections  int
  log             log.Level
}


func LoadConfig() (*config) {
  c := &config{}

  secret := envflag.String(
    "SECRET",
    "",
    "JWT secret used to validate access keys.")
  port := envflag.String(
    "PORT",
    "6380",
    "Port in which to serve Philote websocket connections")
  logLevel := envflag.String(
    "LOGLEVEL",
    "info",
    "Log level, accepts: 'debug', 'info', 'warning', 'error', 'fatal', 'panic'")

  maxConnections := envflag.Int(
    "MAX_CONNECTIONS",
    255,
    "Maximum amount of permitted concurrent connections")

  readBufferSize := envflag.Int(
    "READ_BUFFER_SIZE",
    1024,
    "Size (in bytes) for the read buffer")

  writeBufferSize := envflag.Int(
    "WRITE_BUFFER_SIZE",
    1024,
    "Size (in bytes) for the write buffer")

  envflag.Parse()

  c.jwtSecret = []byte(*secret)
  c.port = *port
  c.maxConnections = *maxConnections
  c.readBufferSize = *readBufferSize
  c.writeBufferSize = *writeBufferSize

  var err error
  c.log, err = log.ParseLevel(*logLevel); if err != nil {
    log.WithFields(log.Fields{"error": err}).Panic("Unparsable log level")
  }
  log.SetLevel(c.log)

  return c
}
