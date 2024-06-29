package http_server

import (
  "syscall"
)

type HttpServer struct {
  HttpConfiguration
}

func (server *HttpServer) Listen() (err error) {
  config := server.HttpConfiguration
  sockFd, error := syscall.Socket(syscall.AF_INET6, syscall.SOCK_STREAM, 0)
  
  if error != nil {
    return error
  }
  
  syscall.Listen(sockFd, int(config.Backlog))

  return nil
}

