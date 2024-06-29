package http_server

type HttpConfiguration struct {
  Port int16
}

type HttpServer struct {
  HttpConfiguration
}

func StartServer(config HttpConfiguration) HttpServer {
  return HttpServer{
    HttpConfiguration: config,
  }
}
