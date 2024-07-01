package http_server

import (
	"github.com/Huy-DNA/go-http/server-core"
)

type HttpConnection struct {
  conn <-chan *server.Connection 
  server *HttpServer
}
