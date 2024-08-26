# gproxy

| English | [简体中文](./README.zh-CN.md) |

## Install

### With go

```bash
go install github.com/aak1247/gproxy
```

### Binary

See [Releases](https://github.com/aak1247/gproxy/releases/)

## Usage

### Arguments

gproxy [-m mode] [-p port] targetURL

- ``mode``: http or tcp
- ``port``: local port, example: 8080

### Proxy a http website:

``gproxy -m http -p 8080 http://github.com/``

then, a simple http server will run on local 8080 port to serve the content of github.com

#### HTTPS

``gproxy -m https --key keyPath --cert certPath -p 8080 http://github.com/``

### Proxy a tcp service:

``gproxy -m tcp -p 2333 github.com:443``

then, all tcp request to localhost:2333 will be sent to github.com:443 and response from github.com:443 will be sent
back to client

### Proxy a websocket service:

``gproxy -m ws -p 8081 ws://abc.com``

then, a websocket server will run on port 8081 and proxy all ws requests to abc.com. All path and query will be kept and sent to `abc.com`

#### WSS

``gproxy -m wss --key keyPath --cert certPath -p 8081 ws://abc.com``

## Build

### With makefile

``make release VERSION=$VERSION``: will generate all linux/windows/osx binary

### Run by command

#### Target linux

``CGO_ENABLED=0 GOOS=linux GOARCH=amd64  go build -o gproxy-linux-amd64``

#### Target windows

``CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o gproxy-windows-amd64.exe``

#### Target mac

``CGO_ENABLED=0 GOOS=darwin  GOARCH=amd64 go build -o gproxy-darwin-amd64``
