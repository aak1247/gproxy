# gproxy

| English | [简体中文](./README.zh-CN.md) |

## Usage

### Arguments
gproxy [-m mode] [-p port] targetURL

- ``mode``: http or tcp
- ``port``: local port, example: 8080


### Proxy a http website:
``gproxy -m http -p 8080 http://github.com/``

then, a simple http server will run on local 8080 port to serve the content of github.com

### Proxy a tcp service:
``gproxy -m tcp -p 2333 github.com:443``

then, all tcp request to localhost:2333 will be sent to github.com:443 and response from github.com:443 will be sent back to client
