package http_server_test

import (
	"fmt"
	"net"
	"testing"
	"time"
	"github.com/Huy-DNA/go-http/http-server"
	"github.com/stretchr/testify/assert"
)

func TestBuildConfiguration(t *testing.T) {
  assert := assert.New(t)

  config := (http_server.HttpConfiguration{
    Port: 80,
  }).Build()

  assert.EqualExportedValues(config, http_server.HttpConfiguration{
    Ip: net.ParseIP("127.0.0.1"),
    Port: 80,
    Backlog: 100000,
  }, "Built config should be equal to explicitly init config")

  config = (http_server.HttpConfiguration{}).Build()

  assert.EqualExportedValues(config, http_server.HttpConfiguration{
    Ip: net.ParseIP("127.0.0.1"),
    Port: 8000,
    Backlog: 100000,
  }, "Built config should be equal to explicitly init config")
}

func TestHttpServerConn(t *testing.T) {
  assert := assert.New(t)

  config := (http_server.HttpConfiguration{Verbose: true}).Build()
  server := http_server.BuildServer(config)

  error := server.Start()

  assert.Nil(error, "Error should be nil");

  go server.Accept()

  time.Sleep(1000)
  simulateTcpConnect(8000)
  time.Sleep(1000)
}

func TestHttpServerLoop(t *testing.T) {
  assert := assert.New(t)

  config := (http_server.HttpConfiguration{Verbose: true}).Build()
  server := http_server.BuildServer(config)

  error := server.Start()

  assert.Nil(error, "Error should be nil");

  server.Loop()

  time.Sleep(1000)
  // simulateTcpConnect(8000)
  time.Sleep(1000)
}

func simulateTcpConnect(port uint16) {
  servAddr := fmt.Sprintf("localhost:%v", port)
  tcpAddr, _ := net.ResolveTCPAddr("tcp", servAddr)
  tcpConn, error := net.DialTCP("tcp", nil, tcpAddr)
  if error != nil {
    panic("Failed to dial TCP!")
  }
  tcpConn.Write([]byte("Hello world!"));
}
