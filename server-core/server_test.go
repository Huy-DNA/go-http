package server_test

import (
	"fmt"
	"github.com/Huy-DNA/go-http/server-core"
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
	"time"
)

func TestBuildConfiguration(t *testing.T) {
	assert := assert.New(t)

	config := (server.Configuration{
		Port: 80,
	}).Build()

	assert.EqualExportedValues(config, server.Configuration{
		Ip:      net.ParseIP("127.0.0.1"),
		Port:    80,
		Backlog: 100000,
	}, "Built config should be equal to explicitly init config")

	config = (server.Configuration{}).Build()

	assert.EqualExportedValues(config, server.Configuration{
		Ip:      net.ParseIP("127.0.0.1"),
		Port:    8000,
		Backlog: 100000,
	}, "Built config should be equal to explicitly init config")
}

// FIXME: flaky test
func TestServerConnEarlyEventHandler(t *testing.T) {
	assert := assert.New(t)

	config := (server.Configuration{Verbose: true}).Build()
	server := server.New(config)

	connChan, error := server.Start()
	defer server.Stop()

	assert.Nil(error, "Error should be nil")

	time.Sleep(1000)
	tcpConn := simulateTcpConnect(8000)
	conn := <-connChan
	conn.OnMessage(func(data []byte) {
		fmt.Println(string(data))
	})
	tcpConn.Write([]byte("Hello World"))
	tcpConn.Close()
	time.Sleep(1000000000)
}

// FIXME: flaky test
func TestServerConnLateEventHandler(t *testing.T) {
	assert := assert.New(t)

	config := (server.Configuration{Verbose: true, Port: 8080}).Build()
	server := server.New(config)

	connChan, error := server.Start()
	defer server.Stop()

	assert.Nil(error, "Error should be nil")

	time.Sleep(1000)
	tcpConn := simulateTcpConnect(8080)
	conn := <-connChan
	tcpConn.Write([]byte("Hello World"))
	time.Sleep(1000)
	conn.OnMessage(func(data []byte) {
		fmt.Println(string(data))
	})
	tcpConn.Close()
	time.Sleep(1000000000)
}

func simulateTcpConnect(port uint16) *net.TCPConn {
	servAddr := fmt.Sprintf("localhost:%v", port)
	tcpAddr, _ := net.ResolveTCPAddr("tcp", servAddr)
	tcpConn, error := net.DialTCP("tcp", nil, tcpAddr)
	if error != nil {
		panic("Failed to dial TCP!")
	}

	return tcpConn
}
