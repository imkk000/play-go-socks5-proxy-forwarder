# socks5-server

A lightweight SOCKS5 proxy server with optional upstream proxy chaining support.

## Features

- SOCKS5 proxy server listening on port `9001`
- Optional upstream proxy chaining (chain multiple SOCKS5 proxies)
- Optional request logging with timestamp, source, and destination
- Graceful shutdown on `SIGINT` / `SIGTERM`

## Requirements

- Go 1.26+

## Build

```sh
go build -o socks5-server .
```

Or with Docker:

```sh
docker build -t socks5-server .
```

## Usage

### Run directly

```sh
./socks5-server
```

### With proxy chaining

Pass a comma-separated list of upstream SOCKS5 proxies via `-chains`. Traffic is forwarded through them in order (left to right).

```sh
./socks5-server -chains 10.0.0.1:1080,10.0.0.2:1080
```

### With logging enabled

```sh
./socks5-server -log
./socks5-server -chains 10.0.0.1:1080,10.0.0.2:1080 -log
```

### Run with Docker

```sh
docker run -p 9001:9001 socks5-server
```

With chains:

```sh
docker run -p 9001:9001 socks5-server -chains 10.0.0.1:1080,10.0.0.2:1080
```

## Flags

| Flag      | Default  | Description                                               |
|-----------|----------|-----------------------------------------------------------|
| `-chains` | _(none)_ | Comma-separated list of upstream SOCKS5 proxies to chain  |
| `-log`    | `false`  | Enable request logging to stdout                          |

## Server address

The server listens on `0.0.0.0:9001` by default.

## Logging

Logging is disabled by default. Enable it with `-log` to print to stdout:

```
[2026-04-20T10:00:00.000000000Z] start on :9001
[2026-04-20T10:00:01.000000000Z] chain 1: 10.0.0.1:1080
[2026-04-20T10:00:01.000000000Z] chain 2: 10.0.0.2:1080
[2026-04-20T10:00:02.000000000Z] from: 192.168.1.5:54321 -> example.com:443 (1)
```

## Dependencies

- [things-go/go-socks5](https://github.com/things-go/go-socks5)
- [golang.org/x/net/proxy](https://pkg.go.dev/golang.org/x/net/proxy)
