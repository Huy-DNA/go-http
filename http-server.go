package http_server

import (
	"io"
	"log"
	"os"
	"syscall"
)

type HttpServer struct {
  HttpConfiguration
  sockFd uint16
  epollFd uint16
}

func (server *HttpServer) Listen() (err error) {
  if server.sockFd != 0 {
    syscall.Close(int(server.sockFd))
    server.sockFd = 0
  }

  config := server.HttpConfiguration

  var loggerDest io.Writer = os.Stdout
  if !config.Verbose {
    loggerDest = io.Discard
  }
  logger := log.New(loggerDest, "Log: ", log.LstdFlags)

  // socket creation
  sockFd, error := syscall.Socket(syscall.AF_INET6, syscall.SOCK_STREAM | syscall.SOCK_CLOEXEC, 0)

  if error != nil {
    logger.Printf("Socket creation failed")
    syscall.Close(sockFd)
    return error
  }
 
  // socket binding
  sockaddr := &syscall.SockaddrInet6{
		Port: int(config.Port),
    Addr: [16]byte(config.Ip.To16()),
	  ZoneId: 0,
	}

  error = syscall.Bind(sockFd, sockaddr)

  if error != nil {
    logger.Printf("Socket binding to %v:%v failed", config.Ip, config.Port)
    syscall.Close(sockFd)
    return error
  }

  // socket listening
  error = syscall.Listen(sockFd, int(config.Backlog))
  
  if error != nil {
    logger.Printf("Server failed to listen on port %v", config.Port)
    syscall.Close(sockFd)
    return error
  } else {
    logger.Printf("Server is listening on port %v", config.Port)
  }

  error = server.initEpoll()
  if error != nil {
    logger.Printf("Server failed to init epoll")
    return error
  }

  server.sockFd = uint16(sockFd) 

  return nil
}

func (server *HttpServer) initEpoll() (err error) {
  if server.epollFd != 0 {
    syscall.Close(int(server.epollFd)) 
    server.epollFd = 0
  }

  epollFd, error := syscall.EpollCreate1(syscall.EPOLL_CLOEXEC)
 
  server.epollFd = uint16(epollFd)

  if error != nil {
    return error
  }
  
  return nil
}

func (server *HttpServer) Accept() (conn *Connection, err error) {
  config := server.HttpConfiguration
  nfd, cliAddr, err := syscall.Accept(int(server.sockFd)) 

  srvAddr := &syscall.SockaddrInet6{
		Port: int(config.Port),
    Addr: [16]byte(config.Ip.To16()),
	  ZoneId: 0,
	}

  if err != nil {
    return nil, err
  }

  return &Connection{nfd: uint16(nfd), cliAddr: cliAddr, srvAddr: srvAddr}, nil
}
