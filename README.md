# go-http
This is a challenge which aims at implementing an HTTP server from scratch using Go.

## Constraints

* Libraries that wrap around `libc` are permitted.
* No libraries that abstract away the internals of `libc` are allowed.
* Utilities library for compression, error handling, etc. are allowed.

## Requirements

- [ ] Should be HTTP/1.1-compliant (P2)
- [X] Able to handle multiple concurrent TCP connection. (P1)
- [ ] Should perform request assembly. (P1)
- [ ] Should be able to parse HTTP requests. (P1)
- [ ] Allow defining handlers for requests. (P1)
- [ ] Allow defining plugins to perform logging/monitoring. (P2)
- [ ] Can act as a static HTTP server as the fallback behavior. (P2)
- [ ] Allow server configuration. (P2)
- [ ] Perform response compression. (P3)

## References

* HTTP/1.1 spec: https://www.rfc-editor.org/rfc/rfc7230#section-3.3.3
