package http_server_test

import (
	"fmt"
	"net"
	"testing"
	"time"

	go_http "github.com/Huy-DNA/go-http"
	"github.com/stretchr/testify/assert"
)

func TestBuildConfiguration(t *testing.T) {
  assert := assert.New(t)

  config := (go_http.HttpConfiguration{
    Port: 80,
  }).Build()

  assert.EqualExportedValues(config, go_http.HttpConfiguration{
    Ip: net.ParseIP("127.0.0.1"),
    Port: 80,
    Backlog: 100000,
  }, "Built config should be equal to explicitly init config")

  config = (go_http.HttpConfiguration{}).Build()

  assert.EqualExportedValues(config, go_http.HttpConfiguration{
    Ip: net.ParseIP("127.0.0.1"),
    Port: 8000,
    Backlog: 100000,
  }, "Built config should be equal to explicitly init config")
}

func TestHttpServerConn(t *testing.T) {
  assert := assert.New(t)

  config := (go_http.HttpConfiguration{Verbose: true}).Build()
  server := go_http.HttpServer{HttpConfiguration: config}

  error := server.Listen()

  assert.Nil(error, "Error should be nil");

  go server.Accept()

  time.Sleep(1000)
  simulateTcpConnect(8000)
  time.Sleep(1000)
}

func simulateTcpConnect(port uint16) {
  servAddr := fmt.Sprintf("localhost:%v", port)
  tcpAddr, _ := net.ResolveTCPAddr("tcp", servAddr)
  net.DialTCP("tcp", nil, tcpAddr)
}
