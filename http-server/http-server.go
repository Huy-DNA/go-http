package http_server

import (
	"github.com/Huy-DNA/go-http/server-core"
)

type HttpServer struct {
  server server.Server
}

type HttpServerConfiguration struct {
  server.Configuration
}

func New(config HttpServerConfiguration) HttpServer {
  return HttpServer{
    server: server.New(config.Configuration),
  }
}

func (server *HttpServer) Start() (connChan *HttpConnection, err error) {
  conn, error := server.server.Start()
  if error != nil {
    return nil, error
  }
  return &HttpConnection{conn: conn}, nil
}

func (server *HttpServer) Stop() {
  server.server.Stop()
}
