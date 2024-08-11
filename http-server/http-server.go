package http_server

import (
	"github.com/Huy-DNA/go-http/server-core"
)

type HttpServer struct {
	server server.Server
	conn   <-chan *server.Connection
	req    chan Request
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

func getReqChan(connChan <-chan *server.Connection) chan Request {
	reqChan := make(chan Request)
	go func() {
		for conn := range connChan {
			remoteAddr := RemoteAddress{Ip: conn.RemoteIp(), Port: conn.RemotePort()}
			req := Request{remoteAddr: remoteAddr}
			buffer := make([]byte, 0)
			conn.OnMessage(func(bytes []byte) {
				buffer := append(buffer, bytes...)
				done := partialFill(&req, &buffer)
				if done {
					reqChan <- req
					req = Request{remoteAddr: remoteAddr}
				}
			})
		}
	}()

	return reqChan
}

func partialFill(req *Request, buffer *[]byte) (done bool) {
	return false
}
