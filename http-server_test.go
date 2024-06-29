package http_server_test

import (
	"testing"
	go_http "github.com/Huy-DNA/go-http"
	"github.com/stretchr/testify/assert"
)

func TestBuildConfiguration(t *testing.T) {
  assert := assert.New(t)

  config := (go_http.HttpConfiguration{
    Port: 80,
  }).Build()

  assert.EqualExportedValues(config, go_http.HttpConfiguration{
    Port: 80,
    Backlog: 100000,
  }, "Built config should be equal to explicitly init config")

  config = (go_http.HttpConfiguration{}).Build()

  assert.EqualExportedValues(config, go_http.HttpConfiguration{
    Port: 8000,
    Backlog: 100000,
  }, "Built config should be equal to explicitly init config")
}

func TestHttpServerListent(t *testing.T) {
  assert := assert.New(t)

  config := (go_http.HttpConfiguration{Verbose: true}).Build()
  server := go_http.HttpServer{HttpConfiguration: config}

  error := server.Listen()

  assert.Nil(error, "Error should be nil");
}
