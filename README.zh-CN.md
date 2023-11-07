# gproxy

| [English](./README.md) | 简体中文 |

## 使用

### 参数
gproxy [-m 模式] [-p 端口] 目标地址

- ``mode``: http 或 tcp
- ``port``: 本地端口, 比如: 8080


### 代理一个http网站:
``gproxy -m http -p 8080 http://github.com/``

本地8080端口即可以直接访问github

### 代理一个tcp服务:
``gproxy -m tcp -p 2333 github.com:443``

本地2333端口的所有tcp流量都会代理到github.com:443对应的tcp服务