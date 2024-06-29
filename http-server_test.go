package http_server_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
  go_http "github.com/Huy-DNA/go-http"
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
  config := (go_http.HttpConfiguration{}).Build()
  server := go_http.HttpServer{HttpConfiguration: config}

  server.Listen()
}
