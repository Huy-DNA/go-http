package http_server

type HttpConfiguration struct {
  Port int16
  Backlog uint64      `default:"100000"`
}
