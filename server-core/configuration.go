package server

import (
  "reflect"
  "strconv"
  "net"
)

type Configuration struct {
  Ip net.IP           `default:"127.0.0.1"`
  Port uint16         `default:"8000"`
  Backlog uint64      `default:"100000"`
  Verbose bool
}

func (config Configuration) Build() Configuration {
  typ := reflect.TypeOf(config)

  ip := config.Ip;
  if ip == nil {
    f, _ := typ.FieldByName("Ip") 
    ip = net.ParseIP(f.Tag.Get("default"))
  }

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

  return Configuration{
    Ip: ip,
    Port: port,
    Backlog: backlog,
    Verbose: config.Verbose,
  }
}
