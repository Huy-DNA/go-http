# go-http
This is a challenge which aims at implementing an HTTP server from scratch using Go.

## Constraints

* Libraries that wrap around `libc` are permitted.
* No libraries that abstract away the internals of `libc` are allowed.
* Utilities library for compression, error handling, etc. are allowed.

## Requirements

* Able to handle multiple concurrent TCP connection. (P1)
* Should perform request assembly. (P1)
* Should be able to parse HTTP requests. (P1)
* Allow defining plugins to plug into the request handling process. (P1)
* Can act as a static HTTP server as the fallback behavior. (P2)
* Perform response compression. (P3)
