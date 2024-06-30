package server

import (
  "syscall"
)

type Connection struct {
  nfd uint16
  cliAddr syscall.Sockaddr
  server *Server
}

func (conn *Connection) OnMessage(callback func([] byte)) {
  buffer := make(chan []byte)
  conn.server.mesSubs.Store(conn.nfd, messageSubscriber{conn: conn, buffer: buffer})

  go func() {
    for data := range buffer {
      callback(data)
    }
  }()
}
