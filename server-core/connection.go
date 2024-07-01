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
    for {
      select {
      case data := <-buffer:
        callback(data)
      case <-conn.server.quitChan:
        return
      }
    }
  }()
}

func (conn *Connection) RemoteIp() string {
  return string(conn.cliAddr.(*syscall.SockaddrInet6).Addr[:])
}

func (conn *Connection) RemotePort() uint16 {
  return uint16(conn.cliAddr.(*syscall.SockaddrInet6).Port)
}
