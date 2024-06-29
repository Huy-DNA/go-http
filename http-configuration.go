package http_server

import (
	"reflect"
	"strconv"
)

type HttpConfiguration struct {
  Port uint16         `default:"8000"`
  Backlog uint64      `default:"100000"`
  Verbose bool
}

func (config HttpConfiguration) Build() HttpConfiguration {
  typ := reflect.TypeOf(config)

  port := config.Port
  if port == 0 {
    f, _ := typ.FieldByName("Port")
    _port, _ := strconv.ParseUint(f.Tag.Get("default"), 10, 16)
    port = uint16(_port)
  }

  backlog := config.Backlog
  if backlog == 0 {
    f, _ := typ.FieldByName("Backlog")
    backlog, _ = strconv.ParseUint(f.Tag.Get("default"), 10, 64)
  }

  return HttpConfiguration{
    Port: port,
    Backlog: backlog,
    Verbose: config.Verbose,
  }
}
