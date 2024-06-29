package http_server

import (
	"fmt"
	"syscall"
)

type HttpServer struct {
  HttpConfiguration
}

func (server *HttpServer) Listen() (err error) {
  config := server.HttpConfiguration
  sockFd, error := syscall.Socket(syscall.AF_INET6, syscall.SOCK_STREAM, 0)
  
  if error != nil {
    if config.Verbose {
      fmt.Printf("Socket creation failed")
    }
    return error
  }
  
  error = syscall.Listen(sockFd, int(config.Backlog))
  
  if error != nil {
    if config.Verbose {
      fmt.Printf("Server failed to listen on port %v", config.Port)
    }

    return error
  } else {
    if config.Verbose {
      fmt.Printf("Server is listening on port %v", config.Port)
    }
  }

  return nil
}

