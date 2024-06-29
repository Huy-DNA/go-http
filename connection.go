package http_server

import "syscall"

type Connection struct {
  nfd uint16
  cliAddr syscall.Sockaddr
  server *HttpServer
}


