package http_server

import (
	"github.com/Huy-DNA/go-http/server-core"
)

type HttpServer struct {
  server server.Server
  conn <-chan *server.Connection
  req chan Request
}

type HttpServerConfiguration struct {
  server.Configuration
}

func New(config HttpServerConfiguration) HttpServer {
  return HttpServer{
    server: server.New(config.Configuration),
  }
}

func (server *HttpServer) Start() (reqChan <-chan Request, err error) {
  if server.conn != nil {
    return server.req, nil 
  }

  conn, error := server.server.Start()
  if error != nil {
    return nil, error
  }

  server.conn = conn
  server.req = getReqChan(conn) 
  
  return server.req, nil
}

func (server *HttpServer) Stop() {
  server.server.Stop()
  server.conn = nil
  close(server.req)
  server.req = nil
}

func getReqChan(conn <-chan *server.Connection) chan Request {
  return nil
}
