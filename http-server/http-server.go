package http_server

import (
	"io"
	"log"
	"os"
	"syscall"
  "sync"
)

type messageSubscriber struct {
  conn *Connection
  callback func([]byte)
}

type eventsController struct {
  mesSubs sync.Map
  connChan chan *Connection
  stopped bool
}

type HttpServer struct {
  config HttpConfiguration
  eventsController
  sockFd uint16
  epollFd uint16
}

func BuildServer(config HttpConfiguration) HttpServer {
  return HttpServer{
    config: config,
    sockFd: 0,
    epollFd: 0,
  }
}

func (server *HttpServer) Start() (connChan <-chan *Connection, err error) {
  if server.sockFd != 0 {
    syscall.Close(int(server.sockFd))
    server.sockFd = 0
  }

  var loggerDest io.Writer = os.Stdout
  if !server.config.Verbose {
    loggerDest = io.Discard
  }
  logger := log.New(loggerDest, "Log: ", log.LstdFlags)

  // socket creation
  sockFd, error := syscall.Socket(syscall.AF_INET6, syscall.SOCK_STREAM | syscall.SOCK_CLOEXEC, 0)

  if error != nil {
    logger.Printf("Socket creation failed")
    syscall.Close(sockFd)
    return nil, error
  }
 
  // socket binding
  sockaddr := &syscall.SockaddrInet6{
		Port: int(server.config.Port),
    Addr: [16]byte(server.config.Ip.To16()),
	  ZoneId: 0,
	}

  error = syscall.Bind(sockFd, sockaddr)

  if error != nil {
    logger.Printf("Socket binding to %v:%v failed", server.config.Ip, server.config.Port)
    syscall.Close(sockFd)
    return nil, error
  }

  // socket listening
  error = syscall.Listen(sockFd, int(server.config.Backlog))
  
  if error != nil {
    logger.Printf("Server failed to listen on port %v", server.config.Port)
    syscall.Close(sockFd)
    return nil, error
  } else {
    logger.Printf("Server is listening on port %v", server.config.Port)
  }

  error = server.initEpoll()
  if error != nil {
    logger.Printf("Server failed to init epoll")
    return nil, error
  }

  server.sockFd = uint16(sockFd) 
  server.stopped = false

  server.loop()

  return server.connChan, nil
}

func (server *HttpServer) Stop() {
  server.stopped = true
  if server.sockFd != 0 {
    syscall.Close(int(server.sockFd))
  }
  server.sockFd = 0
}

// Internal methods

func (server *HttpServer) loop() <-chan *Connection {
  go server.loopAccept()
  go server.loopMessage()

  return server.connChan
}

func (server *HttpServer) loopAccept() {
  for !server.stopped {
    conn, _ := server.accept()
    server.connChan <- conn
  }
}

func (server *HttpServer) loopMessage() {
  var loggerDest io.Writer = os.Stdout
  if !server.config.Verbose {
    loggerDest = io.Discard
  }
  logger := log.New(loggerDest, "Log: ", log.LstdFlags)

  for !server.stopped {
    epollEvent := make([] syscall.EpollEvent, 100)
    n, error := syscall.EpollWait(int(server.epollFd), epollEvent, -1)
    if error != nil {
      logger.Printf("error while epolling for messages: %v", error.Error())
      continue
    }

    for i := range n {
      event := epollEvent[i]
      switch {
      case event.Events & syscall.EPOLLERR > 0:
        logger.Printf("EPOLLERR on socket with fd %v", event.Fd)
        syscall.Close(int(event.Fd))
        server.mesSubs.Delete(uint16(event.Fd))
      case event.Events & syscall.EPOLLRDHUP > 0:
        logger.Printf("EPOLLRDHUP on socket with fd %v", event.Fd)
        syscall.Close(int(event.Fd))
        server.mesSubs.Delete(uint16(event.Fd))
      case event.Events & syscall.EPOLLIN > 0:
        logger.Printf("EPOLLIN on socket with fd %v", event.Fd)
        sub, ok := server.mesSubs.Load(uint16(event.Fd))
        if !ok {
          continue
        }
        data := make([]byte, 0)
        for {
          const BUFFER_SIZE uint = 1000
          pdata := make([]byte, BUFFER_SIZE)
          n, err := syscall.Read(int(event.Fd), data)
          if err != nil {
            break
          }
          data = append(data, pdata...)
          if n < int(BUFFER_SIZE) {
            break
          }
        }
        sub.(messageSubscriber).callback(data)
      }
    }
  }
}

func (server *HttpServer) addConnToServerEpoll(nfd uint16) (err error) {
  epollEvents := syscall.EpollEvent {
    Events: syscall.EPOLLIN | syscall.EPOLLRDHUP | (syscall.EPOLLET & 0xffffffff), // level-triggered mode, wait for reading readiness/reading disconnection from the remote peer
    Fd: int32(server.sockFd),
  }
  error := syscall.EpollCtl(int(server.epollFd), syscall.EPOLL_CTL_ADD, int(nfd), &epollEvents)

  if error != nil {
    return error
  }
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

func (server *HttpServer) accept() (conn *Connection, err error) {
  nfd, cliAddr, error := syscall.Accept(int(server.sockFd)) 

  if error != nil {
    return nil, error
  }

  var loggerDest io.Writer = os.Stdout
  if !server.config.Verbose {
    loggerDest = io.Discard
  }
  logger := log.New(loggerDest, "Log: ", log.LstdFlags)
  
  logger.Printf("New connection accepted")
  
  error = server.addConnToServerEpoll(uint16(nfd))

  if error != nil {
    logger.Printf("Failed to add connection to epoll")
    return nil, error
  }

  return &Connection{nfd: uint16(nfd), cliAddr: cliAddr, server: server}, nil
}
