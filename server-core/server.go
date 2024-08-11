package server

import (
	"io"
	"log"
	"os"
	"sync"
	"syscall"
)

type messageSubscriber struct {
	conn   *Connection
	buffer chan []byte
}

type eventsController struct {
	mesSubs  sync.Map
	connChan chan *Connection
	quitChan chan struct{}
}

type Server struct {
	config Configuration
	eventsController
	sockFd  uint16
	epollFd uint16
}

func New(config Configuration) Server {
	return Server{
		config:  config,
		sockFd:  0,
		epollFd: 0,
		eventsController: eventsController{
			mesSubs:  sync.Map{},
			connChan: make(chan *Connection, config.Backlog),
			quitChan: make(chan struct{}),
		},
	}
}

func (server *Server) Start() (connChan <-chan *Connection, err error) {
	if server.sockFd != 0 {
		syscall.Close(int(server.sockFd))
		server.sockFd = 0
	}

	logger := server.getLogger()

	// socket creation
	sockFd, error := syscall.Socket(syscall.AF_INET6, syscall.SOCK_STREAM|syscall.SOCK_CLOEXEC, 0)

	if error != nil {
		logger.Printf("Socket creation failed")
		syscall.Close(sockFd)
		return nil, error
	}

	// socket binding
	sockaddr := &syscall.SockaddrInet6{
		Port:   int(server.config.Port),
		Addr:   [16]byte(server.config.Ip.To16()),
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

	server.loop()

	return server.connChan, nil
}

func (server *Server) Stop() {
	if server.sockFd != 0 {
		syscall.Close(int(server.sockFd))
	}
	server.sockFd = 0
	if server.epollFd != 0 {
		syscall.Close(int(server.epollFd))
	}
	server.epollFd = 0
	close(server.quitChan)
}

// Internal methods

func (server *Server) loop() <-chan *Connection {
	go server.loopAccept()
	go server.loopMessage()

	return server.connChan
}

func (server *Server) loopAccept() {
	var loggerDest io.Writer = os.Stdout
	if !server.config.Verbose {
		loggerDest = io.Discard
	}
	logger := log.New(loggerDest, "Log: ", log.LstdFlags)

	for {
		select {
		case <-server.quitChan:
			logger.Printf("Quitting accept loop")
			return // will be executed if the server is closed
		default:
			conn, _ := server.accept()
			server.connChan <- conn
		}
	}
}

func (server *Server) loopMessage() {
	logger := server.getLogger()
	for {
		select {
		case <-server.quitChan:
			logger.Printf("Quitting message loop")
			return // will be executed if the server is closed
		default:
			epollEvent := make([]syscall.EpollEvent, 100)
			n, error := syscall.EpollWait(int(server.epollFd), epollEvent, 1000)
			if error != nil {
				logger.Printf("error while epolling for messages: %v", error.Error())
				continue
			}

			for i := range n {
				event := epollEvent[i]
				switch {
				case event.Events&syscall.EPOLLERR > 0:
					logger.Printf("EPOLLERR on socket with fd %v", event.Fd)
					syscall.Close(int(event.Fd))
					server.mesSubs.Delete(uint16(event.Fd))
				case event.Events&syscall.EPOLLRDHUP > 0:
					logger.Printf("EPOLLRDHUP on socket with fd %v", event.Fd)
					syscall.Close(int(event.Fd))
					server.mesSubs.Delete(uint16(event.Fd))
				case event.Events&syscall.EPOLLIN > 0:
					logger.Printf("EPOLLIN on socket with fd %v", event.Fd)
					sub, ok := server.mesSubs.Load(uint16(event.Fd))
					if !ok {
						continue
					}
					data := make([]byte, 0)
					for {
						const BUFFER_SIZE uint = 1000
						pdata := make([]byte, BUFFER_SIZE)
						n, err := syscall.Read(int(event.Fd), pdata)
						if err != nil {
							break
						}
						if n > 0 {
							data = append(data, pdata[:n]...)
						}
						if n < int(BUFFER_SIZE) {
							break
						}
					}
					sub.(messageSubscriber).buffer <- data
				}
			}
		}
	}
}

func (server *Server) addConnToServerEpoll(nfd uint16) (err error) {
	epollEvents := syscall.EpollEvent{
		Events: syscall.EPOLLIN | syscall.EPOLLRDHUP | (syscall.EPOLLET & 0xffffffff), // edge-triggered mode, wait for reading readiness/reading disconnection from the remote peer
		Fd:     int32(nfd),
	}
	error := syscall.EpollCtl(int(server.epollFd), syscall.EPOLL_CTL_ADD, int(nfd), &epollEvents)

	if error != nil {
		return error
	}
	return nil
}

func (server *Server) initEpoll() (err error) {
	epollFd, error := syscall.EpollCreate1(syscall.EPOLL_CLOEXEC)

	server.epollFd = uint16(epollFd)

	if error != nil {
		return error
	}

	return nil
}

func (server *Server) accept() (conn *Connection, err error) {
	nfd, cliAddr, error := syscall.Accept(int(server.sockFd))

	if error != nil {
		return nil, error
	}

	logger := server.getLogger()

	logger.Printf("New connection accepted")

	error = server.addConnToServerEpoll(uint16(nfd))

	if error != nil {
		logger.Printf("Failed to add connection to epoll")
		return nil, error
	}

	return &Connection{nfd: uint16(nfd), cliAddr: cliAddr, server: server}, nil
}

func (server *Server) getLogger() *log.Logger {
	var loggerDest io.Writer = os.Stdout
	if !server.config.Verbose {
		loggerDest = io.Discard
	}
	return log.New(loggerDest, "Log: ", log.LstdFlags)
}
