package http_server

import (
	"io"
	"log"
	"os"
	"syscall"
)

type HttpServer struct {
  HttpConfiguration
}

func (server *HttpServer) Listen() (err error) {
  config := server.HttpConfiguration

  var loggerDest io.Writer = os.Stdout
  if !config.Verbose {
    loggerDest = io.Discard
  }
  logger := log.New(loggerDest, "Log: ", log.LstdFlags)

  // socket creation
  sockFd, error := syscall.Socket(syscall.AF_INET6, syscall.SOCK_STREAM, 0)
  
  if error != nil {
    logger.Printf("Socket creation failed")
    return error
  }
  
  // socket listening
  error = syscall.Listen(sockFd, int(config.Backlog))
  
  if error != nil {
    logger.Printf("Server failed to listen on port %v", config.Port)
    return error
  } else {
    logger.Printf("Server is listening on port %v", config.Port)
  }

  return nil
}

