package http_server

func StartServer(config HttpConfiguration) HttpServer {
  return HttpServer{
    HttpConfiguration: config,
    sockFd: 0,
    epollFd: 0,
  }
}
