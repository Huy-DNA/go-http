package http_server

import (
	"reflect"
	"strconv"
)

type HttpConfiguration struct {
  Port int16
  Backlog uint64      `default:"100000"`
}

func (config *HttpConfiguration) Build() HttpConfiguration {
  typ := reflect.TypeOf(config)

  port := config.Port
  backlog := config.Backlog

  if backlog == 0 {
    f, _ := typ.FieldByName("Backlog")
    backlog, _ = strconv.ParseUint(f.Tag.Get("default"), 10, 64)
  }

  return HttpConfiguration{
    Port: port,
    Backlog: backlog,
  }
}
